package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bengtan/websqlview/dialog"
	"github.com/bengtan/websqlview/native"
	"github.com/bengtan/websqlview/sqlite"
	"github.com/bengtan/websqlview/webviewex"
	"github.com/zserge/webview"
)

func main() {
	os.Exit(mainExitCode())
}

func mainExitCode() (exitCode int) {
	debug := true

	flag.BoolVar(&debug, "d", debug, "debug")
	flag.BoolVar(&debug, "debug", debug, "debug")
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		fmt.Println("Please supply a URI")
		return
	}

	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle(filename)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(filename)

	ex := webviewex.New(w)
	native.Init(ex, &exitCode)
	sqlite.Init(ex)
	defer sqlite.Shutdown()
	dialog.Init(ex)

	w.Run()
	return
}
