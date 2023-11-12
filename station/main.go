package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Expose-Headers", "Link")
			next.ServeHTTP(w, r)
		})
	})

	var (
		website = endpoint.NewStaticWeb()
		api     = endpoint.NewAPI(db)
	)
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		api.Wire(r)
	})
	router.Method("GET", "/*", website)

	log.Printf("listening on %s", ":8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("http server: %v", err)
	}
}
