package main

import (
	"Twitter-Backend/data"
	"Twitter-Backend/handlers"
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	logger := log.New(os.Stdout, "[product-api]", log.LstdFlags)

	env := godotenv.Load()
	if env != nil {
		logger.Fatal(env)
	}

	port := os.Getenv("app_port")
	if len(port) == 0 {
		port = "8080"
	}

	tweetRepo, err := data.NewPostgreSql(logger)

	if err != nil {
		logger.Fatal(err)
	}

	tweetsHandler := handlers.New(logger, &tweetRepo)

	router := mux.NewRouter()
	router.Use(tweetsHandler.MiddlewareContentTypeSet)

	getAllRouter := router.Methods(http.MethodGet).Subrouter()
	getAllRouter.HandleFunc("/all", tweetsHandler.GetTweets)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/add", tweetsHandler.PostTweet)
	postRouter.Use(tweetsHandler.MiddlewareTweetValidation)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", tweetsHandler.PutTweet)
	putRouter.Use(tweetsHandler.MiddlewareTweetValidation)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", tweetsHandler.DeleteTweet)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGKILL)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)
	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)

	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}

}
