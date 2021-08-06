package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo/app"
	"todo/repository"
	"todo/todo"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("error:", err)
	}

	viper.SetDefault("app.port", ":9090")
}

func main() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	logger, _ := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
	}.Build()

	defer logger.Sync()

	r := app.New()
	r.GET("/auth", func(c *app.Ctx) {
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

	r.PUT("/todos", todoApp.Add)
	r.PUT("/todos/:id", todoApp.MarkDone)
	r.GET("/todos", todoApp.ListTask)

	srv := &http.Server{
		Addr:    viper.GetString("app.port"),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
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
