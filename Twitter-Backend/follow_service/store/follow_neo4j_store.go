package store

import (
	"context"
	"fmt"
	"follow_service/domain"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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

	isExist, err := session.ExecuteRead(ctx,
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

	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	requests, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			"MATCH (r:Request)-[:REQUEST_TO]->(u:User) "+
				"WHERE r.receiver = $username AND u.username = $username AND r.status = 1"+
				"RETURN r.id as id, r.requester as requester, r.receiver as receiver, r.status as status",
			map[string]any{"username": username})
		if err != nil {
			log.Printf("Error in getting pending requests of user: %s", err.Error())
			return nil, err
		}

		var requests []*domain.FollowRequest
		if result.Next(ctx) {
			record := result.Record()
			id, _ := record.Get("id")
			requester, _ := record.Get("requester")
			receiver, _ := record.Get("receiver")
			status, _ := record.Get("status")
			requests = append(requests, &domain.FollowRequest{
				ID:        id.(string),
				Requester: requester.(string),
				Receiver:  receiver.(string),
				Status:    status.(domain.Status),
			})
		}

		return requests, nil
	})
	if err != nil {
		return nil, err
	}

	return requests.([]*domain.FollowRequest), nil
}

func (store *FollowNeo4JStore) GetFollowingsOfUser(username string) ([]*domain.User, error) {
	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	followings, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			"MATCH (f:User)-[:FOLLOWS]->(u:User) "+
				"WHERE u.username = $username "+
				"RETURN f.id as id, f.username as username, f.age as age, f.residence as residence",
			map[string]any{"username": username})
		if err != nil {
			log.Printf("Error in getting followings of user: %s", err.Error())
			return nil, err
		}

		var followings []*domain.User
		if result.Next(ctx) {
			record := result.Record()
			id, _ := record.Get("id")
			username, _ := record.Get("username")
			age, _ := record.Get("age")
			residence, _ := record.Get("residence")
			followings = append(followings, &domain.User{
				ID:        id.(string),
				Username:  username.(string),
				Age:       age.(int),
				Residence: residence.(string),
			})
		}

		return followings, nil
	})
	if err != nil {
		return nil, err
	}

	return followings.([]*domain.User), nil
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
				return rid, nil
			}

			return nil, fmt.Errorf("neo4j node and relationships didn't save")
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

func (store *FollowNeo4JStore) SaveFollow(requestID *string) error {
	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx,
				"MATCH (requester:User), (receiver:User), (r:Request) "+
					"WHERE r.id = $requestID requester.username = r.requester AND receiver.username = r.receiver "+
					"CREATE f = (requester)-[:FOLLOWS]->(receiver)",
				map[string]any{"requestID": requestID})
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

func (store *FollowNeo4JStore) AcceptRequest(id *string) error {
	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx,
				"MATCH (r:Request) "+
					"WHERE r.id = $id "+
					"SET r.status = 3",
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
func (store *FollowNeo4JStore) DeclineRequest(id *string) error {
	ctx := context.Background()
	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx,
				"MATCH (r:Request) "+
					"WHERE r.id = $id "+
					"SET r.status = 2",
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
