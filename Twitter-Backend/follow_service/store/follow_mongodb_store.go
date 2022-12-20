package store

import (
	"context"
	"follow_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const (
	DATABASE   = "follow"
	COLLECTION = "follows"
)

type FollowMongoDBStore struct {
	follows *mongo.Collection
}

func NewFollowMongoDBStore(client *mongo.Client) domain.FollowRequestStore {
	follows := client.Database(DATABASE).Collection(COLLECTION)
	return &FollowMongoDBStore{
		follows: follows,
	}
}

func (store *FollowMongoDBStore) GetAll() ([]*domain.FollowRequest, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *FollowMongoDBStore) GetRequestsForUser(username string) ([]*domain.FollowRequest, error) {
	filter := bson.D{{"receiver", username}, {"status", 1}}
	result, err := store.filter(filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (store *FollowMongoDBStore) GetFollowingsOfUser(username string) ([]*domain.FollowRequest, error) {
	filter := bson.D{{"requester", username}, {"status", domain.Accepted}}
	result, err := store.filter(filter)
	log.Println("FOLLOW MONGO RESULT: ")
	log.Println(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (store *FollowMongoDBStore) SaveRequest(request *domain.FollowRequest) (*domain.FollowRequest, error) {

	result, err := store.follows.InsertOne(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	request.ID = result.InsertedID.(primitive.ObjectID)

	return request, nil
}

func (store *FollowMongoDBStore) AcceptRequest(id primitive.ObjectID) error {

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"status", 3}}}}

	_, err := store.follows.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (store *FollowMongoDBStore) DeclineRequest(id primitive.ObjectID) error {

	filter := bson.D{{"_id", id}}

	_, err := store.follows.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (store *FollowMongoDBStore) HandleRequest() {
	//TODO implement me
	panic("implement me")
}

func (store *FollowMongoDBStore) filter(filter interface{}) ([]*domain.FollowRequest, error) {
	cursor, err := store.follows.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	return decode(cursor)
}

func (store *FollowMongoDBStore) filterOne(filter interface{}) (follow *domain.FollowRequest, err error) {
	result := store.follows.FindOne(context.TODO(), filter)
	err = result.Decode(&follow)
	return
}

func decode(cursor *mongo.Cursor) (follows []*domain.FollowRequest, err error) {
	for cursor.Next(context.TODO()) {
		var follow domain.FollowRequest
		err = cursor.Decode(&follow)
		if err != nil {
			return
		}
		follows = append(follows, &follow)
	}
	err = cursor.Err()
	return
}

//func (store *FollowMongoDBStore) GetAll() ([]*domain.User, error) {
//	filter := bson.D{{}}
//	return store.filter(filter)
//}
//
//func (store *FollowMongoDBStore) Get(id primitive.ObjectID) (*domain.User, error) {
//	filter := bson.M{"_id": id}
//	return store.filterOne(filter)
//}
//
//func (store *FollowMongoDBStore) GetByEmail(email string) (*domain.User, error) {
//	filter := bson.M{"email": email}
//	return store.filterOne(filter)
//}
//
//func (store *FollowMongoDBStore) Post(user *domain.User) (*domain.User, error) {
//	user.ID = primitive.NewObjectID()
//
//	result, err := store.users.InsertOne(context.TODO(), user)
//
//	if err != nil {
//		return nil, err
//	}
//
//	user.ID = result.InsertedID.(primitive.ObjectID)
//
//	return user, nil
//}
//
//func (store *FollowMongoDBStore) UpdateUser(user *domain.User) error {
//	_, err := store.users.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": user})
//	if err != nil {
//		log.Printf("Updating user error mongodb: %s", err.Error())
//		return err
//	}
//
//	return nil
//}
//

//
//func (store *FollowMongoDBStore) GetOneUser(username string) (*domain.User, error) {
//
//	filter := bson.M{"username": username}
//
//	user, err := store.filterOne(filter)
//	if err != nil {
//		return nil, err
//	}
//
//	return user, nil
//}
