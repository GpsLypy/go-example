package mysqlib

import (
	"database/sql"

	"github.com/siddontang/go-mysql/replication"
)

type Position map[string]string

type InboundEvent struct {
	Event *replication.BinlogEvent
	// From which
	Index int
	// Error
	Err error
}

type IInbound interface {
	Start([]Position) error
	Close()
	Pause() error
	Resume([]Position) error
	AddAddrWrapper(string, string, string, string, string, int64) error
	AddFixedHostSource(string, int, string, string) error
	GetEventChannel() <-chan *InboundEvent
	GetDataSourceCount() int
	GetStartPoint(int) Position
	GetQueryDB(int) *sql.DB
	GetAddress(int) string
}
