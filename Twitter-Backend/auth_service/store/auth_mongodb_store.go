package store

import (
	"auth_service/domain"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "user_credentials"
	COLLECTION = "credentials"
)

type AuthMongoDBStore struct {
	credentials *mongo.Collection
}

func (store *AuthMongoDBStore) GetAll() ([]*domain.User, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func NewAuthMongoDBStore(client *mongo.Client) domain.AuthStore {
	auths := client.Database(DATABASE).Collection(COLLECTION)
	return &AuthMongoDBStore{
		credentials: auths,
	}
}

func (store *AuthMongoDBStore) Register(user *domain.Credentials) error {
	result, err := store.credentials.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (store *AuthMongoDBStore) ChangePassword(user *domain.User) error {

	fmt.Println(user)
	newState, err := store.credentials.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	fmt.Println(newState)
	return nil
}

func (store *AuthMongoDBStore) GetOneUser(username string) (*domain.User, error) {
	filter := bson.M{"username": username}

	user, err := store.filterOne(filter)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (store *AuthMongoDBStore) GetOneUserByID(id primitive.ObjectID) *domain.User {
	filter := bson.M{"id": id}

	var user *domain.User
	err := store.credentials.FindOne(context.TODO(), filter, nil).Decode(user)
	if err != nil {
		return nil
	}

	return user
}

func (store *AuthMongoDBStore) filter(filter interface{}) ([]*domain.User, error) {
	cursor, err := store.credentials.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *AuthMongoDBStore) filterOne(filter interface{}) (user *domain.User, err error) {
	result := store.credentials.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func decode(cursor *mongo.Cursor) (users []*domain.User, err error) {
	for cursor.Next(context.TODO()) {
		var user domain.User
		err = cursor.Decode(&user)
		if err != nil {
			return
		}
		users = append(users, &user)
	}
	err = cursor.Err()
	return
}
