package application

import (
	"auth_service/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

var (
	userServiceHost = os.Getenv("USER_SERVICE_HOST")
	userServicePort = os.Getenv("USER_SERVICE_PORT")
)

type AuthService struct {
	store domain.AuthStore
}

func NewAuthService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (service *AuthService) Register(user *domain.User) error {
	validatedUser, err := validateUserType(user)
	if err != nil {
		return err
	}

	pass := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	body, err := json.Marshal(validatedUser)
	if err != nil {
		return err
	}

	userServiceEndpoint := fmt.Sprintf("http://%s:%s/", userServiceHost, userServicePort)

	userServiceRequest, _ := http.NewRequest("POST", userServiceEndpoint, bytes.NewReader(body))
	_, err = http.DefaultClient.Do(userServiceRequest)

	if err != nil {
		return err
	}

	credentials := domain.Credentials{Username: user.Username, Password: user.Password, UserType: user.UserType}
	credentials.ID = primitive.NewObjectID()

	return service.store.Register(&credentials)
}

func (service *AuthService) Login(credentials *domain.Credentials) error {
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
	}

	return nil, fmt.Errorf("invalid user data")
}

func isBusiness(user *domain.User) bool {
	if len(user.CompanyName) != 0 && len(user.Email) != 0 && len(user.WebSite) != 0 {
		return false
	}

	return true
}

func isRegular(user *domain.User) bool {
	if len(user.FirstName) != 0 && len(user.LastName) != 0 && len(user.Username) != 0 && len(user.Gender) != 0 {
		return false
	}

	return true
}
