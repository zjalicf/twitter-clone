package application

import (
	"fmt"
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

func (service *UserService) Post(user *domain.User) (*domain.User, error) {
	user.ID = primitive.NewObjectID()

	validatedUser, err := validateUserType(user)
	if err != nil {
		return nil, err
	}

	retUser, err := service.store.Post(validatedUser)
	if err != nil {
		return nil, err
	}

	return retUser, nil
}

func validateUserType(user *domain.User) (*domain.User, error) {

	business := isBusiness(user)
	regular := isRegular(user)

	if business && regular {
		return nil, fmt.Errorf("invalid user format")
	} else if business {
		user.UserType = domain.Business
		return user, nil
	} else if regular {
		user.UserType = domain.Regular
		return user, nil
	}

	return nil, fmt.Errorf("invalid user data")
}

func isBusiness(user *domain.User) bool {
	if len(user.CompanyName) != 0 && len(user.Website) != 0 && len(user.Email) != 0 &&
		len(user.Username) != 0 && len(user.Password) != 0 {
		return true
	}

	return false
}

func isRegular(user *domain.User) bool {
	if len(user.Firstname) != 0 && len(user.Lastname) != 0 &&
		len(user.Gender) != 0 && user.Age <= 0 &&
		len(user.Residence) != 0 && len(user.Email) != 0 &&
		len(user.Username) != 0 && len(user.Password) != 0 {
		return true
	}

	return false
}
