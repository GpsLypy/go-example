package replication

import (
	"fmt"
	"testing"

	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/mysql"
)

func doMerge(start, end int, mstart, mend int) error {
	gset, err := mysql.ParseMysqlGTIDSet(fmt.Sprintf("8dcb7f49-6d43-11e7-a0ae-6c0b84d53df8:%d-%d", start, end))
	if nil != err {
		return err
	}
	mgset := gset.(*mysql.MysqlGTIDSet)
	uset, _ := mysql.ParseUUIDSet(fmt.Sprintf("8dcb7f49-6d43-11e7-a0ae-6c0b84d53df8:%d-%d", mstart, mend))
	mgset.AddSet(uset)

	wantStart := start
	if mstart < wantStart {
		wantStart = mstart
	}
	wantEnd := end
	if mend > wantEnd {
		wantEnd = mend
	}
	if mgset.String() != fmt.Sprintf("8dcb7f49-6d43-11e7-a0ae-6c0b84d53df8:%d-%d", wantStart, wantEnd) {
		return errors.New("Merge failed")
	}
	return nil
}

func TestMergeGtid(t *testing.T) {
	if err := doMerge(1, 100, 1, 1000); nil != err {
		t.Error(err)
	}
	if err := doMerge(1, 1000, 1, 100); nil != err {
		t.Error(err)
	}
}
