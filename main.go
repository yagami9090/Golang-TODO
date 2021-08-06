package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"todo/serialize"
	"todo/todo"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/auth", func(rw http.ResponseWriter, r *http.Request) {
		mySigningKey := []byte("password")
		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
			Issuer:    "test",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(mySigningKey)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(rw).Encode(map[string]string{
			"token": ss,
		})
	})

	api := r.NewRoute().Subrouter()
	api.Use(authMiddleware)

	todoApp := todo.NewApp(serialize.NewMsgPkgSerializer())

	r.HandleFunc("/todos", todoApp.AddTask).Methods(http.MethodPut)

	api.HandleFunc("/todos/{index}", todo.MarkDone).Methods(http.MethodPut)

	r.HandleFunc("/todos", todo.ListTask).Methods(http.MethodGet)

	http.ListenAndServe(":9090", r)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

		mySigningKey := []byte("password")
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return mySigningKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
