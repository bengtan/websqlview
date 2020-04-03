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
	w.Init(_sqliteJs)
}

// Shutdown should be called at program exit. Closes all database connections.
func Shutdown() {
	for _, db := range connections {
		if (db != nil) {
			db.Close()
		}
	}
}

func mux(op string, args ...interface{}) (result interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s: %v", op, e)
		}
	}()

	switch op {
	case "open":
		if name, ok := args[0].(string); ok {
			return open(name)
		}
	case "close":
		if handle, ok := args[0].(float64); ok {
			return nil, close(int(handle))
		}
	case "exec":
		handle, ok0 := args[0].(float64)
		q, ok1 := args[1].(string)
		if ok0 && ok1 {
			return exec(int(handle), q, args[2:]...)
		}
	case "query":
		handle, ok0 := args[0].(float64)
		q, ok1 := args[1].(string)
		if ok0 && ok1 {
			return query(int(handle), q, args[2:]...)
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

	err = db.Ping()
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

func exec(handle int, query string, args ...interface{}) (result interface{}, err error) {
	if (handle < 0 || handle >= len(connections) || connections[handle] == nil) {
		return nil, fmt.Errorf("Invalid handle %d", handle)
	}

	code, err := connections[handle].Exec(query, args...)
	if (err != nil) {
		return nil, err
	}

	lastInsertID, _ := code.LastInsertId()
	rowsAffected, _ := code.RowsAffected()

	return map[string]interface{}{
		"lastInsertId": lastInsertID,
		"rowsAffected": rowsAffected,
	}, err
}

func query(handle int, query string, args ...interface{}) (result interface{}, err error) {
	if (handle < 0 || handle >= len(connections) || connections[handle] == nil) {
		return nil, fmt.Errorf("Invalid handle %d", handle)
	}

	rows, err := connections[handle].Query(query, args...)
	if (err != nil) {
		return nil, err
	}
	defer rows.Close()

	// Prepare placeholders for scanning
	types, _ := rows.ColumnTypes()
	columns := make([]interface{}, len(types), len(types))
	references := make([]interface{}, 0, len(types))
	for i := range types {
		references = append(references, &columns[i])
	}

	data := make([]map[string]interface{}, 0, len(types))
	for rows.Next() {
		err := rows.Scan(references...)
		if err != nil {
			return data, err
		}

		object := map[string]interface{}{}
		for i, t := range types {
			object[t.Name()] = columns[i]
		}
		data = append(data, object)
	}

	err = rows.Err()
	return data, err
}
