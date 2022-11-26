package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"tweet_service/application"
	"tweet_service/domain"
)

type TweetHandler struct {
	service *application.TweetService
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

func NewTweetHandler(service *application.TweetService) *TweetHandler {
	return &TweetHandler{
		service: service,
	}
}

func (handler *TweetHandler) Init(router *mux.Router) {
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	//router.HandleFunc("/{id}", handler.Get).Methods("GET")
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
func (handler *TweetHandler) Post(writer http.ResponseWriter, req *http.Request) {
	var request domain.Tweet
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Header["Token"] == nil {
		fmt.Print("ovde")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	UserID := handler.GetID(handler.GetClaims(req.Header["Token"][0]))
	id, err := gocql.UUIDFromBytes([]byte(UserID))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	tweet, err := handler.service.Post(&request, id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	jsonResponse(tweet, writer)
}

func (handler *TweetHandler) GetClaims(tokenString string) jwt.MapClaims {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println(ok)
		}
		return token, nil
	})

	if err != nil {
		fmt.Println(err)
		return nil

	}

	return token.Claims.(jwt.MapClaims)
}

func (handler *TweetHandler) GetID(claims jwt.MapClaims) string {

	userId := claims["UserID"]
	fmt.Println(userId, claims["Username"].(string))
	return userId.(string)
}
