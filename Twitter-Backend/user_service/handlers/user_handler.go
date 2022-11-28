package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"user_service/application"
	"user_service/domain"
	"user_service/errors"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Init(router *mux.Router) {
	router.HandleFunc("/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	http.Handle("/", router)
}

func (handler *UserHandler) Login(writer http.ResponseWriter, request *http.Request) {

	var user domain.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	token, err := handler.service.Login(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}
	writer.Write([]byte(token))

}

func (handler *UserHandler) Register(writer http.ResponseWriter, req *http.Request) {
	var user domain.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	savedUser, err := handler.service.Register(&user)
	if err != nil {
		if err.Error() == errors.DatabaseError {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
		return
	}

  writer.WriteHeader(http.StatusCreated)
	//jsonResponse(saved, writer)
}

func (handler *UserHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	users, err := handler.service.GetAll()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(users, writer)
}

func (handler *UserHandler) Get(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := handler.service.Get(objectID)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	jsonResponse(user, writer)
}
