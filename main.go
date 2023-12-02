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
	r.Get("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	r.Mount("/bar", exampleSubRouter())
	http.ListenAndServe("localhost:3000", r)
}

func exampleSubRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(ExampleMiddleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	return r
}

func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello, from middleware!")
		next.ServeHTTP(w, r)
	})
}
