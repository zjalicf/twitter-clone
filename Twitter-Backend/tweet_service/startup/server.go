package startup

import (
	"Twitter-Backend/application"
	"Twitter-Backend/domain"
	"Twitter-Backend/handlers"
	"Twitter-Backend/startup/config"
	store2 "Twitter-Backend/store"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	mongoClient := server.initMongoClient()
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {

		}
	}(mongoClient, context.Background())

	tweetStore := server.initTweetStore(mongoClient)
	tweetService := server.initTweetService(tweetStore)

	tweetHandler := server.initTweetHandler(tweetService)

	server.start(tweetHandler)

}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := store2.GetClient(server.config.TweetDBHost, server.config.TweetDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initTweetStore(client *mongo.Client) domain.TweetStore {
	store := store2.NewTweetMongoDBStore(client)

	//Delete everything from the database on server start
	//store.DeleteAll()
	return store
}

func (server *Server) initTweetService(store domain.TweetStore) *application.TweetService {
	return application.NewTweetService(store)
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
