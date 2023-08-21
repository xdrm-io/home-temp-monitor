package main

import (
	"log"
	"os"
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

	storage, err := NewStorage("./db.sqlite")
	if err != nil {
		log.Fatalf("storage: %v", err)
	}
	defer storage.Close()

	collector, err := NewCollector(*cnf, storage)
	if err != nil {
		log.Fatalf("collector: %v", err)
	}
	defer collector.Close()

	if err := collector.Subscribe(); err != nil {
		log.Fatalf("cannot subscribe: %v", err)
	}

	for {
	}
}
