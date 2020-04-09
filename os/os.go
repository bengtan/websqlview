package os

import (
	"fmt"
	"os"
	"reflect"

	"github.com/zserge/webview"
)

var exitCodePtr *int

// Init binds some custom functionality
func Init(w webview.WebView, p *int) {
	exitCodePtr = p
	w.Bind("_osMux", mux)
	w.Init(_osJs)
}

func mux(w webview.WebView, op string, args ...interface{}) (result interface{}, err error) {
	if w.GetURI()[0:7] != "file://" {
		return nil, fmt.Errorf("Access denied")
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s: %v", op, e)
		}
	}()

	switch op {
	case "exit":
		if exitCode, ok := args[0].(float64); ok {
			w.Terminate()
			*exitCodePtr = int(exitCode)
			return nil, nil
		}
	case "remove":
		if name, ok := args[0].(string); ok {
			return nil, os.Remove(name)
		}
	}

	signature := []string{}
	for _, arg := range args {
		signature = append(signature, reflect.TypeOf(arg).Name())
	}
	return nil, fmt.Errorf("Unknown operation %s with signature %v", op, signature)
}
