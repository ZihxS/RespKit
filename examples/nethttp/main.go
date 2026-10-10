package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/ZihxS/RespKit/resp"
)

func main() {
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		if err := resp.OK(resp.HTTP(w, r), map[string]any{"user_id": 7, "full_name": "Saleh"}); err != nil {
			log.Print(err)
		}
	})
	http.HandleFunc("/not-found", func(w http.ResponseWriter, r *http.Request) {
		if err := resp.Fail(resp.HTTP(w, r), resp.ErrNotFound); err != nil {
			log.Print(err)
		}
	})
	http.HandleFunc("/internal", func(w http.ResponseWriter, r *http.Request) {
		if err := resp.Fail(resp.HTTP(w, r), errors.New("connection string: DO NOT LEAK")); err != nil {
			log.Print(err)
		}
	})
	addr := os.Getenv("RESPKIT_ADDR")
	if addr == "" {
		addr = ":18090"
	}
	log.Printf("Try curl -i -H 'Accept-Language: en-US' localhost%s/not-found", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
