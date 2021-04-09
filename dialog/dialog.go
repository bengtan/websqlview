package dialog

import (
	_ "embed"
	"fmt"
	"reflect"

	"github.com/bengtan/websqlview/webviewex"
	"github.com/sqweek/dialog"
)

//go:embed dialog.js
var _dialogJs string

// Init binds the js->go bridge for dialog functionality
func Init(ex *webviewex.WebViewEx) {
	ex.W.Bind("_dialogMux", func(op string, args ...interface{}) (result interface{}, err error) {
		return mux(ex, op, args...)
	})
	ex.W.Init(_dialogJs)
}

func mux(ex *webviewex.WebViewEx, op string, args ...interface{}) (result interface{}, err error) {
	if ex.URI[0:7] != "file://" && ex.URI[0:5] != "data:" {
		return nil, fmt.Errorf("Access denied")
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s: %v", op, e)
		}
	}()

	switch op {
	case "message":
		if config, ok := args[0].(map[string]interface{}); ok {
			return message(config)
		}
	case "directory":
		if config, ok := args[0].(map[string]interface{}); ok {
			return directory(config)
		}
	case "file":
		if config, ok := args[0].(map[string]interface{}); ok {
			return file(config)
		}
	}

	signature := []string{}
	for _, arg := range args {
		signature = append(signature, reflect.TypeOf(arg).Name())
	}
	return nil, fmt.Errorf("Unknown operation %s with signature %v", op, signature)
}

func message(config map[string]interface{}) (result interface{}, err error) {
	if message, ok := config["message"]; ok {
		d := dialog.Message(message.(string))

		if title, ok := config["title"]; ok {
			d.Title(title.(string))
		}

		if action, ok := config["type"]; ok {
			switch action {
			case "info":
				d.Info()
				return nil, nil
			case "error":
				d.Error()
				return nil, nil
			case "confirm":
				return d.YesNo(), nil
			}
			return nil, fmt.Errorf("unknown type")
		}
		return nil, fmt.Errorf("type is required")
	}
	return nil, fmt.Errorf("message is required")
}

func directory(config map[string]interface{}) (result interface{}, err error) {
	d := dialog.Directory()
	if title, ok := config["title"]; ok {
		d.Title(title.(string))
	}
	return d.Browse()
}

func file(config map[string]interface{}) (result interface{}, err error) {
	d := dialog.File()
	if title, ok := config["title"]; ok {
		d.Title(title.(string))
	}
	if startDir, ok := config["startDir"]; ok {
		d.SetStartDir(startDir.(string))
	}
	if filters, ok := config["filters"].(map[string]interface{}); ok {
		for desc, v0 := range filters {
			v1 := v0.([]interface{})
			extensions := make([]string, 0, len(v1))
			for _, extension := range v1 {
				extensions = append(extensions, extension.(string))
			}
			d.Filter(desc, extensions...)
		}
	}
	if action, ok := config["type"]; ok {
		switch action {
		case "load":
			return d.Load()
		case "save":
			return d.Save()
		}
	}
	return nil, fmt.Errorf("type is required")
}
