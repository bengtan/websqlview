package webviewex

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/webview/webview"
)

// WebViewEx is an extension/wrapper of WebView
type WebViewEx struct {
	W     webview.WebView
	token string
	URI   string
}

// New creates a WebViewEx extension/wrapper around the supplied WebView
func New(w webview.WebView) *WebViewEx {
	var token string

	data := make([]byte, 16)
	_, err := rand.Read(data)
	if err == nil {
		token = base64.StdEncoding.EncodeToString(data)
	} else {
		fmt.Println("Error: ", err)
		// Fallback token in case of error
		token = "0123456789abcdef"
	}

	ex := &WebViewEx{w, token, "blah"}

	// This is a round-about way to track the current uri of the webview
	w.Bind("__updateLocation", func(uri string, token string) {
		updateLocation(ex, uri, token)
	})
	w.Init(fmt.Sprintf("__updateLocation(window.location.toString(), '%s')", token))

	return ex
}

func updateLocation(ex *WebViewEx, uri string, token string) {
	if token == ex.token {
		ex.URI = uri
	} else {
		// This should never happen. If it happens, maybe an attacker is trying
		// to spoof the updateLocation() call? Hence, abort immediately
		fmt.Println("Error: updateLocation: invalid token. Aborting")
		ex.W.Terminate()
	}
}
