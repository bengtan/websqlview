package native

import (
	"fmt"
	"os"
	"reflect"

	"github.com/bengtan/websqlview/webviewex"
)

var exitCodePtr *int

// Init binds some custom functionality
func Init(ex *webviewex.WebViewEx, p *int) {
	exitCodePtr = p
	ex.W.Bind("_nativeMux", func(op string, args ...interface{}) (result interface{}, err error) {
		return mux(ex, op, args...)
	})
	ex.W.Init(_nativeJs)
}

func mux(ex *webviewex.WebViewEx, op string, args ...interface{}) (result interface{}, err error) {
	if ex.URI[0:7] != "file://" {
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
			ex.W.Terminate()
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
