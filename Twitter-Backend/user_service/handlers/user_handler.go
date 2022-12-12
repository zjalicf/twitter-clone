package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/cristalhq/jwt/v4"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
	"strings"
	"user_service/application"
	"user_service/authorization"
	"user_service/domain"
	"user_service/errors"
)

type UserHandler struct {
	service *application.UserService
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))
var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Init(router *mux.Router) {

	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	log.Println("sucessful init of enforcer")
	if err != nil {
		log.Fatal(err)
	}

	//test

	router.HandleFunc("/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/", handler.Register).Methods("POST")
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	router.HandleFunc("/getOne/{username}", handler.GetOne).Methods("GET")
	router.HandleFunc("/getMe/", handler.GetMe).Methods("GET")
	router.HandleFunc("/mailExist/{mail}", handler.MailExist).Methods("GET")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8002", authorization.Authorizer(authEnforcer)(router)))
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

func (handler *UserHandler) GetOne(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	username := vars["username"]

	user, err := handler.service.GetOneUser(username)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusNotFound)
	}
	jsonResponse(user, writer)
}

func (handler *UserHandler) GetMe(writer http.ResponseWriter, request *http.Request) {
	bearer := request.Header.Get("Authorization")
	bearerToken := strings.Split(bearer, "Bearer ")
	tokenString := bearerToken[1]
	fmt.Println(tokenString)
	token, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		log.Println(err)
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims := authorization.GetMapClaims(token.Bytes())
	username := claims["username"]

	user, err := handler.service.GetOneUser(username)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusNotFound)
	}
	jsonResponse(user, writer)
}
