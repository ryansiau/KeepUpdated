package main

import (
	"github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
)

func NewSqliteDB(uri string) (*gorm.DB, error) {
	conn, err := driver.Open(uri, nil) // sample conn: "file:demo.db"
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(gormlite.OpenDB(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
