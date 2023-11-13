package endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/xdrm-io/home-temp-monitor/storage"
)

type Service interface {
	Wire(r chi.Router)
}

type service struct {
	storage storage.Storage
}

func NewAPI(s storage.Storage) Service {
	return &service{storage: s}
}

func (s *service) Wire(r chi.Router) {
	r.MethodFunc("GET", "/series", s.getSeries)
	r.MethodFunc("GET", "/rooms", s.getRooms)
	r.MethodFunc("GET", "/last", s.getLast)

}

type seriesReq struct {
	Rooms   []string
	From    time.Time
	To      time.Time
	By      storage.TimeBy
	RefRoom *string
}

// Parse Request from http.Request
func (req *seriesReq) Parse(r *http.Request) error {
	req.From = time.Time{}
	req.To = time.Now()
	req.By = storage.ByHour
	req.Rooms = []string{}

	// from to filters
	q := r.URL.Query()

	uFrom, err := strconv.ParseUint(q.Get("from"), 10, 64)
	if err == nil {
		req.From = time.Unix(int64(uFrom), 0)
	}

	uTo, err := strconv.ParseUint(q.Get("to"), 10, 64)
	if err == nil {
		req.To = time.Unix(int64(uTo), 0)
	}

	sBy := storage.TimeBy(q.Get("by"))
	if err := sBy.OK(); err == nil {
		req.By = sBy
	}
	req.Rooms = q["rooms"]

	if q.Get("ref") != "" {
		ref := q.Get("ref")
		req.RefRoom = &ref
	}

	// fail if too much values are requested
	if req.To.Sub(req.From) > storage.MaxRows*req.By.Duration() {
		return fmt.Errorf("too much data requested : %d samples > %d", req.To.Sub(req.From)/req.By.Duration(), storage.MaxRows)
	}
	return nil
}

func (s *service) getSeries(w http.ResponseWriter, r *http.Request) {
	req := &seriesReq{}
	if err := req.Parse(r); err != nil {
		log.Printf("invalid request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("GetSeries(%s, %s, %v, %v)", req.From.Format(time.RFC3339), req.To.Format(time.RFC3339), req.By, req.Rooms)

	// ensure ref room is in rooms
	if req.RefRoom != nil {
		var found bool
		for _, r := range req.Rooms {
			if r == *req.RefRoom {
				found = true
				break
			}
		}
		if !found {
			req.Rooms = append(req.Rooms, *req.RefRoom)
		}
	}

	var (
		entries storage.Entries
		err     error
	)
	entries, err = s.storage.GetAll(r.Context(), req.From, req.To, req.By, req.Rooms)
	if err != nil {
		log.Printf("cannot get entries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// no reference room, return
	if req.RefRoom == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			log.Printf("cannot encode entries: %v", err)
		}
		return
	}

	// ref room, compute delta
	relative, err := relativeEntries(*req.RefRoom, req.By, entries)
	if err != nil {
		log.Printf("cannot compute relative entries: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(relative); err != nil {
		log.Printf("cannot encode entries: %v", err)
	}
}

func relativeEntries(refRoom string, by storage.TimeBy, entries storage.Entries) (storage.Entries, error) {
	// fetch refSeries room from entries
	refSeries, ok := entries[refRoom]
	if !ok {
		return nil, fmt.Errorf("ref room not found in entries")
	}

	// index @refSeries by timestamp (rounded to @by)
	timeBuckets := make(map[int64]storage.Entry, len(refSeries))
	for _, data := range refSeries {
		// round bucket with modSec
		bucket := time.Unix(data.Timestamp, 0).Round(by.Duration()).Unix()

		// store the first data for each bucket
		if _, ok := timeBuckets[bucket]; !ok {
			timeBuckets[data.Timestamp] = data
		}
	}

	relative := make(storage.Entries, len(entries))
	for room, series := range entries {
		if room == refRoom {
			continue
		}
		relative[room] = make([]storage.Entry, 0, len(series))
		for _, data := range series {
			bucket := time.Unix(data.Timestamp, 0).Round(by.Duration()).Unix()
			refBucket, ok := timeBuckets[bucket]
			if !ok {
				continue
			}
			relative[room] = append(relative[room], storage.Entry{
				Timestamp: data.Timestamp,
				TempMin:   data.TempMin - refBucket.TempMin,
				TempAvg:   data.TempAvg - refBucket.TempAvg,
				TempMax:   data.TempMax - refBucket.TempMax,
				HumMin:    data.HumMin - refBucket.HumMin,
				HumAvg:    data.HumAvg - refBucket.HumAvg,
				HumMax:    data.HumMax - refBucket.HumMax,
			})
		}
	}
	return relative, nil
}

func (s *service) getRooms(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetRooms()")

	rooms, err := s.storage.GetRooms(r.Context(), time.Now().Add(-30*24*time.Hour))
	if err != nil {
		log.Printf("cannot get rooms: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		log.Printf("cannot encode rooms: %v", err)
	}
}
func (s *service) getLast(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetLast()")

	last, err := s.storage.GetLast(r.Context())
	if err != nil {
		log.Printf("cannot get last: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(last); err != nil {
		log.Printf("cannot encode last: %v", err)
	}
}
