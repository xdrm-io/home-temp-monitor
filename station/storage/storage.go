package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

const MaxRows = 5000

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
	Timestamp int64   `json:"t"`
	TempMin   float32 `json:"tmin"`
	TempAvg   float32 `json:"tavg"`
	TempMax   float32 `json:"tmax"`
	HumMin    float32 `json:"hmin"`
	HumAvg    float32 `json:"havg"`
	HumMax    float32 `json:"hmax"`
}

type TimeBy string

// avail time grouping
const (
	ByMin   = TimeBy("minute")
	ByHour  = TimeBy("hour")
	ByDay   = TimeBy("day")
	ByMonth = TimeBy("month") // 30 days
	ByYear  = TimeBy("year")
)

func (t TimeBy) OK() error {
	switch t {
	case ByMin, ByHour, ByDay, ByMonth, ByYear:
		return nil
	}
	return fmt.Errorf("unknown time grouping %q ; [minute, hour, day, month, year]", t)
}

func (t TimeBy) Duration() time.Duration {
	switch t {
	case ByMin:
		return time.Minute
	case ByHour:
		return time.Hour
	case ByDay:
		return 24 * time.Hour
	case ByMonth:
		return 30 * 24 * time.Hour
	case ByYear:
		return 365 * 24 * time.Hour
	}
	log.Printf("unknown time grouping: %q ; using minute", t)
	return time.Minute
}

type Storage interface {
	GetAll(ctx context.Context, from, to time.Time, by TimeBy) (Entries, error)
	GetRooms(ctx context.Context, from, to time.Time, by TimeBy, rooms []string) (Entries, error)

	Append(ctx context.Context, m Measure) error
	io.Closer
}
