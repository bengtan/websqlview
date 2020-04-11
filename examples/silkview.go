package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bengtan/silk/dialog"
	"github.com/bengtan/silk/native"
	"github.com/bengtan/silk/sqlite"
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

	native.Init(w, &exitCode)
	sqlite.Init(w)
	defer sqlite.Shutdown()
	dialog.Init(w)

	w.Run()
	return
}
