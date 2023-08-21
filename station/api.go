package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type service struct {
	storage Storage
}

func NewServer(storage Storage) http.Handler {
	return &service{storage: storage}
}

// ServeHTTP implements http.Handler
func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// from to filters
	var (
		q = r.URL.Query()

		from = time.Time{}
		to   = time.Now().Add(time.Hour)
	)

	uFrom, err := strconv.ParseUint(q.Get("from"), 10, 64)
	if err == nil {
		from = time.Unix(int64(uFrom), 0)
	}

	uTo, err := strconv.ParseUint(q.Get("to"), 10, 64)
	if err == nil {
		to = time.Unix(int64(uTo), 0)
	}

	// rooms filter
	rooms := q["rooms"]
	if len(rooms) == 0 {
		s.getAll(r.Context(), w, from, to)
		return
	}

	s.getRooms(r.Context(), w, from, to, rooms) // FIXME: split rooms
}

func (s *service) getAll(ctx context.Context, w http.ResponseWriter, from, to time.Time) {
	log.Printf("GetAll(%d, %d)", from.Unix(), to.Unix())
	entries, err := s.storage.GetAll(ctx, from, to)
	if err != nil {
		log.Printf("cannot get entries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("cannot encode entries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *service) getRooms(ctx context.Context, w http.ResponseWriter, from, to time.Time, rooms []string) {
	log.Printf("GetRooms(%d, %d, %v)", from.Unix(), to.Unix(), rooms)
	entries, err := s.storage.GetRooms(ctx, from, to, rooms)
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
