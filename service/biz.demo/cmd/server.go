package main

import (
	"mond/service/biz.demo/handler"
	"mond/wind"
)

func main() {
	wind.InitFrame(handler.NewHook())
}
