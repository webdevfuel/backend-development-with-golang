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
	r.Delete("/foo", func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:   "message",
			Value:  "Hello, from cookies!",
			MaxAge: -1,
		}
		http.SetCookie(w, &cookie)
	})
	r.Post("/foo", func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:  "message",
			Value: "Hello, from cookies!",
		}
		http.SetCookie(w, &cookie)
	})
	r.Get("/bar", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("message")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Println(cookie.Value)
	})
	http.ListenAndServe("localhost:3000", r)
}
