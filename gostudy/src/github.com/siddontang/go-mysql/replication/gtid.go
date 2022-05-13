package replication

import (
	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/mysql"
)

func getReplicationGtidSet(purgedGtid string, gset mysql.GTIDSet) (mysql.GTIDSet, error) {
	// No binlog was purged
	if purgedGtid == "" {
		return gset, nil
	}

	purgedGset, err := mysql.ParseGTIDSet(mysql.MySQLFlavor, purgedGtid)
	if gset.String() == "" {
		// Starts with the purged gtid
		return purgedGset, nil
	}

	// Directly add to purged gset
	purgedMysqlGset := purgedGset.(*mysql.MysqlGTIDSet)
	uset, err := mysql.ParseUUIDSet(gset.String())
	if nil != err {
		return nil, errors.Trace(err)
	}
	// We must check the SID of the uuidSet in purged gtid to avoid forming a hole
	if pset, ok := purgedMysqlGset.Sets[uset.SID.String()]; ok {
		if pset.Contain(uset) {
			// If pset contains uset, it represents the current start gtid was purged.
			// Simply add the uset to purgedMysqlGset will lose transactions.
			// Because the transactions in uset was executed, so if the maximum gtid is equal
			// to the maximum gtid in pset, it should be ok.
			maxInUset := int64(0)
			for _, iv := range uset.Intervals {
				if iv.Stop > maxInUset {
					maxInUset = iv.Stop
				}
			}
			maxInPset := int64(0)
			for _, iv := range pset.Intervals {
				if iv.Stop > maxInPset {
					maxInPset = iv.Stop
				}
			}
			if maxInPset != maxInUset {
				return nil, errors.Errorf("Request gtid %s is purged by the gset %s",
					uset.String(), pset.String())
			}
		}
	}
	purgedMysqlGset.AddSet(uset)
	return purgedMysqlGset, nil
}
