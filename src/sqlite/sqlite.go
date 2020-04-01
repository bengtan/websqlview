package sqlite

import (
	"database/sql"
	"fmt"
	"reflect"
	// SQLite database driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/zserge/webview"
)

var connections []*sql.DB

// Init binds the js->go bridge for sqlite functionality
func Init(w webview.WebView) {
	w.Bind("_sqliteMux", mux)
}

// Shutdown should be called at program exit. Closes all database connections.
func Shutdown() {
	for _, db := range connections {
		db.Close()
	}
}

func mux(op string, args ...interface{}) (result interface{}, err error) {
	switch op {
	case "open":
		if a0, ok := args[0].(string); ok {
			return open(a0)
		}
	case "close":
		if a0, ok := args[0].(float64); ok {
			return nil, close(int(a0))
		}
	}

	signature := []string{}
	for _, arg := range args {
		signature = append(signature, reflect.TypeOf(arg).Name())
	}
	return nil, fmt.Errorf("Unknown operation %s with signature %v", op, signature)
}

func open(name string) (result interface{}, err error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return -1, fmt.Errorf("open(%s): %s", name, err.Error())
	}

	handle := len(connections)
	connections = append(connections, db)
	return handle, nil
}

func close(handle int) (err error) {
	if (handle < 0 || handle >= len(connections) || connections[handle] == nil) {
		return fmt.Errorf("Invalid handle %d", handle)
	}
	db := connections[handle]
	connections[handle] = nil
	return db.Close()
}
