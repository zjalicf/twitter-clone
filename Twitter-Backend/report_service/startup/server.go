package startup

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
	"github.com/zjalicf/twitter-clone-common/common/saga/messaging/nats"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
	"os/signal"
	"report_service/application"
	"report_service/domain"
	"report_service/handlers"
	"report_service/startup/config"
	store2 "report_service/store"
	"syscall"
	"time"
)

type Server struct {
	config *config.Config
}

const (
	QueueGroup = "report_service"
)

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

//
func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {

		}
	}(mongoClient, context.Background())

	cfg := config.NewConfig()

	ctx := context.Background()
	exp, err := newExporter(cfg.JaegerAddress)
	if err != nil {
		log.Fatalf("Failed to Initialize Exporter: %v", err)
	}

	tp := newTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("report_service")

	reportStore := server.initReportStore(mongoClient, tracer)

	reportService := server.initReportService(reportStore, tracer)

	replyPublisher := server.initPublisher(server.config.CreateReportReplySubject)
	commandSubscriber := server.initSubscriber(server.config.CreateReportCommandSubject, QueueGroup)

	server.initCreateEventHandler(reportService, replyPublisher, commandSubscriber)
	reportHandler := server.initAuthHandler(reportService, tracer)

	server.start(reportHandler)

}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := store2.GetMongoClient(server.config.ReportDBHost, server.config.ReportDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func newExporter(address string) (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("auth_service"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func (server *Server) initReportStore(client *mongo.Client, tracer trace.Tracer) domain.ReportStore {
	store := store2.NewReportMongoDBStore(client, tracer)
	return store
}

func (server *Server) initReportService(store domain.ReportStore, tracer trace.Tracer) *application.ReportService {
	return application.NewReportService(store, tracer)
}

func (server *Server) initAuthHandler(service *application.ReportService, tracer trace.Tracer) *handlers.ReportHandler {
	return handlers.NewReportHandler(service, tracer)
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

func (server *Server) initCreateEventHandler(reportService *application.ReportService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := handlers.NewCreateEventCommandHandler(reportService, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

// start
func (server *Server) start(authHandler *handlers.ReportHandler) {
	router := mux.NewRouter()
	authHandler.Init(router)

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

//
//	redisClient := server.initRedisClient()
//	authCache := server.initAuthCache(redisClient)
//	authStore := server.initAuthStore(mongoClient, tracer)
//
//	//saga init
//
//	//orchestrator
//	commandPublisher := server.initPublisher(server.config.CreateUserCommandSubject)
//	replySubscriber := server.initSubscriber(server.config.CreateUserReplySubject, QueueGroup)
//
//	//service
//	replyPublisher := server.initPublisher(server.config.CreateUserReplySubject)
//	commandSubscriber := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)
//
//	createUserOrchestrator := server.initCreateUserOrchestrator(commandPublisher, replySubscriber)
//
//	authService := server.initAuthService(authStore, authCache, createUserOrchestrator, tracer)
//
//	server.initCreateUserHandler(authService, replyPublisher, commandSubscriber)
//	authHandler := server.initAuthHandler(authService, tracer)
//
//	server.start(authHandler)
//}
//
//
//
//func (server *Server) initRedisClient() *redis.Client {
//	client, err := store2.GetRedisClient(server.config.AuthCacheHost, server.config.AuthCachePort)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return client
//}
//
//func (server *Server) initAuthStore(client *mongo.Client, tracer trace.Tracer) domain.AuthStore {
//	store := store2.NewAuthMongoDBStore(client, tracer)
//	return store
//}
//
//func (server *Server) initAuthCache(client *redis.Client) domain.AuthCache {
//	cache := store2.NewAuthRedisCache(client)
//	return cache
//}
//
//func (server *Server) initAuthService(store domain.AuthStore, cache domain.AuthCache, orchestrator *application.CreateUserOrchestrator, tracer trace.Tracer) *application.AuthService {
//	return application.NewAuthService(store, cache, orchestrator, tracer)
//}
//
//func (server *Server) initAuthHandler(service *application.AuthService, tracer trace.Tracer) *handlers.AuthHandler {
//	return handlers.NewAuthHandler(service, tracer)
//}
//
////saga
//
//func (server *Server) initPublisher(subject string) saga.Publisher {
//	publisher, err := nats.NewNATSPublisher(
//		server.config.NatsHost, server.config.NatsPort,
//		server.config.NatsUser, server.config.NatsPass, subject)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return publisher
//}
//
//func (server *Server) initSubscriber(subject string, queueGroup string) saga.Subscriber {
//	subscriber, err := nats.NewNATSSubscriber(
//		server.config.NatsHost, server.config.NatsPort,
//		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return subscriber
//}
//
//func (server *Server) initCreateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *application.CreateUserOrchestrator {
//	orchestrator, err := application.NewCreateUserOrchestrator(publisher, subscriber)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return orchestrator
//}
//
//func (server *Server) initCreateUserHandler(service *application.AuthService, publisher saga.Publisher, subscriber saga.Subscriber) {
//	_, err := handlers.NewCreateUserCommandHandler(service, publisher, subscriber)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//// start
//func (server *Server) start(authHandler *handlers.AuthHandler) {
//	router := mux.NewRouter()
//	authHandler.Init(router)
//
//	srv := &http.Server{
//		Addr:    fmt.Sprintf(":%s", server.config.Port),
//		Handler: router,
//	}
//
//	wait := time.Second * 15
//	go func() {
//		if err := srv.ListenAndServe(); err != nil {
//			log.Println(err)
//		}
//	}()
//
//	c := make(chan os.Signal, 1)
//
//	signal.Notify(c, os.Interrupt)
//	signal.Notify(c, syscall.SIGTERM)
//
//	<-c
//
//	ctx, cancel := context.WithTimeout(context.Background(), wait)
//	defer cancel()
//
//	if err := srv.Shutdown(ctx); err != nil {
//		log.Fatalf("Error Shutting Down Server %s", err)
//	}
//	log.Println("Server Gracefully Stopped")
//}
//
//func newExporter(address string) (*jaeger.Exporter, error) {
//	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
//	if err != nil {
//		return nil, err
//	}
//	return exp, nil
//}
//
//func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
//	r, err := resource.Merge(
//		resource.Default(),
//		resource.NewWithAttributes(
//			semconv.SchemaURL,
//			semconv.ServiceNameKey.String("auth_service"),
//		),
//	)
//
//	if err != nil {
//		panic(err)
//	}
//
//	return sdktrace.NewTracerProvider(
//		sdktrace.WithBatcher(exp),
//		sdktrace.WithResource(r),
//	)
//}
