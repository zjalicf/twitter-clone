package handlers

import (
	"encoding/json"
	"github.com/casbin/casbin"
	"github.com/cristalhq/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"tweet_service/application"
	"tweet_service/authorization"
	"tweet_service/domain"
)

type TweetHandler struct {
	service *application.TweetService
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))
var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

func NewTweetHandler(service *application.TweetService) *TweetHandler {
	return &TweetHandler{
		service: service,
	}
}

func (handler *TweetHandler) Init(router *mux.Router) {

	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	log.Println("sucessful init of enforcer")
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/", handler.GetAll).Methods("GET")
	//router.HandleFunc("/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/", Post(handler)).Methods("POST")
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	router.HandleFunc("/user/{username}", handler.GetTweetsByUser).Methods("GET")
	http.Handle("/", router)
	log.Println("Successful")
	log.Fatal(http.ListenAndServe(":8001", authorization.Authorizer(authEnforcer)(router)))
}

func (handler *TweetHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	tweets, err := handler.service.GetAll()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(tweets, writer)
}

func (handler *TweetHandler) GetTweetsByUser(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username, ok := vars["username"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tweets, err := handler.service.GetTweetsByUser(username)
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

	if req.Header.Get("Authorization") == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	bearerToken := strings.Split(req.Header.Get("Authorization"), "Bearer ")
	tokenString := bearerToken[1]
	token := authorization.GetToken(tokenString)

	claims := authorization.GetMapClaims(token.Bytes())
	username := claims["username"]

	tweet, err := handler.service.Post(&request, username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	jsonResponse(tweet, writer)
}

func Post(handler *TweetHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.Post(w, r)
	}
}
