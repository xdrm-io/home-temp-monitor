package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func NewDB(filename string) (Storage, error) {
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

	return &DB{db: db}, nil
}

func (s *DB) Close() error {
	return s.db.Close()
}

func (s *DB) Append(ctx context.Context, m Measure) error {
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

func (s *DB) GetAll(ctx context.Context, from, to time.Time, by TimeBy) (Entries, error) {
	var (
		grouping = "%Y-%m-%d %H:%M"
		suffix   = ":00"
	)
	switch by {
	case ByMin:
		grouping = "%Y-%m-%d %H:%M"
		suffix = ":00"
	case ByHour:
		grouping = "%Y-%m-%d %H"
		suffix = ":00:00"
	case ByDay:
		grouping = "%Y-%m-%d"
		suffix = " 00:00:00"
	case ByMonth:
		grouping = "%Y-%m"
		suffix = "-01 00:00:00"
	case ByYear:
		grouping = "%Y"
		suffix = "-01-01 00:00:00"
	}

	format := `
	SELECT
		room,
		unixepoch((strftime('%s', datetime(at, 'unixepoch')) || '%s')) as t,
		MIN(temperature) as tmin,
		AVG(temperature) as tavg,
		MAX(temperature) as tmax,
		MIN(humidity) as hmin,
		AVG(humidity) as havg,
		MAX(humidity) as hmax
	FROM environment
	WHERE at >= %d
	  AND at <= %d
	GROUP BY t, room
	ORDER BY room, at ASC
	LIMIT %d`

	query := fmt.Sprintf(
		format,
		grouping,
		suffix,
		from.Unix(),
		to.Unix(),
		MaxRows,
	)
	log.Printf("SQL query: %s", query)

	// get all measures for every room inside the given time range
	rows, err := s.db.QueryContext(ctx, query)
	if err == sql.ErrNoRows {
		log.Printf("sql: no rows")
		return nil, nil
	}
	if err != nil {
		log.Printf("sql: %v", err)
		return nil, err
	}
	defer rows.Close()

	entries := make(Entries, 0)
	for rows.Next() {
		var (
			room string
			e    Entry
		)
		err := rows.Scan(
			&room,
			&e.Timestamp,
			&e.TempMin,
			&e.TempAvg,
			&e.TempMax,
			&e.HumMin,
			&e.HumAvg,
			&e.HumMax,
		)
		if err != nil {
			return nil, err
		}
		entries[room] = append(entries[room], e)
	}
	return entries, nil
}

func (s *DB) GetRooms(ctx context.Context, from, to time.Time, by TimeBy, rooms []string) (Entries, error) {
	entries, err := s.GetAll(ctx, from, to, by)
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
