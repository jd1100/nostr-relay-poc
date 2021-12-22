// +build lite
package main

import (
	//"github.com/jmoiron/sqlx"
	"github.com/timshannon/badgerhold/v4"
	//"github.com/dgraph-io/badger/v3"
	//_ "github.com/mattn/go-sqlite3"
)
var dbPath = "/tmp/db"
var options = badgerhold.DefaultOptions

func initDB() (*badgerhold.Store, error){
	options.Dir = "/tmp/db"
	options.ValueDir = "/tmp/db"
	db, err := badgerhold.Open(options)
	if err != nil {
		return nil, err
	}


	return db, nil
}
/*
func initDB() (*sqlx.DB, error) {
	db, err := badger.Open(badger.DefaultOptions(dbPath2))
	db, err := sqlx.Connect("sqlite3", s.SQLiteDatabase)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
CREATE TABLE event (
  id text NOT NULL,
  pubkey text NOT NULL,
  created_at integer NOT NULL,
  kind integer NOT NULL,
  tags text NOT NULL,
  content text NOT NULL,
  sig text NOT NULL
);

CREATE UNIQUE INDEX ididx ON event (id);
CREATE INDEX pubkeytimeidx ON event (pubkey, created_at);
    `)
	return db, nil
}

const relatedEventsCondition = `tags LIKE '%' || ? || '%'`
*/