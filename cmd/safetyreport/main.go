package main

import (
	"log"
	"net/http"
	"os"
	"timber-safety/internal/api"
	"timber-safety/internal/service"
	"timber-safety/internal/store"
)

func main() {
	path := os.Getenv("TIMBER_DB")
	if path == "" {
		path = "./data/timber.db"
	}
	s, err := store.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	svc := service.New(s)
	srv := api.New(svc)
	addr := os.Getenv("TIMBER_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
