package replication

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/juju/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/siddontang/go-mysql/mysql"
)

type PreviousGtidsLogEvent struct {
	numberOfSID uint64
	GtidSet     mysql.GTIDSet
}

func (e *PreviousGtidsLogEvent) Decode(data []byte) error {
	length := len(data)
	if length < 8 {
		return errors.New("EOF when reading length of sid")
	}
	// Read number of SIDs
	pos := 0
	e.numberOfSID = binary.LittleEndian.Uint64(data)
	pos += 8

	for i := uint64(0); i < e.numberOfSID; i++ {
		if length-pos < 16+8 {
			return errors.New("EOF when reading sid")
		}
		// Copy 16 bytes as sid
		var sid [16]byte
		copy(sid[:], data[pos:])
		pos += 16
		// Read 8 bytes as interval
		var intervals uint64
		intervals = binary.LittleEndian.Uint64(data[pos:])
		pos += 8
		// Read all intervals
		if length-pos < 2*8*int(intervals) {
			return errors.New("EOF when reading intervals")
		}
		var last int64
		for i := uint64(0); i < intervals; i++ {
			var start, end int64
			start = int64(binary.LittleEndian.Uint64(data[pos:]))
			pos += 8
			end = int64(binary.LittleEndian.Uint64(data[pos:]))
			pos += 8

			if start <= last || end <= start {
				return errors.Errorf("Bad intervals, last=%d, start=%d, end=%d",
					last, start, end)
			}
			last = end

			if start > end-1 {
				continue
			}
			u, err := uuid.FromBytes(sid[:])
			if nil != err {
				return errors.Trace(err)
			}
			gtid := fmt.Sprintf("%s:%d-%d", u.String(), start, end-1)
			if nil == e.GtidSet {
				e.GtidSet, err = mysql.ParseGTIDSet(mysql.MySQLFlavor, gtid)
				if nil != err {
					return errors.Trace(err)
				}
			} else {
				e.GtidSet.Update(gtid)
			}
		}
	}

	return nil
}

func (e *PreviousGtidsLogEvent) Dump(w io.Writer) {
	if nil == e.GtidSet {
		fmt.Fprint(w, "Null")
	} else {
		fmt.Fprint(w, e.GtidSet.String())
	}
	fmt.Fprintln(w)
}
