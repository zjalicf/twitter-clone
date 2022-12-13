package application

import (
	"follow_service/domain"
)

type FollowService struct {
	store domain.FollowRequestStore
}

func NewFollowService(store domain.FollowRequestStore) *FollowService {
	return &FollowService{
		store: store,
	}
}

func (service *FollowService) GetAll() ([]domain.FollowRequest, error) {
	return service.store.GetAll()
}

func (service *FollowService) SendRequest() (*domain.FollowRequest, error) {
	return service.store.SendRequest()
}

func (service *FollowService) DeclineRequest() (int, error) {
	return service.store.DeclineRequest()
}

func (service *FollowService) HandleRequest(followRequest *domain.FollowRequest) (*domain.FollowRequest, error) {
	//todo
	return service.store.Post(followRequest)
}
