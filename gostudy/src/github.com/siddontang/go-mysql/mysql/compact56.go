package mysql

import "fmt"

func TimeStringFromInt64TimePacked(tm int64) string {
	hms := int64(0)
	sign := ""
	if tm < 0 {
		tm = -tm
		sign = "-"
	}

	hms = tm >> 24

	hour := (hms >> 12) % (1 << 10) /* 10 bits starting at 12th */
	minute := (hms >> 6) % (1 << 6) /* 6 bits starting at 6th   */
	second := hms % (1 << 6)        /* 6 bits starting at 0th   */
	secPart := tm % (1 << 24)

	if secPart != 0 {
		return fmt.Sprintf("%s%02d:%02d:%02d.%06d", sign, hour, minute, second, secPart)
	}

	return fmt.Sprintf("%s%02d:%02d:%02d", sign, hour, minute, second)
}
