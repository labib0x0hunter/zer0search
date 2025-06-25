package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// CREATE DATABASE testDb;
// SHOW DATABASES;
// USE testDb;

/**
CREATE TABLE items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL
);
**/

var (
	user       = "root"
	password   = ""
	ip         = "127.0.0.1"
	port       = "3306"
	insertStmt = `INSERT INTO items (content) VALUES (?)`
	queryStmt  = `SELECT content FROM items WHERE id = ?`

	NoEntryError = "sql: Scan error on column index 0"
)

type DbService struct {
	db *sql.DB
}

func NewDbService() (*DbService, error) {
	dbName := "testDb"
	dsn := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DbService{
		db: db,
	}, nil
}

func (d *DbService) Insert(document string) error {
	_, err := d.db.Exec(insertStmt, document)
	return err
}

func (d *DbService) Query(id int) (string, error) {
	var document string
	if err := d.db.QueryRow(queryStmt, id).Scan(&document); err != nil {
		return "", err
	}
	return document, nil
}

func (d *DbService) GetDocId() (int, error) {
	var lastID int
	row := d.db.QueryRow(`SELECT MAX(id) FROM items`)
	err := row.Scan(&lastID)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

func (d *DbService) Close() error {
	return d.db.Close()
}
