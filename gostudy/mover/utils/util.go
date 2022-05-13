package utils

import (
	"app/tools/mover/logging"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	InvalidGroupIndex = 99999
	SlaveRole         = "SLAVE"
	MasterRole        = "MASTER"
	SlaveToolRole     = "SLAVE_TOOL"
	SlaveBakRole      = "SLAVE_BAK"
)

func CastUnsigned(data interface{}, unsigned bool) interface{} {
	if !unsigned {
		return data
	}

	switch v := data.(type) {
	case int:
		return uint(v)
	case int8:
		return uint8(v)
	case int16:
		return uint16(v)
	case int32:
		return uint32(v)
	case int64:
		return strconv.FormatUint(uint64(v), 10)
	}

	return data
}

func FormatFieldValue(data interface{}) string {
	if nil == data {
		return "NULL"
	} else {
		switch fv := data.(type) {
		case string:
			{
				return "'" + mysqlQuoteString(fv) + "'"
			}
		case []byte:
			{
				if nil == fv {
					return "NULL"
				} else {
					sv := string(fv)
					return "'" + mysqlQuoteString(sv) + "'"
				}
			}
		default:
			{
				return fmt.Sprintf("%v", fv)
			}
		}
	}
}

func mysqlQuoteString(i string) string {
	buf := bytes.NewBuffer(nil)
	rs := []rune(i)

	for i, v := range rs {
		var escape rune

		switch v {
		case 0:
			{
				escape = '0'
			}
		case '\n':
			{
				escape = 'n'
			}
		case '\r':
			{
				escape = 'r'
			}
		case '\\':
			{
				// Check if has the next rune
				if i == len(rs)-1 {
					escape = '\\'
				} else {
					if rs[i+1] != '%' && rs[i+1] != '_' {
						escape = '\\'
					}
				}
			}
		case '\'':
			{
				escape = '\''
			}
		case '"':
			{
				escape = '"'
			}
		case '\032':
			{
				escape = 'Z'
			}
		}

		if 0 != escape {
			buf.WriteRune('\\')
			buf.WriteRune(escape)
		} else {
			buf.WriteRune(v)
		}
	}

	return buf.String()
}

func GetAllInterface() (map[string]string, error) {
	ifaces, err := net.Interfaces()
	if nil != err {
		return nil, err
	}

	var ifaMap = make(map[string]string)
	var ifaName, ipv4 string
	for _, ifa := range ifaces {
		ifaName = ifa.Name
		addrs, err := ifa.Addrs()
		if nil != err {
			//LogError(EC_Net_PKG_ERR, "Get interface addr error, err: %v", err)
			logging.LogError(logging.EC_Net_PKG_ERR, "Get interface addr error, err: %v", err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}

			if ip.To4() != nil {
				ipv4 = ip.String()
				ifaMap[ifaName] = ipv4
			}
		}
	}

	return ifaMap, nil
}

func Bool2int(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func ByteToString(data []uint8) string {
	strByte := make([]byte, len(data))
	for i, v := range data {
		strByte[i] = byte(v)
	}
	return string(strByte)
}

func StringToInt64(bs string) (int64, error) {
	i, err := strconv.ParseInt(bs, 10, 64)
	return i, err
}

func StringToUInt64(bs string) (uint64, error) {
	i, err := strconv.ParseUint(bs, 10, 64)
	return i, err
}

func StringToFloat64(bs string) (float64, error) {
	v, err := strconv.ParseFloat(bs, 64)
	return v, err
}

func ReadData(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

// setBytesBit set the byte array bit, bit is the index of the byte array
// Note: index is from 0
func setBytesBit(buf []byte, bit int) error {
	byteIndex := bit / 8
	bitIndex := uint8(bit - byteIndex*8)

	if byteIndex >= len(buf) {
		return errors.New("Buffer overflow")
	}
	prevByte := buf[byteIndex]
	mask := uint8(1 << (8 - bitIndex - 1))
	prevByte |= mask
	buf[byteIndex] = prevByte

	return nil
}

func IsMaster(role string) bool {
	return strings.ToUpper(role) == MasterRole
}

func IsSlave(role string) bool {
	return strings.ToUpper(role) == SlaveRole
}

func IsToolSlave(role string) bool {
	return strings.ToUpper(role) == SlaveToolRole
}

func IsBakSlave(role string) bool {
	return strings.ToUpper(role) == SlaveBakRole
}
