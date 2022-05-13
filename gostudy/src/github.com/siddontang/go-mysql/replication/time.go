package replication

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

var (
	fracTimeFormat []string
	usecFormat     []string
)

// fracTime is a help structure wrapping Golang Time.
type fracTime struct {
	time.Time

	// Dec must in [0, 6]
	Dec int

	timestampStringLocation *time.Location
}

func (t fracTime) String() string {
	tt := t.Time
	if nil != t.timestampStringLocation {
		tt = tt.In(t.timestampStringLocation)
	}
	return t.Format(fracTimeFormat[t.Dec])
}

func formatZeroTime(frac int, dec int) string {
	if dec == 0 {
		return "0000-00-00 00:00:00"
	}

	s := fmt.Sprintf("0000-00-00 00:00:00.%06d", frac)

	// dec must < 6, if frac is 924000, but dec is 3, we must output 924 here.
	return s[0 : len(s)-(6-dec)]
}

func formatBeforeUnixZeroTime(year, month, day, hour, minute, second, frac, dec int) string {
	if dec == 0 {
		return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
	}

	s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d", year, month, day, hour, minute, second, frac)

	// dec must < 6, if frac is 924000, but dec is 3, we must output 924 here.
	return s[0 : len(s)-(6-dec)]
}

const (
	mysqlTimestampNone     = -2
	mysqlTimestampError    = -1
	mysqlTimestampDate     = 0
	mysqlTimestampDatetime = 1
	mysqlTimestampTime     = 2

	DATETIME_MAX_DECIMALS = 6
)

var log10Int []uint64 = []uint64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
	10000000000000000000,
}

type mysqlTime struct {
	year       uint32
	month      uint32
	day        uint32
	hour       uint32
	minute     uint32
	second     uint32
	secondPart uint64 /* Microseconds */
	neg        bool
	timeType   int
	dec        int
}

func (m *mysqlTime) String() string {
	buf := bytes.NewBuffer(nil)
	var temp, temp2 uint32

	// Year
	temp = m.year / 100
	buf.Write([]byte{'0' + byte(temp/10), '0' + byte(temp%10)})
	temp = m.year % 100
	buf.Write([]byte{'0' + byte(temp/10), '0' + byte(temp%10), '-'})
	// Month
	temp = m.month
	temp2 = temp / 10
	temp = temp - temp2*10
	buf.Write([]byte{'0' + byte(temp2), '0' + byte(temp), '-'})
	// Day
	temp = m.day
	temp2 = temp / 10
	temp = temp - temp2*10
	buf.Write([]byte{'0' + byte(temp2), '0' + byte(temp), ' '})
	// Hour
	temp = m.hour
	temp2 = temp / 10
	temp = temp - temp2*10
	buf.Write([]byte{'0' + byte(temp2), '0' + byte(temp), ':'})
	// Minute
	temp = m.minute
	temp2 = temp / 10
	temp = temp - temp2*10
	buf.Write([]byte{'0' + byte(temp2), '0' + byte(temp), ':'})
	// Second
	temp = m.second
	temp2 = temp / 10
	temp = temp - temp2*10
	buf.Write([]byte{'0' + byte(temp2), '0' + byte(temp)})

	// Microseconds
	if 0 != m.dec {
		usecStr := fmt.Sprintf(usecFormat[m.dec],
			uint64(m.secondPart)/uint64(log10Int[DATETIME_MAX_DECIMALS-m.dec]))
		buf.WriteString(usecStr)
	}
	return buf.String()
}

func init() {
	fracTimeFormat = make([]string, 7)
	usecFormat = make([]string, 7)
	fracTimeFormat[0] = "2006-01-02 15:04:05"

	for i := 1; i <= 6; i++ {
		usecFormat[i] = fmt.Sprintf(".%%0%dd", i)
		fracTimeFormat[i] = fmt.Sprintf("2006-01-02 15:04:05.%s", strings.Repeat("0", i))
	}
}
