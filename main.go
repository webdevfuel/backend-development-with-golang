package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/foo", func(w http.ResponseWriter, r *http.Request) {
		message := r.FormValue("message")
		if message == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Println(message)
		w.Write([]byte("Hello, world!"))
	})
	r.Post("/bar", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		var messages []string
		for k, v := range r.Form {
			if k == "message" {
				messages = append(messages, v...)
			}
		}
		for _, message := range messages {
			log.Println(message)
		}
		w.Write([]byte("Hello, world!"))
	})
	http.ListenAndServe("localhost:3000", r)
}
