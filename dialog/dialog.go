package dialog

import (
	"fmt"
	"reflect"

	"github.com/sqweek/dialog"
	"github.com/zserge/webview"
)

// Init binds the js->go bridge for dialog functionality
func Init(w webview.WebView) {
	w.Bind("_dialogMux", mux)
	w.Init(_dialogJs)
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
