package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
	"log"
	"user_service/domain"
)

const (
	DATABASE   = "user"
	COLLECTION = "users"
)

type UserMongoDBStore struct {
	users  *mongo.Collection
	tracer trace.Tracer
}

func NewUserMongoDBStore(client *mongo.Client, tracer trace.Tracer) domain.UserStore {
	users := client.Database(DATABASE).Collection(COLLECTION)
	return &UserMongoDBStore{
		users:  users,
		tracer: tracer,
	}
}

func (store *UserMongoDBStore) GetAll(ctx context.Context) ([]*domain.User, error) {
	ctx, span := store.tracer.Start(ctx, "UserStore.GetAll")
	defer span.End()

	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	ctx, span := store.tracer.Start(ctx, "UserStore.Get")
	defer span.End()

	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := store.tracer.Start(ctx, "UserStore.GetByEmail")
	defer span.End()

	filter := bson.M{"email": email}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) Post(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, span := store.tracer.Start(ctx, "UserStore.Post")
	defer span.End()

	result, err := store.users.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (store *UserMongoDBStore) UpdateUser(ctx context.Context, user *domain.User) error {
	ctx, span := store.tracer.Start(ctx, "UserStore.UpdateUser")
	defer span.End()

	_, err := store.users.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": user})
	if err != nil {
		log.Printf("Updating user error mongodb: %s", err.Error())
		return err
	}

	return nil
}

func (store *UserMongoDBStore) GetOneUser(ctx context.Context, username string) (*domain.User, error) {
	ctx, span := store.tracer.Start(ctx, "UserStore.GetOneUser")
	defer span.End()

	filter := bson.M{"username": username}

	user, err := store.filterOne(filter)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (store *UserMongoDBStore) filter(filter interface{}) ([]*domain.User, error) {
	cursor, err := store.users.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	return decode(cursor)
}

func (store *UserMongoDBStore) filterOne(filter interface{}) (user *domain.User, err error) {
	result := store.users.FindOne(context.TODO(), filter)
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
