package handlers

import (
	"encoding/json"
	"follow_service/application"
	"follow_service/domain"
	"github.com/casbin/casbin"
	"github.com/cristalhq/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

type FollowHandler struct {
	service *application.FollowService
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))
var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

func NewTweetHandler(service *application.FollowService) *FollowHandler {
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
	router.HandleFunc("/", Post(handler)).Methods("POST")
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	http.Handle("/", router)
	log.Println("Successful")
	//log.Fatal(http.ListenAndServe(":8001", authorization.Authorizer(authEnforcer)(router)))
}

func (handler *FollowHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	tweets, err := handler.service.GetAll()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(tweets, writer)
}

func (handler *FollowHandler) GetTweetsByUser(writer http.ResponseWriter, req *http.Request) {
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
func (handler *FollowHandler) Post(writer http.ResponseWriter, req *http.Request) {

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
