package handlers

import (
	"encoding/json"
	"fmt"
	"follow_service/application"
	"follow_service/authorization"
	"follow_service/domain"
	"follow_service/errors"
	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

type FollowHandler struct {
	service            *application.FollowService
	counterUnavailable int
	tracer             trace.Tracer
}

func NewFollowHandler(service *application.FollowService, tracer trace.Tracer) *FollowHandler {
	return &FollowHandler{
		service:            service,
		counterUnavailable: 3,
		tracer:             tracer,
	}
}

func (handler *FollowHandler) Init(router *mux.Router) {

	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	log.Println("successful init of enforcer")
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/requests/", handler.GetRequestsForUser).Methods("GET")
	router.HandleFunc("/requests/{visibility}", handler.CreateRequest).Methods("POST")
	router.HandleFunc("/acceptRequest/{id}", handler.AcceptRequest).Methods("PUT")
	router.HandleFunc("/declineRequest/{id}", handler.DeclineRequest).Methods("PUT")
	router.HandleFunc("/followings", handler.GetFollowingsByUser).Methods("GET")
	router.HandleFunc("/recommendations", handler.GetRecommendationsForUser).Methods("GET")
	router.HandleFunc("/ad", handler.SaveAd).Methods("POST")

	http.Handle("/", router)
	log.Println("Successful")
	log.Fatal(http.ListenAndServe(":8004", authorization.Authorizer(authEnforcer)(router)))
}

func (handler *FollowHandler) GetRequestsForUser(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "FollowHandler.GetRequestsForUser")
	defer span.End()

	token, err := authorization.GetToken(req)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := authorization.GetMapClaims(token.Bytes())

	returnRequests, err := handler.service.GetRequestsForUser(ctx, claims["username"])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(returnRequests)

	jsonResponse(returnRequests, writer)
}

func (handler *FollowHandler) GetFollowingsByUser(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "FollowHandler.GetFollowingsByUser")
	defer span.End()

	token, _ := authorization.GetToken(req)
	claims := authorization.GetMapClaims(token.Bytes())
	username := claims["username"]

	log.Printf("username is: %s", username)
	users, err := handler.service.GetFollowingsOfUser(ctx, username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(users, writer)

}

func (handler *FollowHandler) GetRecommendationsForUser(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "FollowHandler.GetRecommendationsForUser")
	defer span.End()

	token, _ := authorization.GetToken(req)
	claims := authorization.GetMapClaims(token.Bytes())
	username := claims["username"]

	users, err := handler.service.GetRecommendationsByUsername(ctx, username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(users, writer)
}

func (handler *FollowHandler) CreateRequest(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "FollowHandler.CreateRequest")
	defer span.End()

	var request domain.FollowRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Header.Get("Authorization") == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := authorization.GetToken(req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
	}
	claims := authorization.GetMapClaims(token.Bytes())

	vars := mux.Vars(req)
	fmt.Println(vars)
	var visibility bool
	if vars["visibility"] == "private" {
		visibility = true
	} else {
		visibility = false
	}

	err = handler.service.CreateRequest(ctx, &request, claims["username"], visibility)
	if err != nil {
		log.Println("ERR")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *FollowHandler) AcceptRequest(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "FollowHandler.AcceptRequest")
	defer span.End()

	vars := mux.Vars(req)
	followId, ok := vars["id"]
	if !ok {
		http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	}

	err := handler.service.AcceptRequest(ctx, &followId)
	if err != nil {
		http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	}

	writer.WriteHeader(http.StatusOK)

}

func (handler *FollowHandler) DeclineRequest(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "FollowHandler.GetAll")
	defer span.End()

	vars := mux.Vars(req)
	followId, ok := vars["id"]
	if !ok {
		http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	}

	err := handler.service.DeclineRequest(ctx, &followId)
	if err != nil {
		http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	}

	writer.WriteHeader(http.StatusOK)

}

func (handler *FollowHandler) SaveAd(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "FollowHandler.SaveAd")
	defer span.End()

	var request domain.Ad
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(writer, "bad request", http.StatusBadRequest)
		return
	}

	err = handler.service.SaveAd(ctx, &request)
	if err != nil {
		http.Error(writer, "internal server error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
