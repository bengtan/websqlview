package silk

import (
	"fmt"
	"github.com/zserge/webview"
)

// Init binds some custom functionality
func Init(w webview.WebView, exitCode *int) {
	w.Bind("exit", func(w webview.WebView, i int) (err error) {
		if w.GetURI()[0:7] != "file://" {
			return fmt.Errorf("Access denied")
		}
		w.Terminate()
		*exitCode = i
		return nil
	})
}

