package utils

import (
	"fmt"
	"strings"
	"time"
)

//获取当前时间 ms
func CurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
func CurrentSecond() int64 {
	return time.Now().Unix()
}

func UseTimeToStr(useTime time.Duration) string {
	useTimeStr := ""
	//useTime := fmt.Sprintf("%.2f", float64(time.Now().UnixNano()-start.UnixNano())/float64(time.Microsecond))
	if useTime < time.Millisecond {
		useTimeStr = fmt.Sprintf("%.2fus", float64(useTime)/float64(time.Microsecond))
	} else if useTime < time.Second {
		useTimeStr = fmt.Sprintf("%.2fms", float64(useTime)/float64(time.Millisecond))
	} else {
		useTimeStr = fmt.Sprintf("%.2fs", float64(useTime)/float64(time.Second))
	}
	return useTimeStr
}

func FormatTime(t time.Time, format string) string {
	format = strings.ReplaceAll(format, "YYYY", "2006")
	format = strings.ReplaceAll(format, "YY", "06")
	format = strings.ReplaceAll(format, "MM", "01")
	format = strings.ReplaceAll(format, "DD", "02")
	format = strings.ReplaceAll(format, "HH", "15")
	format = strings.ReplaceAll(format, "mm", "04")
	format = strings.ReplaceAll(format, "ss", "05")
	return t.Format(format)
}