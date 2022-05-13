package replication

import (
	"strings"
	"testing"

	"github.com/siddontang/go-mysql/mysql"
)

func TestMySQLGtidHasSID(t *testing.T) {
	type testcase struct {
		gset string
		sid  string
		has  bool
	}

	tcs := []testcase{
		{
			gset: "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-100",
			sid:  "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68",
			has:  true,
		},
		{
			gset: "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-100,1a46f393-ec67-11e7-9dbb-6c92bf5b8b69:1-100",
			sid:  "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68",
			has:  true,
		},
		{
			gset: "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-100,1a46f393-ec67-11e7-9dbb-6c92bf5b8b69:1-100",
			sid:  "1a46f393-ec67-11e7-9dbb-6c92bf5b8b67",
			has:  false,
		},
		{
			gset: "",
			sid:  "1a46f393-ec67-11e7-9dbb-6c92bf5b8b67",
			has:  false,
		},
	}

	for _, tc := range tcs {
		gset, err := mysql.ParseGTIDSet(mysql.MySQLFlavor, tc.gset)
		if nil != err {
			t.Fatal(err)
		}
		if gset.HasSID(tc.sid) != tc.has {
			t.Errorf("Gset: %v, sid: %v, has: %v", tc.gset, tc.sid, tc.has)
		}
	}
}

func TestMariaGtidMerge(t *testing.T) {
	type testcase struct {
		gset string
		sid  string
		has  bool
	}

	tcs := []testcase{
		{
			gset: "1018006083-1018006083-100",
			sid:  "1018006083",
			has:  true,
		},
		{
			gset: "1018006083-1018006083-100,1018006083-1018006084-100",
			sid:  "1018006083",
			has:  true,
		},
		{
			gset: "1018006083-1018006083-100,1018006083-1018006084-100",
			sid:  "1018006082",
			has:  false,
		},
		{
			gset: "",
			sid:  "1018006083",
			has:  false,
		},
	}

	for _, tc := range tcs {
		gset, err := mysql.ParseGTIDSet(mysql.MariaDBFlavor, tc.gset)
		if nil != err {
			t.Fatal(err)
		}
		if gset.HasSID(tc.sid) != tc.has {
			t.Errorf("Gset: %v, sid: %v, has: %v", tc.gset, tc.sid, tc.has)
		}
	}
}

func TestGetReplicationGtidSet(t *testing.T) {
	type testcase struct {
		purged   string
		gset     string
		result   bool
		wantGset string
	}

	tcs := []testcase{
		{
			purged:   "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1000",
			gset:     "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-999",
			result:   false,
			wantGset: "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1000",
		},
		{
			purged:   "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1000",
			gset:     "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1000",
			result:   true,
			wantGset: "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1000",
		},
		{
			purged:   "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1000",
			gset:     "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1250",
			result:   true,
			wantGset: "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1250",
		},
		{
			purged:   "",
			gset:     "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1250",
			result:   true,
			wantGset: "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1250",
		},
		{
			purged:   "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1000",
			gset:     "",
			result:   true,
			wantGset: "1a46f393-ec67-11e7-9dbb-6c92bf5b8b68:1-1000",
		},
		{
			purged:   "",
			gset:     "",
			result:   true,
			wantGset: "",
		},
	}

	for i := range tcs {
		tc := &tcs[i]
		gset, err := mysql.ParseMysqlGTIDSet(tc.gset)
		if nil != err {
			t.Errorf("Invalid gset %s: %v", tc.gset, err)
			continue
		}
		res, err := getReplicationGtidSet(tc.purged, gset)
		if nil != err {
			if !tc.result {
				if strings.Contains(err.Error(), "Request gtid") {
					continue
				}
			}
			t.Errorf("Error :%v", err)
			continue
		}
		if res.String() != tc.wantGset {
			t.Errorf("Want %v, got %v", tc.wantGset, res.String())
		}
	}
}
