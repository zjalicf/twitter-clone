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
	"github.com/sirupsen/logrus"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
	"github.com/zjalicf/twitter-clone-common/common/saga/messaging/nats"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Logger = logrus.New()

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

const (
	QueueGroup = "follow_service"
)

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

func initLogger() {
	file, err := os.OpenFile("/app/logs/application.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(err)
	}

	Logger.SetOutput(file)

	rotationInterval := 24 * time.Hour
	ticker := time.NewTicker(rotationInterval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			rotateLogs(file)
		}
	}()
}

func rotateLogs(file *os.File) {
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	err := os.Rename("/app/logs/application.log", "/app/logs/application_"+currentTime+".log")
	if err != nil {
		Logger.Error(err)
	}
	file.Close()

	file, err = os.OpenFile("/app/logs/application.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		Logger.Error(err)
	}

	Logger.SetOutput(file)
}

func (server *Server) Start() {

	initLogger()

	neo4jDriver := server.initNeo4JDriver()
	followStore := server.initFollowStore(neo4jDriver)
	followService := server.initFollowService(followStore)
	followHandler := server.initFollowHandler(followService)

	//saga init
	replyPublisher := server.initPublisher(server.config.CreateUserReplySubject)
	commandSubscriber := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)

	server.initCreateUserHandler(followService, replyPublisher, commandSubscriber)

	server.start(followHandler)
}

func (server *Server) initFollowService(store domain.FollowRequestStore) *application.FollowService {
	return application.NewFollowService(store)
}

func (server *Server) initFollowHandler(service *application.FollowService) *handlers.FollowHandler {
	return handlers.NewFollowHandler(service)
}

func (server *Server) initCreateUserHandler(service *application.FollowService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := handlers.NewCreateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject string, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
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
