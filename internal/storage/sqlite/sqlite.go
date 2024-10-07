package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}
type ApiKey struct {
	Id     int64  `db:"id" json:"id"`
	ApiKey string `db:"api_key" json:"api_key"`
	Name   string `db:"name" json:"name"`
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening SQLite database: %v", err)
	}

	statement, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS api_keys(
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    api_key TEXT NOT NULL UNIQUE,
	    name TEXT NOT NULL
	);
`)

	if err != nil {
		return nil, fmt.Errorf("error preparing SQLite table: %v", err)
	}

	if _, err := statement.Exec(); err != nil {
		return nil, fmt.Errorf("error executing SQLite table: %v", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddApiKey(apiKey string, name string) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO api_keys(api_key, name) VALUES (?, ?)")
	if err != nil {
		println(err.Error())
	}
	result, err := stmt.Exec(apiKey, name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Storage) GetApiKeys() ([]*ApiKey, error) {
	keys := make([]*ApiKey, 0)
	stmt, err := s.db.Prepare("SELECT id, api_key, name FROM api_keys LIMIT 20")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		apiKey := new(ApiKey)
		if err := rows.Scan(&apiKey.Id, &apiKey.ApiKey, &apiKey.Name); err != nil {
			return nil, err
		}
		keys = append(keys, apiKey)
	}
	return keys, nil
}

func (s *Storage) DeleteApiKey(id int64) error {
	stmt, err := s.db.Prepare("DELETE FROM api_keys WHERE id = ?")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(id); err != nil {
		return err
	}
	return nil
}
