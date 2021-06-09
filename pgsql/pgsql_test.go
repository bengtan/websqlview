package pgsql

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/bengtan/websqlview/native"
	"github.com/bengtan/websqlview/webviewex"
	"github.com/webview/webview"
)

// Note: This assumes that an postgresql database 'test' exists.
const DB_DSN = "test:test@localhost:5432/test"

var (
	wv            webview.WebView
	sema          sync.Mutex
	dummyExitCode int
	passFn        func()
	failFn        func(s string)
)

func TestMain(m *testing.M) {
	os.Exit(_testMain(m))
}

func _testMain(m *testing.M) (result int) {
	wv = webview.New(true)
	defer wv.Destroy()
	wv.SetTitle("Testing")
	wv.SetSize(800, 600, webview.HintNone)

	ex := webviewex.New(wv)
	native.Init(ex, &dummyExitCode)
	Init(ex)
	defer Shutdown()

	// Override with pgsql.js
	pgsqlJs, err := ioutil.ReadFile("pgsql.js")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		wv.Init(string(pgsqlJs))
	}

	// Inject DB_DSN
	wv.Init(`const DB_DSN = "` + DB_DSN + `"`)

	// Future callbacks
	wv.Bind("pass", func() {
		if passFn != nil {
			passFn()
		}
	})
	wv.Bind("fail", func(s string) {
		if failFn != nil {
			failFn(s)
		}
	})

	// Some coordination with semaphores and asynchronous magic
	wv.Bind("__TestHarnessInit", func() {
		sema.Unlock()
	})
	wv.Init("window.onload = function() { __TestHarnessInit(); }")
	wv.Navigate("data:text/text,TestHarness")

	result = 1
	sema.Lock()
	go func() {
		sema.Lock()
		result = m.Run()
		sema.Unlock()
		if result == 0 {
			wv.Terminate()
		}
	}()

	wv.Run()
	return
}

func TestJS(t *testing.T) {
	matches, _ := filepath.Glob("test/*.js")
	for _, filename := range matches {
		fmt.Printf("Testing: %s", filename)
		failure := testOneJS(filename)
		if failure == "" {
			fmt.Println(" - pass")
		} else {
			fmt.Println(" - FAIL")
			t.Errorf(failure)
			return
		}
	}
}

func testOneJS(filename string) (failure string) {
	var mutex sync.Mutex
	failure = ""

	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return err.Error()
	}

	mutex.Lock()
	passFn = func() {
		mutex.Unlock()
	}
	failFn = func(s string) {
		failure = s
		mutex.Unlock()
	}

	// Execute on main thread otherwise weird stuff happens (on some platforms)
	wv.Dispatch(func() {
		wv.SetTitle("Testing: " + filename)
		wv.Eval(string(text))
		wv.Eval(`runTest().then(errString => {
			errString ? fail(errString.toString()) : pass()
		}).catch(e => {
			fail('Unhandled promise rejection: ' + e.toString())
		})`)
	})

	// Wait for the result
	mutex.Lock()
	{
	} // Avoid warning: empty critical section (SA2001)
	mutex.Unlock()

	passFn = nil
	failFn = nil
	return
}
