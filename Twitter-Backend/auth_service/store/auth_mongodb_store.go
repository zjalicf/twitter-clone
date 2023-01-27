package store

import (
	"auth_service/domain"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
	"log"
)

const (
	DATABASE   = "user_credentials"
	COLLECTION = "credentials"
)

type AuthMongoDBStore struct {
	credentials *mongo.Collection
	tracer      trace.Tracer
}

func NewAuthMongoDBStore(client *mongo.Client, tracer trace.Tracer) domain.AuthStore {
	auths := client.Database(DATABASE).Collection(COLLECTION)
	return &AuthMongoDBStore{
		credentials: auths,
		tracer:      tracer,
	}
}

func (store *AuthMongoDBStore) GetAll(ctx context.Context) ([]*domain.Credentials, error) {
	ctx, span := store.tracer.Start(ctx, "AuthStore.GetAll")
	defer span.End()

	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *AuthMongoDBStore) Register(ctx context.Context, user *domain.Credentials) error {
	ctx, span := store.tracer.Start(ctx, "AuthStore.Register")
	defer span.End()

	//vratiti u jednom trenutku
	user.Verified = true

	result, err := store.credentials.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (store *AuthMongoDBStore) UpdateUser(ctx context.Context, user *domain.Credentials) error {
	ctx, span := store.tracer.Start(ctx, "AuthStore.UpdateUser")
	defer span.End()

	newState, err := store.credentials.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	fmt.Println(newState)
	return nil
}

func (store *AuthMongoDBStore) GetOneUser(ctx context.Context, username string) (*domain.Credentials, error) {
	ctx, span := store.tracer.Start(ctx, "AuthStore.GetOneUser")
	defer span.End()

	filter := bson.M{"username": username}

	user, err := store.filterOne(filter)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Println("SVE OK")

	return user, nil
}

func (store *AuthMongoDBStore) GetOneUserByID(ctx context.Context, id primitive.ObjectID) *domain.Credentials {
	ctx, span := store.tracer.Start(ctx, "AuthStore.GetOneUserByID")
	defer span.End()

	filter := bson.M{"_id": id}

	var user domain.Credentials
	err := store.credentials.FindOne(context.TODO(), filter, nil).Decode(&user)
	if err != nil {
		return nil
	}

	return &user
}

func (store *AuthMongoDBStore) DeleteUserByID(ctx context.Context, id primitive.ObjectID) error {
	ctx, span := store.tracer.Start(ctx, "AuthStore.DeleteUserByID")
	defer span.End()

	_, err := store.credentials.DeleteMany(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (store *AuthMongoDBStore) filter(filter interface{}) ([]*domain.Credentials, error) {
	cursor, err := store.credentials.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *AuthMongoDBStore) filterOne(filter interface{}) (user *domain.Credentials, err error) {
	result := store.credentials.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func decode(cursor *mongo.Cursor) (users []*domain.Credentials, err error) {
	for cursor.Next(context.TODO()) {
		var user domain.Credentials
		err = cursor.Decode(&user)
		if err != nil {
			return
		}
		users = append(users, &user)
	}
	err = cursor.Err()
	return
}
