package application

import (
	"fmt"
	"follow_service/domain"
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

func (service *FollowService) GetRequestsForUser(username string) ([]*domain.FollowRequest, error) {
	return service.store.GetRequestsForUser(username)
}

func (service *FollowService) CreateRequest(request *domain.FollowRequest, username string) (*domain.FollowRequest, error) {

	request.ID = primitive.NewObjectID()
	request.Requester = username
	request.Status = 1

	retFollow, err := service.store.SaveRequest(request)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Follow not inserted in db")
	}
	return retFollow, nil
}

func (service *FollowService) AcceptRequest(id primitive.ObjectID) error {
	return service.store.AcceptRequest(id)
}

func (service *FollowService) DeclineRequest(id primitive.ObjectID) error {
	return service.store.DeclineRequest(id)
}

func (service *FollowService) HandleRequest(followRequest *domain.FollowRequest) (*domain.FollowRequest, error) {
	return service.store.SaveRequest(followRequest)
}
