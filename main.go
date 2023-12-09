package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	key string = "e7c18e1c80da49adab37f45f67da726e"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Delete("/foo", func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:   "message",
			MaxAge: -1,
		}
		http.SetCookie(w, &cookie)
	})
	r.Post("/foo", func(w http.ResponseWriter, r *http.Request) {
		message, err := encryptMessage("Hello, encrypted from cookies!")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:  "message",
			Value: message,
		}
		http.SetCookie(w, &cookie)
	})
	r.Get("/bar", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("message")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		message, err := decryptMessage(cookie.Value)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Write([]byte(message))
	})
	http.ListenAndServe("localhost:3000", r)
}

func encryptMessage(message string) (string, error) {
	plaintext := []byte(message)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	encryptedMessage := hex.EncodeToString(ciphertext)
	return encryptedMessage, nil
}

func decryptMessage(message string) (string, error) {
	c, _ := hex.DecodeString(message)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := c[:nonceSize], c[nonceSize:]
	m, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(m), nil
}
