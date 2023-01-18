package startup

import (
	"context"
	"fmt"
	"follow_service/application"
	"follow_service/domain"
	"follow_service/handlers"
	"follow_service/startup/config"
	"follow_service/store"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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

func (server *Server) initNeo4JDriver() *neo4j.DriverWithContext {
	driver, err := store.GetClient(server.config.FollowDBHost, server.config.FollowDBPort,
		server.config.FollowDBUser, server.config.FollowDBPass)
	if err != nil {
		log.Fatal(err)
	}
	return driver
}

func (server *Server) initFollowStore(driver *neo4j.DriverWithContext) domain.FollowRequestStore {
	store := store.NewFollowNeo4JStore(driver)

	return store
}

func (server *Server) Start() {

	neo4jDriver := server.initNeo4JDriver()
	followStore := server.initFollowStore(neo4jDriver)
	followService := server.initFollowService(followStore)
	followHandler := server.initFollowHandler(followService)

	server.start(followHandler)
}

func (server *Server) initFollowService(store domain.FollowRequestStore) *application.FollowService {
	return application.NewFollowService(store)
}

func (server *Server) initFollowHandler(service *application.FollowService) *handlers.FollowHandler {
	return handlers.NewFollowHandler(service)
}

func (server *Server) start(followHandler *handlers.FollowHandler) {
	router := mux.NewRouter()
	followHandler.Init(router)

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
