package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
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

		renderChart = q.Has("chart")
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
		s.getAll(r.Context(), w, renderChart, from, to)
		return
	}

	s.getRooms(r.Context(), w, renderChart, from, to, rooms) // FIXME: split rooms
}

func (s *service) getAll(ctx context.Context, w http.ResponseWriter, chart bool, from, to time.Time) {
	log.Printf("GetAll(%d, %d)", from.Unix(), to.Unix())
	entries, err := s.storage.GetAll(ctx, from, to)
	if err != nil {
		log.Printf("cannot get entries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if chart {
		s.renderChart(entries, w)
		return
	}

	// or json
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("cannot encode entries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *service) getRooms(ctx context.Context, w http.ResponseWriter, chart bool, from, to time.Time, rooms []string) {
	log.Printf("GetRooms(%d, %d, %v)", from.Unix(), to.Unix(), rooms)
	entries, err := s.storage.GetRooms(ctx, from, to, rooms)
	if err != nil {
		log.Printf("cannot get entries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if chart {
		s.renderChart(entries, w)
		return
	}

	// or json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("cannot encode entries: %v", err)
		return
	}
}

func (s service) entries2tseries(entries []Entry) []opts.LineData {
	series := make([]opts.LineData, 0, len(entries))
	for _, e := range entries {
		series = append(series, opts.LineData{
			// YAxisIndex: int(e.Temperature * 10),
			Value: []any{
				time.Unix(e.Timestamp, 0).Format(time.RFC3339),
				e.Temperature,
			},
			Symbol: "circle",
		})
	}
	return series
}
func (s service) entries2hseries(entries []Entry) []opts.LineData {
	series := make([]opts.LineData, 0, len(entries))
	for _, e := range entries {
		series = append(series, opts.LineData{
			// YAxisIndex: int(e.Temperature * 10),
			Value: []any{
				time.Unix(e.Timestamp, 0).Format(time.RFC3339),
				e.Humidity,
			},
			Symbol: "circle",
		})
	}
	return series
}

func (s *service) renderChart(entries Entries, w http.ResponseWriter) {
	globalOpts := []charts.GlobalOpts{
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),

		charts.WithYAxisOpts(opts.YAxis{
			Name: "Temperature",
			Type: "value",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time",
			Type: "time",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "slider"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
	}
	seriesOpts := []charts.SeriesOpts{
		charts.WithLineChartOpts(opts.LineChart{
			Smooth:     false,
			ShowSymbol: true,
		}),
	}

	tChart := charts.NewLine()
	tChart.SetGlobalOptions(append(globalOpts, charts.WithTitleOpts(opts.Title{
		Title: "Temperature (°C) by room over time",
	}))...)
	tChart.SetSeriesOptions(seriesOpts...)
	for room, data := range entries {
		tChart.AddSeries(room, s.entries2tseries(data))
	}

	hChart := charts.NewLine()
	hChart.SetGlobalOptions(append(globalOpts, charts.WithTitleOpts(opts.Title{
		Title: "Humidity (%) by room over time",
	}))...)
	hChart.SetSeriesOptions(seriesOpts...)
	for room, data := range entries {
		hChart.AddSeries(room, s.entries2hseries(data))
	}

	tChart.Render(w)
	hChart.Render(w)
}
