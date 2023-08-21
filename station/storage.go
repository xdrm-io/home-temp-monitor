package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	_ "modernc.org/sqlite"
)

// Measure as received from the sensor
type Measure struct {
	Room        string
	Temperature uint16 `json:"t"`
	Humidity    uint16 `json:"h"`
	OffsetSec   uint32 `json:"d"`
}

// Entries indexed by room id
type Entries map[string][]Entry

// Entry as returned by the storage GetXxx() methods
type Entry struct {
	Timestamp   int64   `json:"ts"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

type Storage interface {
	GetAll(ctx context.Context, from, to time.Time) (Entries, error)
	GetRooms(ctx context.Context, from, to time.Time, rooms []string) (Entries, error)

	Append(ctx context.Context, m Measure) error
	io.Closer
}

type storage struct {
	db *sql.DB
}

func NewStorage(filename string) (Storage, error) {
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS environment (
		uid         INTEGER PRIMARY KEY AUTOINCREMENT,
		at          INTEGER NOT NULL,
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
		ts = time.Now().Add(-time.Duration(m.OffsetSec) * time.Second).Unix()
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

func (s *storage) GetAll(ctx context.Context, from, to time.Time) (Entries, error) {
	// get all measures for every room inside the given time range
	rows, err := s.db.QueryContext(ctx, `
		SELECT room, at, temperature, humidity FROM environment
		WHERE at >= ? AND at <= ?
		ORDER BY ROOM, at ASC
	`, from.Unix(), to.Unix())
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make(Entries, 0)
	for rows.Next() {
		var (
			r string
			e Entry
		)
		if err := rows.Scan(&r, &e.Timestamp, &e.Temperature, &e.Humidity); err != nil {
			return nil, err
		}
		entries[r] = append(entries[r], e)
	}
	return entries, nil
}

func (s *storage) GetRooms(ctx context.Context, from, to time.Time, rooms []string) (Entries, error) {
	entries, err := s.GetAll(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// filter entries by room
	filtered := make(Entries, 0)
	for room, e := range entries {
		var found bool
		for _, r := range rooms {
			if room == r {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		filtered[room] = e
	}
	return filtered, nil
}
