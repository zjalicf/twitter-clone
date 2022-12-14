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

func (service *FollowService) CreateRequest(request *domain.FollowRequest) (*domain.FollowRequest, error) {

	//todo insert in follow_db

	request.ID = primitive.NewObjectID()

	//pozivanje upisa ka bazi
	retFollow, err := service.store.SaveRequest(request)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Follow not inserted in db")
	}
	return retFollow, nil
}

func (service *FollowService) SendRequest() (*domain.FollowRequest, error) {
	return service.store.SendRequest()
}

func (service *FollowService) DeclineRequest(followRequest *domain.FollowRequest) (*domain.FollowRequest, error) {
	return service.store.DeclineRequest(followRequest)
}

func (service *FollowService) HandleRequest(followRequest *domain.FollowRequest) (*domain.FollowRequest, error) {
	return service.store.SaveRequest(followRequest)
}
