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

}

type seriesReq struct {
	Rooms []string
	From  time.Time
	To    time.Time
	By    storage.TimeBy
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("GetSeries(%s, %s, %v, %v)", req.From.Format(time.RFC3339), req.To.Format(time.RFC3339), req.By, req.Rooms)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("cannot encode entries: %v", err)
		return
	}
}
func (s *service) getRooms(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetRooms()")

	rooms, err := s.storage.GetRooms(r.Context(), time.Now().Add(-30*24*time.Hour))
	if err != nil {
		log.Printf("cannot get rooms: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		log.Printf("cannot encode rooms: %v", err)
		return
	}
}
