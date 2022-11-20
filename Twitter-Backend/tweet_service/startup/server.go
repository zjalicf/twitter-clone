package startup

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tweet_service/application"
	"tweet_service/handlers"
	"tweet_service/startup/config"
	"tweet_service/store"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {

	tweetStore, err := store.New(log.Default())
	if err != nil {
		log.Fatal(err)
	}
	defer tweetStore.CloseSession()
	tweetStore.CreateTables()

	tweetService := server.initTweetService(*tweetStore)
	tweetHandler := server.initTweetHandler(tweetService)

	server.start(tweetHandler)
}

func (server *Server) initTweetService(store store.TweetRepo) *application.TweetService {
	return application.NewTweetService(&store)
}

func (server *Server) initTweetHandler(service *application.TweetService) *handlers.TweetHandler {
	return handlers.NewTweetHandler(service)
}

func (server *Server) start(tweetHandler *handlers.TweetHandler) {
	router := mux.NewRouter()
	tweetHandler.Init(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.config.Port),
		Handler: router,
	}

	wait := time.Second * 15
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error Shutting Down Server %s", err)
	}
	log.Println("Server Gracefully Stopped")
}
