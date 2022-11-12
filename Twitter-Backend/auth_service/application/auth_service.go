package application

import (
	"auth_service/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	store domain.AuthStore
}

func RegisterService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (service *AuthService) Register(user *domain.User, isBusiness bool) error {
	user.ID = primitive.NewObjectID()

	if isBusiness {
		user.UserType = domain.Business
	} else {
		user.UserType = domain.Regular
	}

	//userService todo
	// user *domain.User -> credentials

	return service.store.Register(user, isBusiness)
}
