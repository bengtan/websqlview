package sqlite

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"github.com/zserge/webview"
)

var (
	wd string
)

func TestMain(m *testing.M) {
	wd, _ = os.Getwd()
	os.Exit(m.Run())
}

func TestJS(t *testing.T) {
	matches, _ := filepath.Glob("test/*.js")
	for _, filename := range matches {
		log.Println("Testing:", filename)
		testOneJS(t, filename)
	}
}

func testOneJS(t *testing.T, filename string) {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Testing: " + filename)
	w.SetSize(800, 600, webview.HintNone)
	w.Bind("pass", func() {
		w.Terminate()
	})
	w.Bind("fail", func(s string) {
		w.Terminate()
		t.Errorf(s)
	})

	Init(w)
	defer Shutdown()

	w.Init(string(text))
	w.Navigate("file://" + wd + "/test/testHarness.html")
	w.Run()
	return
}
