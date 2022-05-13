package replication

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	uuid "github.com/satori/go.uuid"
	"github.com/siddontang/go-mysql/client"
	. "github.com/siddontang/go-mysql/mysql"

	"strconv"
)

const (
	// If retry times exceed max retry times, select next data source
	maxRetryTimes = 60
	// Default binlog heartbeat interval (enable by default)
	defaultBinlogHeartbeatInterval = 30
)

var (
	errSyncRunning = errors.New("Sync is running, must Close first")
)

type DataSource struct {
	Host     string
	Port     uint16
	User     string
	Password string
}

func (s *DataSource) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// BinlogSyncerConfig is the configuration for BinlogSyncer.
type BinlogSyncerConfig struct {
	// ServerID is the unique ID in cluster.
	ServerID uint32
	// Flavor is "mysql" or "mariadb", if not set, use "mysql" default.
	Flavor string

	// Master (primary) data source
	// Host is for MySQL server host.
	Host string
	// Port is for MySQL server port.
	Port uint16
	// User is for MySQL user.
	User string
	// Password is for MySQL password.
	Password string

	// Slaves
	Slaves []DataSource

	// Localhost is local hostname if register salve.
	// If not set, use os.Hostname() instead.
	Localhost string

	// Charset is for MySQL client character set
	Charset string

	// SemiSyncEnabled enables semi-sync or not.
	SemiSyncEnabled bool

	// RawModeEnabled is for not parsing binlog event.
	RawModeEnabled bool

	// If not nil, use the provided tls.Config to connect to the database using TLS/SSL.
	TLSConfig *tls.Config

	// Use replication.Time structure for timestamp and datetime.
	// We will use Local location for timestamp and UTC location for datatime.
	ParseTime bool

	LogLevel string

	// RecvBufferSize sets the size in bytes of the operating system's receive buffer associated with the connection.
	RecvBufferSize int

	// KeepAliveSec enable the tcp keepalive, default time is 60*60*2 seconds
	KeepAliveSec int64

	NetWriteTimeoutSec int64
}

// BinlogSyncerEx syncs binlog event from server.
// Extend with gtid full support
type BinlogSyncerEx struct {
	m sync.RWMutex

	cfg *BinlogSyncerConfig

	c *client.Conn

	wg sync.WaitGroup

	parser *BinlogParser

	// For retry sync with position
	nextPos Position
	// For retry sync with gtid
	// Binlog always send next-gtid event, so we should record next-gtid and current-gtid
	// If sync failed and retry, we should use current-gtid to sync
	// Because gtid event is next gtid, so if using the gtid, it will skip the current events until next
	// gtid is received
	currentGTID   string
	nextGTID      string
	enableGTID    bool
	purgedGTIDs   string
	gtidStrBuffer *bytes.Buffer

	running bool

	ctx    context.Context
	cancel context.CancelFunc

	dataSources     []DataSource
	dataSourceIndex int64
	// Last binlog event received time (heartbeat check)
	lastBinlogEventReceivedTime int64
}

// NewBinlogSyncer creates the BinlogSyncer with cfg.
func NewBinlogSyncerEx(cfg *BinlogSyncerConfig) *BinlogSyncerEx {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	log.SetLevelByString(cfg.LogLevel)

	log.Infof("create BinlogSyncer with config %v", cfg)

	b := new(BinlogSyncerEx)

	b.cfg = cfg
	b.parser = NewBinlogParser()
	b.parser.SetRawMode(b.cfg.RawModeEnabled)
	b.parser.SetParseTime(b.cfg.ParseTime)
	b.running = false
	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.gtidStrBuffer = bytes.NewBuffer(make([]byte, 0, 512))
	// Initialize data source
	b.dataSources = make([]DataSource, 0, 1+len(cfg.Slaves))

	var masterDS DataSource
	masterDS.Host = cfg.Host
	masterDS.Port = cfg.Port
	masterDS.User = cfg.User
	masterDS.Password = cfg.Password
	b.dataSources = append(b.dataSources, masterDS)
	// Backup data sources
	if nil != cfg.Slaves && len(cfg.Slaves) > 0 {
		for _, ds := range cfg.Slaves {
			b.dataSources = append(b.dataSources, ds)
		}
	}
	log.Infof("Replication data source initialized with %d node(s)", len(b.dataSources))

	return b
}

