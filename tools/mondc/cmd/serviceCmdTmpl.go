package cmd

import "text/template"

var cmdTemplate, _ = template.New("").Parse(`
package main

import (
	"meta/frame"
	"meta/service/{{.FolderPath}}/handler"
)

func main() {
	frame.InitFrame(handler.NewHook())
}

`)
