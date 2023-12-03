package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type messageKey struct{}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(ExampleMiddleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		message, ok := ctx.Value(messageKey{}).(string)
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Println(message)
		w.Write([]byte("Hello, world!"))
	})
	http.ListenAndServe("localhost:3000", r)
}

func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message := "Hello, with data from context!"
		ctx := context.WithValue(r.Context(), messageKey{}, message)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
