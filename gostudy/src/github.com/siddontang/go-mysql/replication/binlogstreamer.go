package replication

import (
	"lib/localcmd"

	"golang.org/x/net/context"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/siddontang/go/sync2"
)

var (
	ErrNeedSyncAgain = errors.New("Last sync error or closed, try sync and get event again")
	ErrSyncClosed    = errors.New("Sync was closed")
)

// BinlogStreamer gets the streaming event.
type BinlogStreamer struct {
	ch             chan *BinlogEvent
	ech            chan error
	err            error
	dataBytesParse sync2.AtomicInt64
}

// GetEvent gets the binlog event one by one, it will block until Syncer receives any events from MySQL
// or meets a sync error. You can pass a context (like Cancel or Timeout) to break the block.
func (s *BinlogStreamer) GetEvent(ctx context.Context) (*BinlogEvent, error) {
	if s.err != nil {
		return nil, ErrNeedSyncAgain
	}

	select {
	case c := <-s.ch:
		return c, nil
	case s.err = <-s.ech:
		return nil, s.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *BinlogStreamer) GetEventWithLocalCommand(ctx context.Context,
	localCmdCh chan *localcmd.LocalCommand) (*BinlogEvent, *localcmd.LocalCommand, error) {
	if s.err != nil {
		return nil, nil, ErrNeedSyncAgain
	}

	select {
	case c := <-s.ch:
		return c, nil, nil
	case c := <-localCmdCh:
		return nil, c, nil
	case s.err = <-s.ech:
		return nil, nil, s.err
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
}

func (s *BinlogStreamer) GetLen() int {
	return len(s.ch)
}

func (s *BinlogStreamer) close() {
	s.closeWithError(ErrSyncClosed)
}

func (s *BinlogStreamer) closeWithError(err error) {
	if err == nil {
		err = ErrSyncClosed
	}
	log.Errorf("close sync with err: %v", err)
	select {
	case s.ech <- err:
	default:
	}
}

func (s *BinlogStreamer) setDataBytesParsed(parsed int64) {
	s.dataBytesParse.Set(parsed)
}

func (s *BinlogStreamer) GetDataBytesParsed() int {
	return int(s.dataBytesParse.Get())
}

func newBinlogStreamer(size int) *BinlogStreamer {
	if 0 == size {
		size = 10240
	}
	s := new(BinlogStreamer)

	s.ch = make(chan *BinlogEvent, size)
	s.ech = make(chan error, 4)

	return s
}
