package handlers

import (
	"Twitter-Backend/application"
	"Twitter-Backend/domain"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type TweetHandler struct {
	service *application.TweetService
}

func NewTweetHandler(service *application.TweetService) *TweetHandler {
	return &TweetHandler{
		service: service,
	}
}

func (handler *TweetHandler) Init(router *mux.Router) {
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	router.HandleFunc("/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/", handler.Post).Methods("POST")
	http.Handle("/", router)
}

func (handler *TweetHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	tweets, err := handler.service.GetAll()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(tweets, writer)
}

func (handler *TweetHandler) Get(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	tweet, err := handler.service.Get(objectId)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	jsonResponse(tweet, writer)
}

func (handler *TweetHandler) Post(writer http.ResponseWriter, req *http.Request) {
	var request domain.Tweet
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.service.Post(&request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
