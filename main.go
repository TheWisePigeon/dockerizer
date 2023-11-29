package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello bozo"))
		return
	})
	err := http.ListenAndServe(":6061", r)
	if err != nil {
		panic(err)
	}
}
