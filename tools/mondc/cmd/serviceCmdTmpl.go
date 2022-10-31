package cmd

import "text/template"

var cmdTemplate, _ = template.New("").Parse(`
package main

import (
	"mond/wind"
	"mond/service/{{.FolderPath}}/handler"
)

func main() {
	wind.InitFrame(handler.NewHook())
}

`)
