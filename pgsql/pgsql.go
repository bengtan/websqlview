package pgsql

import (
	"database/sql"
	_ "embed"
	"fmt"
	"reflect"
	"sync"

	"github.com/bengtan/websqlview/webviewex"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//go:embed pgsql.js
var _pgsqlJs string

var (
	databases    []*sql.DB
	transactions []*sql.Tx
	m            sync.Mutex
)

type queryable interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// Init binds the js->go bridge for sql functionality
func Init(ex *webviewex.WebViewEx) {
	// Pre-allocate table to hold up to 8 concurrent transactions
	transactions = make([]*sql.Tx, 0, 8)
	ex.W.Bind("_pgsqlMux", func(op string, args ...interface{}) (result interface{}, err error) {
		return mux(ex, op, args...)
	})
	ex.W.Init(_pgsqlJs)
}

// Shutdown should be called at program exit. Closes all databases.
func Shutdown() {
	for _, db := range databases {
		if db != nil {
			db.Close()
		}
	}
}

func mux(ex *webviewex.WebViewEx, op string, args ...interface{}) (result interface{}, err error) {
	if ex.URI[0:7] != "file://" && ex.URI[0:5] != "data:" {
		fmt.Println("pgsql.go ex.URI:", ex.URI)
		return nil, fmt.Errorf("access denied")
	}

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
			return query(false, int(handle), q, args[2:]...)
		}
	case "queryRow":
		handle, ok0 := args[0].(float64)
		q, ok1 := args[1].(string)
		if ok0 && ok1 {
			return query(true, int(handle), q, args[2:]...)
		}
	case "queryResult":
		handle, ok0 := args[0].(float64)
		q, ok1 := args[1].(string)
		if ok0 && ok1 {
			return queryResult(int(handle), q, args[2:]...)
		}
	case "begin":
		if handle, ok := args[0].(float64); ok {
			return begin(int(handle))
		}

	case "tx.commit":
		if handle, ok := args[0].(float64); ok {
			return nil, txCommitOrRollback(true, int(handle))
		}
	case "tx.rollback":
		if handle, ok := args[0].(float64); ok {
			return nil, txCommitOrRollback(false, int(handle))
		}
	case "tx.exec":
		handle, ok0 := args[0].(float64)
		q, ok1 := args[1].(string)
		if ok0 && ok1 {
			return txExec(int(handle), q, args[2:]...)
		}
	case "tx.query":
		handle, ok0 := args[0].(float64)
		q, ok1 := args[1].(string)
		if ok0 && ok1 {
			return txQuery(false, int(handle), q, args[2:]...)
		}
	case "tx.queryRow":
		handle, ok0 := args[0].(float64)
		q, ok1 := args[1].(string)
		if ok0 && ok1 {
			return txQuery(true, int(handle), q, args[2:]...)
		}
	case "tx.queryResult":
		handle, ok0 := args[0].(float64)
		q, ok1 := args[1].(string)
		if ok0 && ok1 {
			return txQueryResult(int(handle), q, args[2:]...)
		}
	}

	signature := []string{}
	for _, arg := range args {
		signature = append(signature, reflect.TypeOf(arg).Name())
	}
	return nil, fmt.Errorf("unknown operation %s with signature %v", op, signature)
}

func open(name string) (handle int, err error) {
	m.Lock()
	defer m.Unlock()

	// pgx requires a "postgres://" prefix
	sqlDb, err := sql.Open("pgx", "postgres://"+name)
	if err != nil {
		return -1, fmt.Errorf("open(%s): %s", name, err.Error())
	}

	err = sqlDb.Ping()
	if err != nil {
		sqlDb.Close()
		return -1, fmt.Errorf("open(%s): %s", name, err.Error())
	}

	handle = -1
	for i := range databases {
		if databases[i] == nil {
			// Reuse a handle
			handle = i
			databases[i] = sqlDb
			break
		}
	}

	if handle == -1 {
		// Use a new handle
		handle = len(databases)
		databases = append(databases, sqlDb)
	}
	return handle, nil
}

func close(handle int) (err error) {
	m.Lock()
	defer m.Unlock()

	if handle < 0 || handle >= len(databases) || databases[handle] == nil {
		return fmt.Errorf("invalid handle %d", handle)
	}
	db := databases[handle]
	databases[handle] = nil
	return db.Close()
}

