package database

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

type DB struct {
    conn *sql.DB
}

func Connect(connectionString string) (*DB, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    fmt.Println("Connected to the database")
    return &DB{conn: db}, nil
}

func (db *DB) Close() error {
    return db.conn.Close()
}

func (db *DB) InsertID(id int, tableName string) error {
    insertStatement := fmt.Sprintf("INSERT INTO %s (id) VALUES ($1)", tableName)
    _, err := db.conn.Exec(insertStatement, id)
    if err != nil {
        return err
    }

    fmt.Println("ID inserted successfully")
    return nil
}
