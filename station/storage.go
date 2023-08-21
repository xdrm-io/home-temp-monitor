package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Measure struct {
	Room        string
	Temperature uint16 `json:"t"`
	Humidity    uint16 `json:"h"`
	OffsetSec   uint32 `json:"d"`
}

type Storage interface {
	Append(ctx context.Context, m Measure) error
	io.Closer
}

type storage struct {
	db *sql.DB
}

func NewStorage(filename string) (Storage, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS environment (
		uid         INTEGER PRIMARY KEY AUTOINCREMENT,
		at          TIMESTAMP NOT NULL,
		room        TEXT NOT NULL,
		temperature REAL NOT NULL,
		humidity    REAL NOT NULL
	);`)
	if err != nil {
		return nil, fmt.Errorf("cannot create table: %w", err)
	}

	return &storage{db: db}, nil
}

func (s *storage) Close() error {
	return s.db.Close()
}

func (s *storage) Append(ctx context.Context, m Measure) error {
	var (
		ts = time.Now().Add(time.Duration(m.OffsetSec) * time.Second).Unix()
		t  = float64(m.Temperature) / 10.
		h  = float64(m.Humidity) / 10.
	)
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO environment (at, room, temperature, humidity) VALUES (?, ?, ?, ?) `,
		ts,
		m.Room,
		t,
		h,
	)
	return err
}
