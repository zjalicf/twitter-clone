package application

import (
	"fmt"
	"github.com/zjalicf/twitter-clone-common/common/saga/create_user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"user_service/domain"
	"user_service/errors"
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

func (service *UserService) DoesEmailExist(email string) (string, error) {
	user, err := service.store.GetByEmail(email)
	if err != nil {
		return "", err
	}

	return user.ID.Hex(), nil
}

func (service *UserService) GetAll() ([]*domain.User, error) {
	return service.store.GetAll()
}

func (service *UserService) GetOneUser(username string) (*domain.User, error) {
	retUser, err := service.store.GetOneUser(username)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("User not found")
	}
	return retUser, nil
}

func (service *UserService) Register(user *domain.User) (*domain.User, error) {
	fmt.Println(user)
	validatedUser, err := validateUserType(user)
	if err != nil {
		log.Println(errors.ValidationError)
		return nil, fmt.Errorf(errors.ValidationError)
	}
	validatedUser.Visibility = true

	retUser, err := service.store.Post(validatedUser)
	if err != nil {
		log.Println(errors.DatabaseError)
		return nil, fmt.Errorf(errors.DatabaseError)
	}

	return retUser, nil
}

func (service *UserService) ChangeUserVisibility(userID string) error {
	primitiveID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Primitive ID parsing error.")
		return err
	}

	user, err := service.store.Get(primitiveID)
	if err != nil {
		log.Printf("Getting user by id error: %s", err.Error())
		return fmt.Errorf(errors.UserNotFound)
	}

	user.Visibility = !user.Visibility
	err = service.store.UpdateUser(user)
	if err != nil {
		log.Printf("Updating user error in service: %s", err.Error())
		return err
	}

	return nil
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
	if len(user.CompanyName) >= 3 &&
		len(user.Website) >= 3 &&
		len(user.Email) >= 3 &&
		len(user.Username) >= 3 {
		return true
	}

	return false
}

func isRegular(user *domain.User) bool {
	if len(user.Firstname) >= 3 &&
		len(user.Lastname) >= 3 &&
		len(user.Gender) >= 3 &&
		user.Age >= 1 &&
		len(user.Residence) >= 3 &&
		len(user.Username) >= 3 {
		return true
	}

	return false
}

func (service *UserService) UserToDomain(userIn create_user.User) domain.User {
	var user domain.User
	user.ID = userIn.ID
	user.Firstname = userIn.Firstname
	user.Lastname = userIn.Lastname
	if userIn.Gender == "Male" {
		user.Gender = "Male"
	} else {
		user.Gender = "Female"
	}
	user.Age = userIn.Age
	user.Residence = userIn.Residence
	user.Email = userIn.Email
	user.Username = userIn.Username
	if userIn.UserType == "Regular" {
		user.UserType = "Regular"
	} else {
		user.UserType = "Business"
	}
	user.Visibility = userIn.Visibility
	user.CompanyName = userIn.CompanyName
	user.Website = userIn.Website

	return user
}
