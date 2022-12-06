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
	router.HandleFunc("/", handler.Register).Methods("POST")
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	router.HandleFunc("/mailExist/{mail}", handler.MailExist).Methods("GET")
	http.Handle("/", router)
}

func (handler *UserHandler) Register(writer http.ResponseWriter, req *http.Request) {
	var user domain.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	saved, err := handler.service.Register(&user)
	if err != nil {
		if err.Error() == errors.DatabaseError {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
		return
	}

	jsonResponse(saved, writer)
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

func (handler *UserHandler) MailExist(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	mail, ok := vars["mail"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := handler.service.DoesEmailExist(mail)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = writer.Write([]byte(id))
	if err != nil {
		log.Println("error in response user service")
		log.Println(err.Error())
		return
	}
}
