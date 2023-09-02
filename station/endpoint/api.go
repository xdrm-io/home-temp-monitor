package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/xdrm-io/home-temp-monitor/storage"
)

type service struct {
	storage storage.Storage
}

func NewAPI(s storage.Storage) http.Handler {
	return &service{storage: s}
}

type request struct {
	Rooms []string
	From  time.Time
	To    time.Time
	By    storage.TimeBy
}

// Parse Request from http.Request
func (req *request) Parse(r *http.Request) error {
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

	// fail if too much values are requested
	if req.To.Sub(req.From) > storage.MaxRows*req.By.Duration() {
		return fmt.Errorf("too much data requested : %d samples > %d", req.To.Sub(req.From)/req.By.Duration(), storage.MaxRows)
	}
	return nil
}

// ServeHTTP implements http.Handler
func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req := &request{}
	if err := req.Parse(r); err != nil {
		log.Printf("invalid request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	s.getSeries(r.Context(), w, *req)
}

func (s *service) getSeries(ctx context.Context, w http.ResponseWriter, req request) {
	log.Printf("GetSeries(%s, %s, %v, %v)", req.From.Format(time.RFC3339), req.To.Format(time.RFC3339), req.By, req.Rooms)

	var (
		entries storage.Entries
		err     error
	)
	if len(req.Rooms) == 0 {
		entries, err = s.storage.GetAll(ctx, req.From, req.To, req.By)
	} else {
		entries, err = s.storage.GetRooms(ctx, req.From, req.To, req.By, req.Rooms)
	}

	if err != nil {
		log.Printf("cannot get entries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("cannot encode entries: %v", err)
		return
	}
}
