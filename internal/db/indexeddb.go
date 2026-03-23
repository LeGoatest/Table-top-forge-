package db

import (
	"errors"
	"fmt"
	"syscall/js"
)

// Database represents a wrapper around IndexedDB.
type Database struct {
	db js.Value
}

// NewDatabase initializes and opens the IndexedDB database.
func NewDatabase(dbName string) (*Database, error) {
	if js.Global().Get("indexedDB").IsUndefined() {
		return nil, errors.New("indexedDB not supported")
	}

	resCh := make(chan struct {
		db  js.Value
		err error
	})

	request := js.Global().Get("indexedDB").Call("open", dbName, 1)

	request.Set("onupgradeneeded", js.FuncOf(func(this js.Value, args []js.Value) any {
		db := args[0].Get("target").Get("result")
		if !db.Get("objectStoreNames").Call("contains", "games").Bool() {
			db.Call("createObjectStore", "games")
		}
		return nil
	}))

	request.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) any {
		db := args[0].Get("target").Get("result")
		resCh <- struct {
			db  js.Value
			err error
		}{db: db, err: nil}
		return nil
	}))

	request.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
		resCh <- struct {
			db  js.Value
			err error
		}{db: js.Null(), err: fmt.Errorf("failed to open indexedDB: %s", args[0].Get("target").Get("error").Get("message").String())}
		return nil
	}))

	res := <-resCh
	if res.err != nil {
		return nil, res.err
	}

	return &Database{db: res.db}, nil
}

// Save stores a value in the 'games' object store.
func (d *Database) Save(key string, value any) error {
	resCh := make(chan error)

	transaction := d.db.Call("transaction", []any{"games"}, "readwrite")
	store := transaction.Call("objectStore", "games")
	request := store.Call("put", js.ValueOf(value), js.ValueOf(key))

	request.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) any {
		resCh <- nil
		return nil
	}))

	request.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
		resCh <- fmt.Errorf("failed to save: %s", args[0].Get("target").Get("error").Get("message").String())
		return nil
	}))

	return <-resCh
}

// Load retrieves a value from the 'games' object store.
func (d *Database) Load(key string) (js.Value, error) {
	resCh := make(chan struct {
		val js.Value
		err error
	})

	transaction := d.db.Call("transaction", []any{"games"}, "readonly")
	store := transaction.Call("objectStore", "games")
	request := store.Call("get", js.ValueOf(key))

	request.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) any {
		val := args[0].Get("target").Get("result")
		resCh <- struct {
			val js.Value
			err error
		}{val: val, err: nil}
		return nil
	}))

	request.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
		resCh <- struct {
			val js.Value
			err error
		}{val: js.Undefined(), err: fmt.Errorf("failed to load: %s", args[0].Get("target").Get("error").Get("message").String())}
		return nil
	}))

	res := <-resCh
	if res.err != nil {
		return js.Undefined(), res.err
	}

	return res.val, nil
}
