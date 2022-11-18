package application

import (
	"auth_service/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
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
	pass := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	userServiceEndpoint := fmt.Sprintf("http://%s:%s/", userServiceHost, userServicePort)

	userServiceRequest, _ := http.NewRequest("POST", userServiceEndpoint, bytes.NewReader(body))
	responseUser, err := http.DefaultClient.Do(userServiceRequest)
	if err != nil {
		return err
	}

	response, err := io.ReadAll(responseUser.Body)
	if err != nil {
		return err
	}

	var newUser domain.User
	json.Unmarshal(response, &newUser)

	credentials := domain.Credentials{
		ID:       newUser.ID,
		Username: user.Username,
		Password: user.Password,
		UserType: newUser.UserType,
	}

	return service.store.Register(&credentials)
}

func (service *AuthService) Login(credentials *domain.Credentials) error {
	return nil
}
