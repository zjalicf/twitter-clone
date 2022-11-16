package application

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user_service/domain"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) Get(id primitive.ObjectID) (*domain.User, error) {
	return service.store.Get(id)
}

func (service *UserService) GetAll() ([]*domain.User, error) {
	return service.store.GetAll()
}

func (service *UserService) Post(user *domain.User) error {
	user.ID = primitive.NewObjectID()
	return service.store.Post(user)
}
