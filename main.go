package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo/app"
	"todo/repository"
	"todo/todo"
)

func main() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/auth", func(c *gin.Context) {
		mySigningKey := []byte("password")
		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
			Issuer:    "test",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(mySigningKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, map[string]string{
			"token": ss,
		})
	})

	api := r.Group("/")
	api.Use(authMiddleware)

	todoSrv := todo.NewService(repository.NewGormRepository(db))
	todoApp := todo.New(todoSrv)

	app := app.New(r)

	app.PUT("/todos", todoApp.Add)
	api.PUT("/todos/:id", todoApp.MarkDone)
	api.GET("/todos", todoApp.ListTask)

	r.Run(":9090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// func xmain() {
// 	r := mux.NewRouter()

// 	r.HandleFunc("/auth", func(rw http.ResponseWriter, r *http.Request) {
// 		mySigningKey := []byte("password")
// 		claims := &jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
// 			Issuer:    "test",
// 		}

// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 		ss, err := token.SignedString(mySigningKey)
// 		if err != nil {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		json.NewEncoder(rw).Encode(map[string]string{
// 			"token": ss,
// 		})
// 	})

// 	api := r.NewRoute().Subrouter()
// 	api.Use(authMiddleware)

// 	api.HandleFunc("/todos", todo.AddTask).Methods(http.MethodPut)

// 	api.HandleFunc("/todos/{index}", todo.MarkDone).Methods(http.MethodPut)

// 	api.HandleFunc("/todos", todo.ListTask).Methods(http.MethodGet)

// 	http.ListenAndServe(":9090", r)
// }

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.RequestURI)
// 		next.ServeHTTP(w, r)
// 	})
// }

func authMiddleware(c *gin.Context) {
	fmt.Println("start auth mw")
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	mySigningKey := []byte("password")
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})

	if err != nil {
		fmt.Println("err-->", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

// func authMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")
// 		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

// 		mySigningKey := []byte("password")
// 		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 			}

// 			return mySigningKey, nil
// 		})

// 		if err != nil {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
