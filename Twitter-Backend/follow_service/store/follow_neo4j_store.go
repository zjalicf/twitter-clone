package store

import (
	"context"
	"fmt"
	"follow_service/domain"
	"follow_service/errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"log"
)

const (
	DATABASE   = "follow"
	COLLECTION = "follows"
)

type FollowNeo4JStore struct {
	driver  neo4j.DriverWithContext
	logger  *log.Logger
	logging *logrus.Logger
	tracer  trace.Tracer
}

func NewFollowNeo4JStore(driver *neo4j.DriverWithContext, tracer trace.Tracer, logging *logrus.Logger) domain.FollowRequestStore {
	return &FollowNeo4JStore{
		driver:  *driver,
		logger:  log.Default(),
		tracer:  tracer,
		logging: logging,
	}
}

func (store *FollowNeo4JStore) FollowExist(ctx context.Context, followRequest *domain.FollowRequest) (bool, error) {
	ctx, span := store.tracer.Start(ctx, "FollowStore.FollowExist")
	defer span.End()

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

func (store *FollowNeo4JStore) GetRequestsForUser(ctx context.Context, username string) ([]*domain.FollowRequest, error) {
	ctx, span := store.tracer.Start(ctx, "FollowStore.GetRequestsForUser")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	requests, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			"MATCH (r:Request)-[:REQUEST_TO]->(u:User) "+
				"WHERE r.receiver = $username AND u.username = $username AND r.status = 1 "+
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
				Status:    domain.Status(status.(int64)),
			})
		}

		return requests, nil
	})
	if err != nil {
		return nil, err
	}

	return requests.([]*domain.FollowRequest), nil
}

func (store *FollowNeo4JStore) GetRequestByRequesterReceiver(ctx context.Context, requester, receiver *string) (*domain.FollowRequest, error) {
	ctx, span := store.tracer.Start(ctx, "FollowStore.GetRequestByRequesterReceiver")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	requests, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			"MATCH (r:Request) "+
				"WHERE r.requester = $requester AND r.receiver = $receiver "+
				"RETURN r.id as id, r.requester as requester, r.receiver as receiver, r.status as status",
			map[string]any{"requester": requester, "receiver": receiver})
		if err != nil {
			log.Printf("Error in getting request by requester and receiver: %s", err.Error())
			return nil, err
		}

		var request *domain.FollowRequest
		if result.Next(ctx) {
			record := result.Record()
			id, _ := record.Get("id")
			requester, _ := record.Get("requester")
			receiver, _ := record.Get("receiver")
			status, _ := record.Get("status")
			request = &domain.FollowRequest{
				ID:        id.(string),
				Requester: requester.(string),
				Receiver:  receiver.(string),
				Status:    domain.Status(status.(int64)),
			}
		} else {
			return nil, fmt.Errorf(errors.ErrorRequestNotExists)
		}
		return request, nil
	})
	if err != nil {
		return nil, err
	}

	return requests.(*domain.FollowRequest), nil
}

func (store *FollowNeo4JStore) GetFollowingsOfUser(ctx context.Context, username string) ([]*domain.User, error) {
	ctx, span := store.tracer.Start(ctx, "FollowStore.GetFollowingsOfUser")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	followings, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			"MATCH (f:User)-[:FOLLOWS]->(u:User) "+
				"WHERE f.username = $username "+
				"RETURN u.id as id, u.username as username, u.age as age, u.residence as residence",
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
				Age:       int(age.(int64)),
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

func (store *FollowNeo4JStore) SaveRequest(ctx context.Context, request *domain.FollowRequest) error {
	ctx, span := store.tracer.Start(ctx, "FollowStore.SaveRequest")
	defer span.End()

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

func (store *FollowNeo4JStore) UpdateRequest(ctx context.Context, request *domain.FollowRequest) error {
	ctx, span := store.tracer.Start(ctx, "FollowStore.UpdateRequest")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx,
				"MATCH (request:Request) "+
					"WHERE request.id = $id "+
					"SET request.status = $status",
				map[string]any{"id": request.ID, "status": request.Status})
			if err != nil {
				log.Printf("Error in creating request node and relationships because of: %s", err.Error())
				return nil, err
			}

			return nil, nil
		})
	if err != nil {
		return err
	}

	return nil
}

func (store *FollowNeo4JStore) SaveUser(ctx context.Context, user *domain.User) error {
	ctx, span := store.tracer.Start(ctx, "FollowStore.SaveUser")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (u:User) SET u.id = $id, u.username = $username, "+
					"u.age = $age, u.residence = $residence, u.gender = $gender RETURN u.id + ', from node ' + id(u)",
				map[string]any{"id": user.ID, "username": user.Username, "age": user.Age,
					"residence": user.Residence, "gender": user.Gender})
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

