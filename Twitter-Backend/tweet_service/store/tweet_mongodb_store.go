package store

import (
	"Twitter-Backend/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "tweet"
	COLLECTION = "tweet"
)

type TweetMongoDBStore struct {
	tweets *mongo.Collection
}

func NewTweetMongoDBStore(client *mongo.Client) domain.TweetStore {
	tweets := client.Database(DATABASE).Collection(COLLECTION)
	return &TweetMongoDBStore{
		tweets: tweets,
	}
}

func (store *TweetMongoDBStore) GetAll() ([]*domain.Tweet, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *TweetMongoDBStore) Get(id primitive.ObjectID) (*domain.Tweet, error) {
	filter := bson.M{"id": id}
	return store.filterOne(filter)
}

func (store *TweetMongoDBStore) Post(tweet *domain.Tweet) error {
	tweet.ID = primitive.NewObjectID()
	result, err := store.tweets.InsertOne(context.TODO(), tweet)
	if err != nil {
		return err
	}
	tweet.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *TweetMongoDBStore) DeleteAll() {
	store.tweets.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *TweetMongoDBStore) filter(filter interface{}) ([]*domain.Tweet, error) {
	cursor, err := store.tweets.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *TweetMongoDBStore) filterOne(filter interface{}) (Tweet *domain.Tweet, err error) {
	result := store.tweets.FindOne(context.TODO(), filter)
	err = result.Decode(&Tweet)
	return
}

func decode(cursor *mongo.Cursor) (tweets []*domain.Tweet, err error) {
	for cursor.Next(context.TODO()) {
		var Tweet domain.Tweet
		err = cursor.Decode(&Tweet)
		if err != nil {
			return
		}
		tweets = append(tweets, &Tweet)
	}
	err = cursor.Err()
	return
}
