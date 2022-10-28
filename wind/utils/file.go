package utils

import (
	"go/format"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func FormatAndWrite(f *os.File, bs []byte) {
	formatted, err := format.Source(bs)
	MustNil(err)
	f.Write(formatted)
}
