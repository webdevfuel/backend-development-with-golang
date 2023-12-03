package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Session struct {
	Key   string
	Value any
}

var (
	store map[string]Session = make(map[string]Session)
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/foo", func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Key")
		if key == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		store[key] = Session{
			Key:   key,
			Value: "Hello, from within an example session!",
		}
		cookie := http.Cookie{
			Name:  "_example_session",
			Value: key,
		}
		http.SetCookie(w, &cookie)
		w.Write([]byte("Hello, world!"))
	})
	r.Get("/bar", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("_example_session")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		key := cookie.Value
		session, ok := store[key]
		if !ok {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Println(session.Value)
		w.Write([]byte("Hello, world!"))
	})
	http.ListenAndServe("localhost:3000", r)
}
