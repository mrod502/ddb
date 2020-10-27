package ddb

import "github.com/dgraph-io/badger/v2"

var db *badger.DB

func OpenDB() (err error) {
	db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	return
}