func (store *FollowNeo4JStore) DeleteUser(ctx context.Context, id *string) error {
	ctx, span := store.tracer.Start(ctx, "FollowStore.DeleteUser")
	defer span.End()

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

func (store *FollowNeo4JStore) SaveFollow(ctx context.Context, request *domain.FollowRequest) error {
	ctx, span := store.tracer.Start(ctx, "FollowStore.SaveFollow")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx,
				"MATCH (requester:User), (receiver:User) "+
					"WHERE requester.username = $requester AND receiver.username = $receiver "+
					"CREATE f = (requester)-[:FOLLOWS]->(receiver)",
				map[string]any{"requester": request.Requester, "receiver": request.Receiver})
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

func (store *FollowNeo4JStore) AcceptRequest(ctx context.Context, id *string) (*domain.FollowRequest, error) {
	ctx, span := store.tracer.Start(ctx, "FollowStore.AcceptRequest")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	request, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (r:Request) "+
					"WHERE r.id = $id "+
					"SET r.status = 3 "+
					"RETURN r.id as id, r.requester as requester, r.receiver as receiver, r.status as status",
				map[string]any{"id": id})
			if err != nil {
				return nil, err
			}

			var request *domain.FollowRequest
			if result.Next(ctx) {
				record := result.Record()
				id, _ := record.Get("id")
				requester, _ := record.Get("requester")
				receiver, _ := record.Get("receiver")
				status, _ := record.Get("status")
				request = &domain.FollowRequest{
					ID:        id.(string),
					Requester: requester.(string),
					Receiver:  receiver.(string),
					Status:    domain.Status(status.(int64)),
				}
			}

			return request, nil
		})
	if err != nil {
		return nil, err
	}

	return request.(*domain.FollowRequest), nil
}
func (store *FollowNeo4JStore) DeclineRequest(ctx context.Context, id *string) error {
	ctx, span := store.tracer.Start(ctx, "FollowStore.DeclineRequest")
	defer span.End()

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

func (store *FollowNeo4JStore) SaveAd(ctx context.Context, ad *domain.Ad) error {
	ctx, span := store.tracer.Start(ctx, "FollowStore.SaveAd")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (ad:Ad) SET ad.tweetID = $tweetID, ad.ageFrom = $ageFrom, "+
					"ad.ageTo = $ageTo, ad.gender = $gender, ad.residence = $residence "+
					"RETURN ad.id + ', from node ' + id(ad)",
				map[string]any{"tweetID": ad.TweetID, "ageFrom": ad.AgeFrom, "ageTo": ad.AgeTo,
					"gender": ad.Gender, "residence": ad.Residence})
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

func (store *FollowNeo4JStore) CountFollowings(ctx context.Context, username string) (int, error) {
	ctx, span := store.tracer.Start(ctx, "FollowStore.CountFollowings")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	followingsCount, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (u:User)-[:FOLLOWS]->(u2:User) "+
					"WHERE u.username=$username "+
					"RETURN count(u2) as count",
				map[string]any{"username": username})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				count, _ := result.Record().Get("count")
				return count, nil
			}

			return nil, result.Err()
		})
	if err != nil {
		return 0, err
	}
	log.Printf("COUNT FOLLOWINGS: %s", followingsCount.(int64))
	return int(followingsCount.(int64)), nil
}

func (store *FollowNeo4JStore) RecommendWithFollowings(ctx context.Context, username string) ([]string, error) {
	ctx, span := store.tracer.Start(ctx, "FollowStore.RecommendWithFollowings")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	users, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"OPTIONAL MATCH (u1:User)-[:FOLLOWS]->(u2:User)-[:FOLLOWS]->(u4:User) "+
					"WHERE u1.username = $username AND NOT u1 = u4 AND NOT exists((u1)-[:FOLLOWS]->(u4)) "+
					"OPTIONAL MATCH (u1:User)-[:FOLLOWS]->(u3:User)-[:FOLLOWS]->(u4:User) "+
					"WHERE NOT u1 = u4 AND NOT u2 = u3 "+
					"OPTIONAL MATCH (u4:User)-[:FOLLOWS]->(u5:User) "+
					"WHERE NOT u5 = u2 AND NOT u5 = u3 AND NOT u5 = u1 "+
					"MATCH (u1:User)-[:FOLLOWS]->(u2:User)-[:FOLLOWS]->(u6:User) "+
					"WHERE NOT u6 = u1 AND NOT exists((u1:User)-[:FOLLOWS]->(u6:User)) "+
					"WITH collect(distinct u4.username) + collect(distinct u5.username) + "+
					"collect(distinct u6.username) AS undistUsernames "+
					"UNWIND undistUsernames AS distUsernames "+
					//"RETURN DISTINCT distUsernames as usernames",
					"RETURN collect(DISTINCT distUsernames) as usernames",
				map[string]any{"username": username})
			if err != nil {
				return nil, err
			}

			var users []string
			if result.Next(ctx) {
				usernames, _ := result.Record().Get("usernames")
				if usernames == nil {
					return users, nil
				}
				for _, username := range usernames.([]interface{}) {
					users = append(users, username.(string))
				}
			}

			return users, nil
		})
	if err != nil {
		return nil, err
	}
	log.Println("Recommends with followers: %s", users.([]string))
	return users.([]string), nil
}

func (store *FollowNeo4JStore) RecommendationWithoutFollowings(ctx context.Context, username string, recommends []string) ([]string, error) {
	ctx, span := store.tracer.Start(ctx, "FollowStore.RecommendationWithoutFollowings")
	defer span.End()

	session := store.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: DATABASE})
	defer session.Close(ctx)

	users, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (u1:User), (u2:User) "+
					"WHERE u1.username = $username AND NOT u2.username IN $recommends AND NOT u1 = u2 "+
					"AND u1.residence = u2.residence AND u2.age-3 <= u1.age <= u2.age+3 "+
					"AND NOT exists((u1:User)-[:FOLLOWS]->(u2:User)) "+
					"RETURN collect(u2.username) as usernames",
				map[string]any{"username": username, "recommends": recommends})
			if err != nil {
				return nil, err
			}

			var users []string
			if result.Next(ctx) {
				usernames, ok := result.Record().Get("usernames")
				if usernames == nil {
					log.Printf("OK STATUS IS %s", ok)
					return users, nil
				}
				for _, username := range usernames.([]interface{}) {
					users = append(users, username.(string))
				}
			}

			return users, nil
		})
	if err != nil {
		return nil, err
	}
	log.Println("Recommends without followers: %s", users.([]string))
	return users.([]string), nil
}
