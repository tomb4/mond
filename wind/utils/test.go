package utils

import (
	"os"
	"regexp"
)

func ChangePosition() {
	pwd, _ := os.Getwd()
	re := regexp.MustCompile(`(.*?/service/.+?/).*`)
	if data := re.FindAllStringSubmatch(pwd, -1); len(data) > 0 && len(data[0]) > 1 {
		_ = os.Chdir(data[0][1])
	}
}
