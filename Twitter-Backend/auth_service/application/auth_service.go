package application

import (
	"auth_service/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	userServiceHost = os.Getenv("USER_SERVICE_HOST")
	userServicePort = os.Getenv("USER_SERVICE_PORT")
)

var jwtKey = []byte("my_secret_key")

type AuthService struct {
	store domain.AuthStore
}

func NewAuthService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (service *AuthService) Register(user *domain.User) (int, error) {
	pass := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return 500, err
	}
	user.Password = string(hash)

	body, err := json.Marshal(user)
	if err != nil {
		return 500, err
	}

	userServiceEndpoint := fmt.Sprintf("http://%s:%s/", userServiceHost, userServicePort)

	userServiceRequest, _ := http.NewRequest("POST", userServiceEndpoint, bytes.NewReader(body))
	responseUser, err := http.DefaultClient.Do(userServiceRequest)
	if err != nil {
		return 500, err
	}

	if responseUser.StatusCode != 200 {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, responseUser.Body)
		return responseUser.StatusCode, fmt.Errorf(buf.String())
	}

	var newUser domain.User
	err = responseToType(responseUser.Body, newUser)
	if err != nil {
		return 500, err
	}

	credentials := domain.Credentials{
		ID:       newUser.ID,
		Username: user.Username,
		Password: user.Password,
		UserType: newUser.UserType,
	}

	err = service.store.Register(&credentials)
	if err != nil {
		return 500, err
	}
	return 200, nil
}

func (service *AuthService) Login(credentials *domain.Credentials) (string, error) {

	user, err := service.store.GetOneUser(credentials.Username)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))

	if passError != nil {
		fmt.Println(passError)
		return "", err
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &domain.Claims{
		Username: user.Username,
		Role:     user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}

func responseToType(response io.ReadCloser, any any) error {
	responseBodyBytes, err := io.ReadAll(response)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBodyBytes, &any)
	if err != nil {
		return err
	}

	return nil
}
