package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/cristalhq/jwt/v4"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
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
	tracer  trace.Tracer
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))
var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

func NewTweetHandler(service *application.TweetService, tracer trace.Tracer) *TweetHandler {
	return &TweetHandler{
		service: service,
		tracer:  tracer,
	}
}

func (handler *TweetHandler) Init(router *mux.Router) {
	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	log.Println("successful init of enforcer")
	if err != nil {
		log.Fatal(err)
	}

	router.Use(ExtractTraceInfoMiddleware)
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	//router.HandleFunc("/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/", Post(handler)).Methods("POST")
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	router.HandleFunc("/image/{id}", handler.GetTweetImage).Methods("GET")
	router.HandleFunc("/favorite", handler.Favorite).Methods("POST")
	router.HandleFunc("/user/{username}", handler.GetTweetsByUser).Methods("GET")
	router.HandleFunc("/whoLiked/{id}", handler.GetLikesByTweet).Methods("GET")
	router.HandleFunc("/feed", handler.GetFeedByUser).Methods("GET")
	router.HandleFunc("/retweet", handler.Retweet).Methods("POST")
	router.HandleFunc("/timespent", handler.TimespentOnAd).Methods("POST")
	http.Handle("/", router)
	log.Println("Successful")
	log.Fatal(http.ListenAndServe(":8001", authorization.Authorizer(authEnforcer)(router)))
}

func (handler *TweetHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "TweetHandler.GetAll")
	defer span.End()

	tweets, err := handler.service.GetAll(ctx)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(tweets, writer)
}

func (handler *TweetHandler) GetTweetsByUser(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "TweetHandler.GetTweetsByUser")
	defer span.End()

	vars := mux.Vars(req)
	username, ok := vars["username"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tweets, err := handler.service.GetTweetsByUser(ctx, username)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse(tweets, writer)
}

func (handler *TweetHandler) GetLikesByTweet(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "TweetHandler.GetLikesByTweet")
	defer span.End()

	vars := mux.Vars(req)
	tweetID, ok := vars["id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	favorites, err := handler.service.GetLikesByTweet(ctx, tweetID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(favorites, writer)
}

func (handler *TweetHandler) Favorite(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "TweetHandler.Favorite")
	defer span.End()

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

	var tweet domain.Tweet
	err = json.NewDecoder(req.Body).Decode(&tweet)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tweets, err := handler.service.Favorite(ctx, tweet.ID.String(), username, tweet.Advertisement)

	if err != nil {
		log.Printf("Error in tweetHandler Favorite(): %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	jsonResponse(tweets, writer)

}

func (handler *TweetHandler) Post(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "TweetHandler.Post")
	defer span.End()

	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println(err)
		return
	}

	strs := strings.Split(req.Header.Get("Content-Type"), "; boundary")

	if strs[0] != "multipart/form-data" {
		log.Println("Invalid Content-Type")
		http.Error(writer, "Invalid Content-Type", http.StatusBadRequest)
		return
	}

	file, _, err := req.FormFile("image")

	var imageBytes []byte
	if err == nil {
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintln(writer, "Error reading image:", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		imageBytes = bytes
		defer file.Close()
	}

	tweet := req.FormValue("json")

	var tweetVal domain.Tweet
	err = json.Unmarshal([]byte(tweet), &tweetVal)
	if err != nil {
		http.Error(writer, "bad json format", http.StatusBadRequest)
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

	ret, err := handler.service.Post(ctx, &tweetVal, username, &imageBytes)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	jsonResponse(ret, writer)
}

func (handler *TweetHandler) GetTweetImage(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		http.Error(writer, "bad request", http.StatusBadRequest)
		return
	}

	image, err := handler.service.GetTweetImage(id)
	if err != nil {
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(*image)
}

func Post(handler *TweetHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.Post(w, r)
	}
}

func ExtractTraceInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *TweetHandler) GetFeedByUser(writer http.ResponseWriter, req *http.Request) {
	feed, err := handler.service.GetFeedByUser(req.Header.Get("Authorization"))
	if err != nil {
		log.Printf("error: %s", err.Error())
		if err.Error() == "FollowServiceError" {
			http.Error(writer, err.Error(), http.StatusServiceUnavailable)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	} else {
		jsonResponse(feed, writer)

	}
}

func (handler *TweetHandler) Retweet(writer http.ResponseWriter, req *http.Request) {

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

	code, err := handler.service.Retweet(tweetID.ID, username)
	if err != nil {
		http.Error(writer, err.Error(), code)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *TweetHandler) TimespentOnAd(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "TweetHandler.TimespentOnAd")
	defer span.End()

	var timespent domain.Timespent
	err := json.NewDecoder(req.Body).Decode(&timespent)
	if err != nil {
		log.Println("Error in decoding body in handler function TimespentOnAd")
		http.Error(writer, "bad request", http.StatusBadRequest)
	}

	handler.service.TimeSpentOnAd(ctx, &timespent)

	writer.WriteHeader(http.StatusOK)
}
