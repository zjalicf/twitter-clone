package handlers

import (
	"encoding/json"
	"github.com/casbin/casbin"
	"github.com/cristalhq/jwt/v4"
	"github.com/gocql/gocql"
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
	router.HandleFunc("/favorite", handler.Favorite).Methods("POST")
	router.HandleFunc("/user/{username}", handler.GetTweetsByUser).Methods("GET")
	router.HandleFunc("/whoLiked/{id}", handler.GetLikesByTweet).Methods("GET")
	router.HandleFunc("/feed", handler.GetFeedByUser).Methods("GET")
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

func (handler *TweetHandler) GetLikesByTweet(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tweetID, ok := vars["id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	favorites, err := handler.service.GetLikesByTweet(tweetID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(favorites, writer)
}

func (handler *TweetHandler) Favorite(writer http.ResponseWriter, req *http.Request) {

	bearer := req.Header.Get("Authorization")
	bearerToken := strings.Split(bearer, "Bearer ")
	tokenString := bearerToken[1]

	token, err := jwt.Parse([]byte(tokenString), verifier)

	if err != nil {
		log.Println(err)
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims := authorization.GetMapClaims(token.Bytes())
	username := claims["username"]

	var tweetID domain.TweetID
	err = json.NewDecoder(req.Body).Decode(&tweetID)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tweets, err := handler.service.Favorite(tweetID.ID, username)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	jsonResponse(tweets, writer)

}

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

func (handler *TweetHandler) GetFeedByUser(writer http.ResponseWriter, req *http.Request) {
	feed, err := handler.service.GetFeedByUser(req.Header.Get("Authorization"))
	if err != nil {
		log.Printf("error: %s", err.Error())
		return
	}

	jsonResponse(feed, writer)
}
