package store

import (
	"context"
	"fmt"
	"follow_service/domain"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

const (
	DATABASE   = "follow"
	COLLECTION = "follows"
)

type FollowNeo4JStore struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func NewFollowNeo4JStore(driver *neo4j.DriverWithContext) domain.FollowRequestStore {
	return &FollowNeo4JStore{
		driver: *driver,
		logger: log.Default(),
	}
}

func (store *FollowNeo4JStore) GetAll() ([]*domain.FollowRequest, error) {
	//filter := bson.D{{}}
	//return store.filter(filter)
	return nil, nil
}

func (store *FollowNeo4JStore) FollowExist(followRequest *domain.FollowRequest) (bool, error) {

	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	isExist, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (req)-[:FOLLOWS]->(rec) "+
					"WHERE req.username = $requester AND rec.username = $receiver "+
					"RETURN rec as receiver",
				map[string]any{"requester": followRequest.Requester, "receiver": followRequest.Receiver})
			if err != nil {
				log.Printf("Error in creating request node and relationships because of: %s", err.Error())
				return nil, err
			}

			if result.Next(ctx) {
				record := result.Record()
				receiver, ok := record.Get("receiver")
				if !ok || receiver == nil {
					return false, nil
				}
			} else {
				return false, nil
			}
			return true, nil
		})
	if err != nil {
		return false, err
	}

	return isExist.(bool), nil
}

func (store *FollowNeo4JStore) GetRequestsForUser(username string) ([]*domain.FollowRequest, error) {
	//filter := bson.D{{"receiver", username}, {"status", 1}}
	//result, err := store.filter(filter)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

func (store *FollowNeo4JStore) GetFollowingsOfUser(username string) ([]*domain.FollowRequest, error) {
	//filter := bson.D{{"requester", username}, {"status", domain.Accepted}}
	//result, err := store.filter(filter)
	//log.Println("FOLLOW MONGO RESULT: ")
	//log.Println(result)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

func (store *FollowNeo4JStore) SaveRequest(request *domain.FollowRequest) error {

	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (requester:User), (receiver:User) "+
					"WHERE requester.username = $requester AND receiver.username = $receiver "+
					"CREATE (r:Request) SET r.id = $id, r.requester = $requester, "+
					"r.receiver = $receiver, r.status = $status "+
					"CREATE p = (requester)-[:CREATED]->(r)-[:REQUEST_TO]->(receiver) "+
					"RETURN r.id as rid",
				map[string]any{"id": request.ID, "requester": request.Requester, "receiver": request.Receiver,
					"status": request.Status.EnumIndex()})
			if err != nil {
				log.Printf("Error in creating request node and relationships because of: %s", err.Error())
				return nil, err
			}

			if result.Next(ctx) {
				rid, ok := result.Record().Get("rid")
				if !ok || rid == nil {
					return nil, fmt.Errorf("neo4j node and relationships not saved")
				}
				log.Println("OK")
				return rid, nil
			}

			return nil, fmt.Errorf("neo4j node and relationships not saved")
		})
	if err != nil {
		return err
	}

	return nil
}

func (store *FollowNeo4JStore) SaveUser(user *domain.User) error {

	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (u:User) SET u.id = $id, u.username = $username, "+
					"u.age = $age, u.residence = $residence RETURN u.id + ', from node ' + id(u)",
				map[string]any{"id": user.ID, "username": user.Username, "age": user.Age,
					"residence": user.Residence})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		return err
	}

	return nil
}

func (store *FollowNeo4JStore) DeleteUser(id *string) error {
	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		_, err := transaction.Run(ctx,
			"MATCH (u:User) "+
				"WHERE u.id = $id "+
				"DELETE u",
			map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (store *FollowNeo4JStore) AcceptRequest(id primitive.ObjectID) error {

	//filter := bson.D{{"_id", id}}
	//update := bson.D{{"$set", bson.D{{"status", 3}}}}
	//
	//_, err := store.follows.UpdateOne(context.TODO(), filter, update)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (store *FollowNeo4JStore) DeclineRequest(id primitive.ObjectID) error {

	//filter := bson.D{{"_id", id}}
	//
	//_, err := store.follows.DeleteOne(context.TODO(), filter)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (store *FollowNeo4JStore) HandleRequest() {
	//TODO implement me
	panic("implement me")
}

//func (store *FollowNeo4JStore) filter(filter interface{}) ([]*domain.FollowRequest, error) {
//	cursor, err := store.follows.Find(context.TODO(), filter)
//	defer cursor.Close(context.TODO())
//
//	if err != nil {
//		return nil, err
//	}
//
//	return decode(cursor)
//}
//
//func (store *FollowNeo4JStore) filterOne(filter interface{}) (follow *domain.FollowRequest, err error) {
//	result := store.follows.FindOne(context.TODO(), filter)
//	err = result.Decode(&follow)
//	return
//}
//
//func decode(cursor *mongo.Cursor) (follows []*domain.FollowRequest, err error) {
//	for cursor.Next(context.TODO()) {
//		var follow domain.FollowRequest
//		err = cursor.Decode(&follow)
//		if err != nil {
//			return
//		}
//		follows = append(follows, &follow)
//	}
//	err = cursor.Err()
//	return
//}

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
