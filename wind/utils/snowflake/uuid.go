package utils

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

//获取没有 破折号的 UUID
func GetNoDashUUIDStr() string {
	uuidStr := uuid.NewV4()
	str := strings.ReplaceAll(uuidStr.String(), "-", "")
	return str
}