// Close closes the BinlogSyncer.
func (b *BinlogSyncerEx) Close() {
	b.m.Lock()
	defer b.m.Unlock()

	b.close()
}

func (b *BinlogSyncerEx) close() {
	if b.isClosed() {
		return
	}

	log.Info("syncer is closing...")

	b.running = false
	b.cancel()

	if b.c != nil {
		b.c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	}

	b.wg.Wait()

	if b.c != nil {
		b.c.Close()
	}

	log.Info("syncer is closed")
}

func (b *BinlogSyncerEx) isClosed() bool {
	select {
	case <-b.ctx.Done():
		return true
	default:
		return false
	}
}

func (b *BinlogSyncerEx) selectNextDataSource() {
	if !b.enableGTID {
		// Only select in gtid replication mode
		return
	}
	atomic.AddInt64(&b.dataSourceIndex, 1)
}

func (b *BinlogSyncerEx) GetDataSourceIndex() int {
	return int(atomic.LoadInt64(&b.dataSourceIndex)) % len(b.dataSources)
}

func (b *BinlogSyncerEx) getDataSource() *DataSource {
	dsi := int(atomic.LoadInt64(&b.dataSourceIndex))
	di := dsi % len(b.dataSources)
	return &b.dataSources[di]
}

func (b *BinlogSyncerEx) enableBinlogHeartbeat(secs int64) error {
	if secs < 0 {
		return errors.Errorf("Invalid binlog heartbeat interval value %d", secs)
	}
	heartbeatInterval := secs
	if 0 == heartbeatInterval {
		heartbeatInterval = defaultBinlogHeartbeatInterval
	}
	// Unit is nano seconds
	if _, err := b.c.Execute(fmt.Sprintf("SET @master_heartbeat_period = %d", heartbeatInterval*1e9)); nil != err {
		return errors.Annotatef(err, "Failed to set master_heartbeat_period: %v", heartbeatInterval*1e9)
	}
	// Once apply heartbeat, enable read timeout of the replication connection
	b.c.SetReadTimeout(time.Second*time.Duration(heartbeatInterval) + time.Second*time.Duration(heartbeatInterval)/2)
	log.Infof("Enable binlog heartbeat with interval %d seconds", heartbeatInterval)

	return nil
}

func (b *BinlogSyncerEx) disableBinlogChecksum() error {
	//for mysql 5.6+, binlog has a crc32 checksum
	//before mysql 5.6, this will not work, don't matter.:-)
	if r, err := b.c.Execute("SHOW GLOBAL VARIABLES LIKE 'BINLOG_CHECKSUM'"); err != nil {
		return errors.Trace(err)
	} else {
		s, _ := r.GetString(0, 1)
		if s != "" {
			// maybe CRC32 or NONE

			// mysqlbinlog.cc use NONE, see its below comments:
			// Make a notice to the server that this client
			// is checksum-aware. It does not need the first fake Rotate
			// necessary checksummed.
			// That preference is specified below.

			if _, err = b.c.Execute(`SET @master_binlog_checksum='NONE'`); err != nil {
				return errors.Trace(err)
			}

			// if _, err = b.c.Execute(`SET @master_binlog_checksum=@@global.binlog_checksum`); err != nil {
			// 	return errors.Trace(err)
			// }

		}
	}
	return nil
}

