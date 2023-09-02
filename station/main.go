package main

import (
	"log"
	"net/http"
	"os"

	"github.com/xdrm-io/home-temp-monitor/endpoint"
	"github.com/xdrm-io/home-temp-monitor/storage"
)

func env(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing mandatory environment variable %q", key)
	}
	return v
}

func main() {
	cnf, err := ReadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := storage.NewDB(cnf.DBPath)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}
	defer db.Close()

	// launch mqtt persistent collector
	collector, err := NewCollector(*cnf, db)
	if err != nil {
		log.Fatalf("collector: %v", err)
	}
	defer collector.Close()

	if err := collector.Subscribe(); err != nil {
		log.Fatalf("cannot subscribe: %v", err)
	}

	// setup http endpoint
	mux := http.NewServeMux()
	mux.Handle("/api/", endpoint.NewAPI(db))
	mux.Handle("/", endpoint.NewStaticSite())

	log.Printf("listening on %s", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("http server: %v", err)
	}
}
