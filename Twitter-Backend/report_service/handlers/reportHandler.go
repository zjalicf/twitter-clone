package handlers

import (
	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"report_service/application"
	"report_service/authorization"
	"report_service/store"
)

type KeyUser struct{}

type ReportHandler struct {
	service *application.ReportService
	store   *store.ReportMongoDBStore
	tracer  trace.Tracer
}

func NewReportHandler(service *application.ReportService, tracer trace.Tracer) *ReportHandler {
	return &ReportHandler{
		service: service,
		tracer:  tracer,
	}
}

func (handler *ReportHandler) Init(router *mux.Router) {
	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	log.Println("successful init of enforcer")
	if err != nil {
		log.Fatal(err)
	}

	//router.HandleFunc("/login", handler.Login).Methods("POST")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8003", authorization.Authorizer(authEnforcer)(router)))

}
