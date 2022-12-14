package handlers

import (
	"encoding/json"
	"follow_service/application"
	"follow_service/authorization"
	"follow_service/domain"
	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type FollowHandler struct {
	service *application.FollowService
}

func NewFollowHandler(service *application.FollowService) *FollowHandler {
	return &FollowHandler{
		service: service,
	}
}

func (handler *FollowHandler) Init(router *mux.Router) {

	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	log.Println("sucessful init of enforcer")
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/", handler.GetAll).Methods("GET")
	//router.HandleFunc("/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/", handler.CreateRequest).Methods("POST")
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

//func (handler *FollowHandler) GetTweetsByUser(writer http.ResponseWriter, req *http.Request) {
//	vars := mux.Vars(req)
//	username, ok := vars["username"]
//	if !ok {
//		writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	tweets, err := handler.service.GetTweetsByUser(username)
//	if err != nil {
//		writer.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	jsonResponse(tweets, writer)
//}

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

	followRequest, err := handler.service.CreateRequest(&request, claims["username"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	jsonResponse(followRequest, writer)
}
