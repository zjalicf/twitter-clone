package application

import (
	"fmt"
	"follow_service/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type FollowService struct {
	store domain.FollowRequestStore
}

func NewFollowService(store domain.FollowRequestStore) *FollowService {
	return &FollowService{
		store: store,
	}
}

func (service *FollowService) GetAll() ([]*domain.FollowRequest, error) {
	return service.store.GetAll()
}

func (service *FollowService) FollowExist(followRequest *domain.FollowRequest) (bool, error) {
	return service.store.FollowExist(followRequest)
}

func (service *FollowService) GetFollowingsOfUser(username string) ([]*string, error) {
	followings, err := service.store.GetFollowingsOfUser(username)
	if err != nil {
		return nil, err
	}

	var usernameList []*string
	for i := 0; i < len(followings); i++ {
		usernameList = append(usernameList, &followings[i].Receiver)
	}
	log.Println("LIST OF USERNAMES FOLLOW SERVICE: ")
	log.Println(usernameList)
	usernameList = append(usernameList, &username)
	return usernameList, nil
}

func (service *FollowService) GetRequestsForUser(username string) ([]*domain.FollowRequest, error) {
	return service.store.GetRequestsForUser(username)
}

func (service *FollowService) CreateRequest(request *domain.FollowRequest, username string, visibility bool) error {

	request.ID = uuid.New().String()
	request.Requester = username

	if visibility {
		request.Status = 1
	} else {
		request.Status = 3
	}

	isExist, err := service.FollowExist(request)
	if err != nil {
		return err
	}

	if isExist {
		return fmt.Errorf("You already follow this user!")
	}

	err = service.store.SaveRequest(request)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Follow not inserted in db")
	}

	return nil
}

func (service *FollowService) CreateUser(user *domain.User) error {
	err := service.store.SaveUser(user)
	if err != nil {
		log.Printf("Error with saving user node in neo4j: %s", err.Error())
		return err
	}

	return nil
}

func (service *FollowService) AcceptRequest(id primitive.ObjectID) error {
	return service.store.AcceptRequest(id)
}

func (service *FollowService) DeclineRequest(id primitive.ObjectID) error {
	return service.store.DeclineRequest(id)
}

func (service *FollowService) HandleRequest(followRequest *domain.FollowRequest) error {
	return service.store.SaveRequest(followRequest)
}