func (b *BinlogSyncerEx) registerSlave() error {
	if b.c != nil {
		b.c.Close()
	}

	ds := b.getDataSource()
	log.Infof("register slave for master server %s", ds.Address())
	var err error
	b.c, err = client.Connect(ds.Address(), ds.User, ds.Password, "", func(c *client.Conn) {
		//c.TLSConfig = b.cfg.TLSConfig
		c.SetTLSConfig(b.cfg.TLSConfig)
	})
	if err != nil {
		return errors.Trace(err)
	}
	if len(b.cfg.Charset) != 0 {
		b.c.SetCharset(b.cfg.Charset)
	}

	if b.cfg.RecvBufferSize > 0 {
		if tcp, ok := b.c.Conn.Conn.(*net.TCPConn); ok {
			tcp.SetReadBuffer(b.cfg.RecvBufferSize)
		}
	}

	if b.cfg.KeepAliveSec != 0 {
		if tcp, ok := b.c.Conn.Conn.(*net.TCPConn); ok {
			log.Infof("Enable tcp keepalive %vs", b.cfg.KeepAliveSec)
			tcp.SetKeepAlive(true)
			tcp.SetKeepAlivePeriod(time.Second * time.Duration(b.cfg.KeepAliveSec))
		}
	}

	if err = b.disableBinlogChecksum(); nil != err {
		return errors.Trace(err)
	}

	if b.cfg.Flavor == MariaDBFlavor {
		// Refer https://github.com/alibaba/canal/wiki/BinlogChange(MariaDB5&10)
		// Tell the server that we understand GTIDs by setting our slave capability
		// to MARIA_SLAVE_CAPABILITY_GTID = 4 (MariaDB >= 10.0.1).
		if _, err := b.c.Execute("SET @mariadb_slave_capability=4"); err != nil {
			return errors.Errorf("failed to set @mariadb_slave_capability=4: %v", err)
		}
	}

	// Set net_write_timeout session value
	if 0 != b.cfg.NetWriteTimeoutSec {
		if _, err = b.c.Execute(fmt.Sprintf("SET SESSION net_write_timeout=%d", b.cfg.NetWriteTimeoutSec)); nil != err {
			return errors.Errorf("failed to set net_write_timeout: %v", err)
		}
	}

	// If source is mysql, we should get the purged gtid to avoid replication failed
	if b.cfg.Flavor == MySQLFlavor {
		if err = b.initPurgedGtids(); nil != err {
			return errors.Trace(err)
		}
		log.Infof("Get mysql gtid_purged: %s", b.purgedGTIDs)
	}

	// Enable binlog heartbeat
	if err = b.enableBinlogHeartbeat(0); nil != err {
		return errors.Trace(err)
	}

	if err = b.writeRegisterSlaveCommand(); err != nil {
		return errors.Trace(err)
	}

	if _, err = b.c.ReadOKPacket(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncerEx) initPurgedGtids() error {
	res, err := b.c.Execute("SHOW VARIABLES LIKE 'gtid_purged'")
	if nil != err {
		return errors.Trace(err)
	}
	if 0 == res.RowNumber() {
		// The source didn't have a gtid_purged variable
		b.purgedGTIDs = ""
		return nil
	}
	b.purgedGTIDs, err = res.GetString(0, 1)
	if nil != err {
		return errors.Trace(err)
	}
	b.purgedGTIDs = strings.Replace(b.purgedGTIDs, "\n", "", -1)
	return nil
}

func (b *BinlogSyncerEx) enalbeSemiSync() error {
	if !b.cfg.SemiSyncEnabled {
		return nil
	}

	if r, err := b.c.Execute("SHOW VARIABLES LIKE 'rpl_semi_sync_master_enabled';"); err != nil {
		return errors.Trace(err)
	} else {
		s, _ := r.GetString(0, 1)
		if s != "ON" {
			log.Errorf("master does not support semi synchronous replication, use no semi-sync")
			b.cfg.SemiSyncEnabled = false
			return nil
		}
	}

	_, err := b.c.Execute(`SET @rpl_semi_sync_slave = 1;`)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncerEx) prepare() error {
	if b.isClosed() {
		return errors.Trace(ErrSyncClosed)
	}

	if err := b.registerSlave(); err != nil {
		return errors.Trace(err)
	}

	if err := b.enalbeSemiSync(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncerEx) startDumpStream() *BinlogStreamer {
	b.running = true

	s := newBinlogStreamer(0)

	b.wg.Add(1)
	go b.onStream(s)
	return s
}

// StartSync starts syncing from the `pos` position.
func (b *BinlogSyncerEx) StartSync(pos Position) (*BinlogStreamer, error) {
	log.Infof("begin to sync binlog from position %s", pos)

	b.m.Lock()
	defer b.m.Unlock()

	if b.running {
		return nil, errors.Trace(errSyncRunning)
	}

	if err := b.prepareSyncPos(pos); err != nil {
		return nil, errors.Trace(err)
	}

	return b.startDumpStream(), nil
}

// StartSyncGTID starts syncing from the gtid.
func (b *BinlogSyncerEx) StartSyncGTID(gtid string) (*BinlogStreamer, error) {
	log.Infof("begin to sync binlog from GTID %s", gtid)

	b.m.Lock()
	defer b.m.Unlock()

	if b.running {
		return nil, errors.Trace(errSyncRunning)
	}

	b.enableGTID = true
	b.currentGTID = gtid
	b.nextGTID = gtid

	if err := b.prepareSyncGTID(gtid); nil != err {
		return nil, errors.Trace(err)
	}

	return b.startDumpStream(), nil
}

func (b *BinlogSyncerEx) writeBinglogDumpCommand(p Position) error {
	b.c.ResetSequence()

	data := make([]byte, 4+1+4+2+4+len(p.Name))

	pos := 4
	data[pos] = COM_BINLOG_DUMP
	pos++

	binary.LittleEndian.PutUint32(data[pos:], p.Pos)
	pos += 4

	binary.LittleEndian.PutUint16(data[pos:], BINLOG_DUMP_NEVER_STOP)
	pos += 2

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	copy(data[pos:], p.Name)

	return b.c.WritePacket(data)
}

func (b *BinlogSyncerEx) writeBinlogDumpMysqlGTIDCommand(gset GTIDSet) error {
	p := Position{"", 4}
	gtidData := gset.Encode()

	b.c.ResetSequence()

	data := make([]byte, 4+1+2+4+4+len(p.Name)+8+4+len(gtidData))
	pos := 4
	data[pos] = COM_BINLOG_DUMP_GTID
	pos++

	binary.LittleEndian.PutUint16(data[pos:], 0)
	pos += 2

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	binary.LittleEndian.PutUint32(data[pos:], uint32(len(p.Name)))
	pos += 4

	n := copy(data[pos:], p.Name)
	pos += n

	binary.LittleEndian.PutUint64(data[pos:], uint64(p.Pos))
	pos += 8

	binary.LittleEndian.PutUint32(data[pos:], uint32(len(gtidData)))
	pos += 4
	n = copy(data[pos:], gtidData)
	pos += n

	data = data[0:pos]

	return b.c.WritePacket(data)
}

func (b *BinlogSyncerEx) writeBinlogDumpMariadbGTIDCommand(gset GTIDSet) error {
	// Copy from vitess

	startPos := gset.String()

	// Set the slave_connect_state variable before issuing COM_BINLOG_DUMP to
	// provide the start position in GTID form.
	query := fmt.Sprintf("SET @slave_connect_state='%s'", startPos)

	if _, err := b.c.Execute(query); err != nil {
		return errors.Errorf("failed to set @slave_connect_state='%s': %v", startPos, err)
	}

	// Real slaves set this upon connecting if their gtid_strict_mode option was
	// enabled. We always use gtid_strict_mode because we need it to make our
	// internal GTID comparisons safe.
	if _, err := b.c.Execute("SET @slave_gtid_strict_mode=1"); err != nil {
		return errors.Errorf("failed to set @slave_gtid_strict_mode=1: %v", err)
	}

	// Since we use @slave_connect_state, the file and position here are ignored.
	return b.writeBinglogDumpCommand(Position{"", 0})
}

// localHostname returns the hostname that register slave would register as.
func (b *BinlogSyncerEx) localHostname() string {
	if len(b.cfg.Localhost) == 0 {
		h, _ := os.Hostname()
		return h
	}
	return b.cfg.Localhost
}

func (b *BinlogSyncerEx) writeRegisterSlaveCommand() error {
	b.c.ResetSequence()

	hostname := b.localHostname()

	// This should be the name of slave host not the host we are connecting to.
	data := make([]byte, 4+1+4+1+len(hostname)+1+len(b.cfg.User)+1+len(b.cfg.Password)+2+4+4)
	pos := 4

	data[pos] = COM_REGISTER_SLAVE
	pos++

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	// This should be the name of slave hostname not the host we are connecting to.
	data[pos] = uint8(len(hostname))
	pos++
	n := copy(data[pos:], hostname)
	pos += n

	data[pos] = uint8(len(b.cfg.User))
	pos++
	n = copy(data[pos:], b.cfg.User)
	pos += n

	data[pos] = uint8(len(b.cfg.Password))
	pos++
	n = copy(data[pos:], b.cfg.Password)
	pos += n

	binary.LittleEndian.PutUint16(data[pos:], b.cfg.Port)
	pos += 2

	//replication rank, not used
	binary.LittleEndian.PutUint32(data[pos:], 0)
	pos += 4

	// master ID, 0 is OK
	binary.LittleEndian.PutUint32(data[pos:], 0)

	return b.c.WritePacket(data)
}

func (b *BinlogSyncerEx) replySemiSyncACK(p Position) error {
	b.c.ResetSequence()

	data := make([]byte, 4+1+8+len(p.Name))
	pos := 4
	// semi sync indicator
	data[pos] = SemiSyncIndicator
	pos++

	binary.LittleEndian.PutUint64(data[pos:], uint64(p.Pos))
	pos += 8

	copy(data[pos:], p.Name)

	err := b.c.WritePacket(data)
	if err != nil {
		return errors.Trace(err)
	}

	_, err = b.c.ReadOKPacket()
	if err != nil {
	}
	return errors.Trace(err)
}

func (b *BinlogSyncerEx) retrySync() error {
	b.m.Lock()
	defer b.m.Unlock()

	b.parser.Reset()

	if b.enableGTID {
		log.Infof("begin to re-sync from gtid %s", b.currentGTID)

		if err := b.prepareSyncGTID(b.currentGTID); nil != err {
			return errors.Trace(err)
		}
	} else {
		log.Infof("begin to re-sync from pos %s", b.nextPos)

		if err := b.prepareSyncPos(b.nextPos); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (b *BinlogSyncerEx) prepareSyncPos(pos Position) error {
	// always start from position 4
	if pos.Pos < 4 {
		pos.Pos = 4
	}

	if err := b.prepare(); err != nil {
		return errors.Trace(err)
	}

	if err := b.writeBinglogDumpCommand(pos); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncerEx) prepareSyncGTID(gtid string) error {
	gset, err := ParseGTIDSet(b.cfg.Flavor, gtid)

	if nil != err {
		return errors.Trace(err)
	}

	return b.prepareSyncGTIDSet(gset)
}

func (b *BinlogSyncerEx) getReplicationGtidSet(gset GTIDSet) (GTIDSet, error) {
	// Now we can get the gtid_purged after prepare
	if b.cfg.Flavor != MySQLFlavor {
		return gset, nil
	}
	return getReplicationGtidSet(b.purgedGTIDs, gset)
}

func (b *BinlogSyncerEx) prepareSyncGTIDSet(gset GTIDSet) error {
	if err := b.prepare(); err != nil {
		return errors.Trace(err)
	}
	var err error
	gset, err = b.getReplicationGtidSet(gset)
	if nil != err {
		return errors.Trace(err)
	}
	log.Infof("Get real replication gtid set: %s", gset.String())

	if b.cfg.Flavor != MariaDBFlavor {
		// default use MySQL
		err = b.writeBinlogDumpMysqlGTIDCommand(gset)
	} else {
		err = b.writeBinlogDumpMariadbGTIDCommand(gset)
	}

	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *BinlogSyncerEx) onStream(s *BinlogStreamer) {
	defer func() {
		if e := recover(); e != nil {
			s.closeWithError(fmt.Errorf("Err: %v\n Stack: %s", e, Pstack()))
		}
		b.wg.Done()
	}()

	connErrRetryCnt := 0
	for {
		data, err := b.c.ReadPacket()
		if err != nil {
			log.Error(err)

			// we meet connection error, should re-connect again with
			// last nextPos we got.
			if !b.enableGTID {
				// In position replication mode, name is required to retry sync
				if len(b.nextPos.Name) == 0 {
					// we can't get the correct position, close.
					s.closeWithError(err)
					return
				}
			}

			// TODO: add a max retry count.
			for {
				select {
				case <-b.ctx.Done():
					s.close()
					return
				case <-time.After(time.Second):
					if err = b.retrySync(); err != nil {
						log.Errorf("retry sync err: %v, wait 1s and retry again", err)
						if tcpErr, ok := errors.Cause(err).(net.Error); ok && nil != tcpErr {
							connErrRetryCnt++
						} else if tcpErr, ok := errors.Cause(err).(*net.OpError); ok && nil != tcpErr {
							connErrRetryCnt++
						}
						if connErrRetryCnt%maxRetryTimes == 0 &&
							connErrRetryCnt != 0 {
							log.Warnf("Select next data source because of current data source down")
							b.selectNextDataSource()
							log.Infof("Select next data source: %s", b.getDataSource().Address())
						}
						continue
					}
				}

				break
			}

			// we connect the server and begin to re-sync again.
			connErrRetryCnt = 0
			continue
		}

		switch data[0] {
		case OK_HEADER:
			if err = b.parseEvent(s, data); err != nil {
				s.closeWithError(err)
				return
			}
		case ERR_HEADER:
			err = b.c.HandleErrorPacket(data)
			s.closeWithError(err)
			return
		case EOF_HEADER:
			// Refer http://dev.mysql.com/doc/internals/en/packet-EOF_Packet.html
			// In the MySQL client/server protocol, EOF and OK packets serve the same purpose.
			// Some users told me that they received EOF packet here, but I don't know why.
			// So we only log a message and retry ReadPacket.
			log.Info("receive EOF packet, retry ReadPacket")
			continue
		default:
			log.Errorf("invalid stream header %c", data[0])
			continue
		}
	}
}

func (b *BinlogSyncerEx) parseEvent(s *BinlogStreamer, data []byte) error {
	//skip OK byte, 0x00
	data = data[1:]

	needACK := false
	if b.cfg.SemiSyncEnabled && (data[0] == SemiSyncIndicator) {
		needACK = (data[1] == 0x01)
		//skip semi sync header
		data = data[2:]
	}

	e, err := b.parser.Parse(data)
	if err != nil {
		return errors.Trace(err)
	}

	if e.Header.LogPos > 0 {
		// Some events like FormatDescriptionEvent return 0, ignore.
		b.nextPos.Pos = e.Header.LogPos
	}

	// Record binlogs position or gtid
	switch fv := e.Event.(type) {
	case *RotateEvent:
		{
			b.nextPos.Name = string(fv.NextLogName)
			b.nextPos.Pos = uint32(fv.Position)
			log.Infof("rotate to position %s", b.nextPos)
		}
	case *MariadbGTIDEvent:
		{
			b.nextGTID = fv.GTID.String()
		}
	case *GTIDEvent:
		{
			u, err := uuid.FromBytes(fv.SID)
			if nil != err {
				return errors.Trace(err)
			}
			serverUUID := u.String()
			// Gtid interval is [n:m], n-m is executed, so we record the next gtid (actually previous gtid) to current gtid
			b.gtidStrBuffer.Reset()
			b.gtidStrBuffer.WriteString(serverUUID)
			b.gtidStrBuffer.WriteString(":1-")
			b.gtidStrBuffer.WriteString(strconv.FormatInt(fv.GNO, 10))
			b.nextGTID = b.gtidStrBuffer.String()
		}
	case *XIDEvent:
		{
			// When we meet xid event, save next gtid position to current
			// Xid represents a transacation is committed, so we should use the previous next gtid
			b.currentGTID = b.nextGTID
		}
	case *HeartbeatEvent:
		{
			b.lastBinlogEventReceivedTime = int64(e.Header.Timestamp)
			log.Debugf("Binlog heartbeat")
		}
	}

	needStop := false
	select {
	case s.ch <- e:
	case <-b.ctx.Done():
		needStop = true
	}

	if needACK {
		err := b.replySemiSyncACK(b.nextPos)
		if err != nil {
			return errors.Trace(err)
		}
	}

	if needStop {
		return errors.New("sync is been closing...")
	}

	return nil
}
