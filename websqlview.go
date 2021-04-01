package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bengtan/websqlview/dialog"
	"github.com/bengtan/websqlview/mysql"
	"github.com/bengtan/websqlview/native"
	"github.com/bengtan/websqlview/sqlite"
	"github.com/bengtan/websqlview/webviewex"
	"github.com/webview/webview"
)

func processURI(uri string) string {
	if _, error := os.Stat(uri); error == nil {
		if abs, error := filepath.Abs(uri); error == nil {
			return "file://" + abs
		}
	}
	return uri
}

func main() {
	os.Exit(mainExitCode())
}

func mainExitCode() (exitCode int) {
	debug := true

	flag.BoolVar(&debug, "d", debug, "debug")
	flag.BoolVar(&debug, "debug", debug, "debug")
	flag.Parse()

	arg := flag.Arg(0)
	if arg == "" {
		fmt.Println("Please supply a filename or URI (ie. file:///...)")
		return
	}
	uri := processURI(arg)

	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle(uri)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(uri)

	ex := webviewex.New(w)
	native.Init(ex, &exitCode)
	sqlite.Init(ex)
	defer sqlite.Shutdown()
	mysql.Init(ex)
	defer mysql.Shutdown()
	dialog.Init(ex)

	w.Run()
	return
}