func exec(handle int, q string, args ...interface{}) (result interface{}, err error) {
	if handle < 0 || handle >= len(databases) || databases[handle] == nil {
		return nil, fmt.Errorf("invalid handle %d", handle)
	}
	return _exec(databases[handle], q, args...)
}

func _exec(queryInterface queryable, q string, args ...interface{}) (result interface{}, err error) {
	code, err := queryInterface.Exec(q, args...)
	if err != nil {
		return nil, err
	}

	lastInsertID, _ := code.LastInsertId()
	rowsAffected, _ := code.RowsAffected()

	return map[string]interface{}{
		"lastInsertId": lastInsertID,
		"rowsAffected": rowsAffected,
	}, err
}

func query(singleton bool, handle int, q string, args ...interface{}) (result interface{}, err error) {
	if handle < 0 || handle >= len(databases) || databases[handle] == nil {
		return nil, fmt.Errorf("invalid handle %d", handle)
	}
	return _query(singleton, databases[handle], q, args...)
}

func _query(singleton bool, queryInterface queryable, q string, args ...interface{}) (result interface{}, err error) {
	rows, err := queryInterface.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Prepare placeholders for scanning
	types, _ := rows.ColumnTypes()
	columns := make([]sql.RawBytes, len(types))
	references := make([]interface{}, 0, len(types))
	for i := range types {
		references = append(references, &columns[i])
	}

	if singleton {
		if rows.Next() {
			err := rows.Scan(references...)
			if err != nil {
				return nil, err
			}
			object := map[string]interface{}{}
			for i, t := range types {
				object[t.Name()] = string(columns[i])
			}

			return object, rows.Err()
		}
		return nil, nil
	}

	data := make([]map[string]interface{}, 0, len(types))
	for rows.Next() {
		err := rows.Scan(references...)
		if err != nil {
			return data, err
		}

		object := map[string]interface{}{}
		for i, t := range types {
			object[t.Name()] = string(columns[i])
		}
		data = append(data, object)
	}

	return data, rows.Err()
}

func queryResult(handle int, q string, args ...interface{}) (result interface{}, err error) {
	if handle < 0 || handle >= len(databases) || databases[handle] == nil {
		return nil, fmt.Errorf("invalid handle %d", handle)
	}
	return _queryResult(databases[handle], q, args...)
}

func _queryResult(queryInterface queryable, q string, args ...interface{}) (result interface{}, err error) {
	var data interface{}
	err = queryInterface.QueryRow(q, args...).Scan(&data)

	if bytes, ok := data.([]byte); ok {
		return string(bytes), err
	}

	return data, err
}

func begin(handle int) (result interface{}, err error) {
	if handle < 0 || handle >= len(databases) || databases[handle] == nil {
		return nil, fmt.Errorf("invalid handle %d", handle)
	}

	tx, err := databases[handle].Begin()
	if err != nil {
		return -1, fmt.Errorf("begin(): %s", err.Error())
	}

	txHandle := -1
	for i := range transactions {
		if transactions[i] == nil {
			// Reuse a handle
			txHandle = i
			transactions[i] = tx
			break
		}
	}

	if txHandle == -1 {
		// Use a new handle
		txHandle = len(transactions)
		transactions = append(transactions, tx)
	}
	return txHandle, nil
}

func txCommitOrRollback(isCommit bool, handle int) (err error) {
	if handle < 0 || handle >= len(transactions) || transactions[handle] == nil {
		return fmt.Errorf("invalid handle %d", handle)
	}
	tx := transactions[handle]
	transactions[handle] = nil

	if isCommit {
		return tx.Commit()
	}
	return tx.Rollback()
}

func txExec(handle int, q string, args ...interface{}) (result interface{}, err error) {
	if handle < 0 || handle >= len(transactions) || transactions[handle] == nil {
		return nil, fmt.Errorf("invalid handle %d", handle)
	}

	return _exec(transactions[handle], q, args...)
}

func txQuery(singleton bool, handle int, q string, args ...interface{}) (result interface{}, err error) {
	if handle < 0 || handle >= len(transactions) || transactions[handle] == nil {
		return nil, fmt.Errorf("invalid handle %d", handle)
	}
	return _query(singleton, transactions[handle], q, args...)
}

func txQueryResult(handle int, q string, args ...interface{}) (result interface{}, err error) {
	if handle < 0 || handle >= len(transactions) || transactions[handle] == nil {
		return nil, fmt.Errorf("invalid handle %d", handle)
	}
	return _queryResult(transactions[handle], q, args...)
}
