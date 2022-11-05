package handlers

import (
	"Twitter-Backend/data"
	"context"
	"log"
	"net/http"
)

type KeyTweet struct{}

type TweetsHandler struct {
	logger    *log.Logger
	tweetRepo data.TweetRepo
}

func New(logger *log.Logger, tweetRepo data.TweetRepo) *TweetsHandler {
	return &TweetsHandler{logger, tweetRepo}
}

func (tweetsHandler *TweetsHandler) GetTweets(responseWriter http.ResponseWriter, request *http.Request) {
	tweets := tweetsHandler.tweetRepo.GetAll()
	err := tweets.ToJSON(responseWriter)

	if err != nil {
		http.Error(responseWriter, "Unable to convert to JSON", http.StatusInternalServerError)
		tweetsHandler.logger.Println("Unable to convert to JSON :", err)
		return
	}
}

func (tweetsHandler *TweetsHandler) PostTweet(responseWriter http.ResponseWriter, request *http.Request) {
	tweet := request.Context().Value(KeyTweet{}).(*data.Tweet)
	tweetsHandler.tweetRepo.PostTweet(tweet)
	responseWriter.WriteHeader(http.StatusCreated)
}

func (tweetsHandler *TweetsHandler) MiddlewareTweetValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		tweet := &data.Tweet{}
		err := tweet.FromJSON(request.Body)
		if err != nil {
			http.Error(responseWriter, "Unable to Decode JSON", http.StatusBadRequest)
			tweetsHandler.logger.Println(err)
			return
		}

		//err = tweet.Validate()
		//
		//if err != nil {
		//	tweetsHandler.logger.Println("Error Validation Tweet", err)
		//	http.Error(responseWriter, fmt.Sprintf("Error Validating tweet: %s", err), http.StatusBadRequest)
		//	return
		//}

		ctx := context.WithValue(request.Context(), KeyTweet{}, tweet)
		request = request.WithContext(ctx)

		next.ServeHTTP(responseWriter, request)
	})
}

func (tweetsHandler *TweetsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		tweetsHandler.logger.Println("Method [", request.Method, "] - Hit path :", request.URL.Path)

		responseWriter.Header().Add("Content-Type", "application-json")

		next.ServeHTTP(responseWriter, request)
	})
}
