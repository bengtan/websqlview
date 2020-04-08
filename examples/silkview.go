package main

import (
	"flag"
	"fmt"
	"os"

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
	w.Bind("exit", func(w webview.WebView, i int) (err error) {
		if w.GetURI()[0:7] != "file://" {
			return fmt.Errorf("Access denied")
		}
		w.Terminate()
		exitCode = i
		return nil
	})

	sqlite.Init(w)
	defer sqlite.Shutdown()

	w.Run()
	return
}
