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
	"log"
	"net/http"
)

type FollowHandler struct {
	service            *application.FollowService
	counterUnavailable int
}

func NewFollowHandler(service *application.FollowService) *FollowHandler {
	return &FollowHandler{
		service:            service,
		counterUnavailable: 3,
	}
}

func (handler *FollowHandler) Init(router *mux.Router) {

	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	log.Println("successful init of enforcer")
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/", handler.GetAll).Methods("GET")
	router.HandleFunc("/requests/", handler.GetRequestsForUser).Methods("GET")
	router.HandleFunc("/requests/{visibility}", handler.CreateRequest).Methods("POST")
	//router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	router.HandleFunc("/acceptRequest/{id}", handler.AcceptRequest).Methods("PUT")
	router.HandleFunc("/declineRequest/{id}", handler.DeclineRequest).Methods("PUT")
	router.HandleFunc("/followings", handler.GetFollowingsByUser).Methods("GET")
	router.HandleFunc("/recommendations", handler.GetRecommendationsForUser).Methods("GET")
	router.HandleFunc("/ad", handler.SaveAd).Methods("POST")

	http.Handle("/", router)
	log.Println("Successful")
	log.Fatal(http.ListenAndServe(":8004", authorization.Authorizer(authEnforcer)(router)))
}

func (handler *FollowHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	tweets, err := handler.service.GetAll()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(tweets, writer)
}

func (handler *FollowHandler) GetRequestsForUser(writer http.ResponseWriter, req *http.Request) {
	token, err := authorization.GetToken(req)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := authorization.GetMapClaims(token.Bytes())

	returnRequests, err := handler.service.GetRequestsForUser(claims["username"])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(returnRequests)

	jsonResponse(returnRequests, writer)
}

func (handler *FollowHandler) GetFollowingsByUser(writer http.ResponseWriter, req *http.Request) {
	token, _ := authorization.GetToken(req)
	claims := authorization.GetMapClaims(token.Bytes())
	username := claims["username"]

	log.Printf("username is: %s", username)
	users, err := handler.service.GetFollowingsOfUser(username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(users, writer)

}

func (handler *FollowHandler) GetRecommendationsForUser(writer http.ResponseWriter, req *http.Request) {
	token, _ := authorization.GetToken(req)
	claims := authorization.GetMapClaims(token.Bytes())
	username := claims["username"]
	//vars := mux.Vars(req)
	//username, ok := vars["username"]
	//if !ok {
	//	http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	//}

	users, err := handler.service.GetRecommendationsByUsername(username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(users, writer)
}

//	func (handler *TweetHandler) Get(writer http.ResponseWriter, req *http.Request) {
//		vars := mux.Vars(req)
//		id, ok := vars["id"]
//		if !ok {
//			writer.WriteHeader(http.StatusBadRequest)
//			return
//		}
//		objectId, err := primitive.ObjectIDFromHex(id)
//		if err != nil {
//			writer.WriteHeader(http.StatusBadRequest)
//			return
//		}
//		tweet, err := handler.service.Get(objectId)
//		if err != nil {
//			writer.WriteHeader(http.StatusNotFound)
//			return
//		}
//		jsonResponse(tweet, writer)
//	}
func (handler *FollowHandler) CreateRequest(writer http.ResponseWriter, req *http.Request) {

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

	err = handler.service.CreateRequest(&request, claims["username"], visibility)
	if err != nil {
		log.Println("ERR")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	//jsonResponse(followRequest, writer)
}

func (handler *FollowHandler) CreateUser(writer http.ResponseWriter, req *http.Request) {
	var request domain.User
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.service.CreateUser(&request)
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, errors.ServiceUnavailable, http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *FollowHandler) AcceptRequest(writer http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	followId, ok := vars["id"]
	if !ok {
		http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	}

	err := handler.service.AcceptRequest(&followId)
	if err != nil {
		http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	}

	writer.WriteHeader(http.StatusOK)

}

func (handler *FollowHandler) DeclineRequest(writer http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	followId, ok := vars["id"]
	if !ok {
		http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	}

	err := handler.service.DeclineRequest(&followId)
	if err != nil {
		http.Error(writer, errors.BadRequestError, http.StatusBadRequest)
	}

	writer.WriteHeader(http.StatusOK)

}

func (handler *FollowHandler) SaveAd(writer http.ResponseWriter, req *http.Request) {
	log.Println("Uslo u handler")

	var request domain.Ad
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(writer, "bad request", http.StatusBadRequest)
		return
	}

	err = handler.service.SaveAd(&request)
	if err != nil {
		http.Error(writer, "internal server error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
