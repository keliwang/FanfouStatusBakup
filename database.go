package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

// Create a new sqlite database
func CreateFanfouDB(dbName string) error {
	// Remove old database at first
	os.Remove(dbName)

	// Open a new database
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	initDBSqls := []string{
		`CREATE TABLE STATUSES (
			status_num	INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at	TEXT NOT NULL,
			id		TEXT NOT NULL,
			rawid		INTEGER NOT NULL,
			text		TEXT NOT NULL,
			source		TEXT NOT NULL
		);`,
		`CREATE INDEX STATUS_NUM_INDEX ON STATUSES(status_num);`,
		`CREATE INDEX RAWID_INDEX ON STATUSES(rawid);`,
	}
	for _, sql := range initDBSqls {
		_, err = db.Exec(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Open our sqlite database and return a handle
func OpenDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Escape ' and " which are in the sql queries
func escapedQuotes(text string) string {
	escapedText := strings.Replace(text, "'", "''", -1)
	escapedText = strings.Replace(escapedText, `"`, `""`, -1)

	return escapedText
}

// Insert one status into our table
func InsertOneStatus(db *sql.DB, status Status) error {
	insertPattern := `INSERT INTO STATUSES(created_at, id, rawid, text, source) VALUES (
			"%s", "%s", %d, "%s", "%s");`
	insertSql := fmt.Sprintf(insertPattern, status.CreatedAt, status.Id,
		status.RawId, escapedQuotes(status.Text), escapedQuotes(status.Source))
	_, err := db.Exec(insertSql)
	if err != nil {
		return err
	}

	return nil
}

// Insert a slice of stauses into our table
func InsertStatuses(db *sql.DB, statuses *[]Status) error {
	for _, status := range *statuses {
		err := InsertOneStatus(db, status)
		if err != nil {
			return err
		}
	}

	return nil
}
