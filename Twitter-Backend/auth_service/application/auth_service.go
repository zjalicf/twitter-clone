package application

import (
	"auth_service/domain"
	"auth_service/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	userServiceHost = os.Getenv("USER_SERVICE_HOST")
	userServicePort = os.Getenv("USER_SERVICE_PORT")
	smtpServer      = "smtp-mail.outlook.com"
	smtpServerPort  = 587
	smtpEmail       = os.Getenv("SMTP_AUTH_MAIL")
	smtpPassword    = os.Getenv("SMTP_AUTH_PASSWORD")
	jwtKey          = []byte(os.Getenv("SECRET_KEY"))
	//odakle povlazi GetEnv keys?
)

type AuthService struct {
	store domain.AuthStore
	cache domain.AuthCache
}

func NewAuthService(store domain.AuthStore, cache domain.AuthCache) *AuthService {
	return &AuthService{
		store: store,
		cache: cache,
	}
}

func (service *AuthService) GetAll() ([]*domain.User, error) {
	return service.store.GetAll()
}

func (service *AuthService) Register(user *domain.User) (string, int, error) {
	pass := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return "", 500, err
	}
	user.Password = string(hash)

	body, err := json.Marshal(user)
	if err != nil {
		return "", 500, err
	}

	userServiceEndpoint := fmt.Sprintf("http://%s:%s/", userServiceHost, userServicePort)
	userServiceRequest, _ := http.NewRequest("POST", userServiceEndpoint, bytes.NewReader(body))
	responseUser, err := http.DefaultClient.Do(userServiceRequest)
	//if err != nil {
	//	return "",500, err
	//}

	if responseUser.StatusCode != 200 {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, responseUser.Body)
		return "", responseUser.StatusCode, fmt.Errorf(buf.String())
	}

	var newUser domain.User
	err = responseToType(responseUser.Body, &newUser)
	if err != nil {
		return "", 500, err
	}

	credentials := domain.Credentials{
		ID:       newUser.ID,
		Username: user.Username,
		Password: user.Password,
		UserType: newUser.UserType,
	}

	err = service.store.Register(&credentials)
	if err != nil {
		return "", 500, err
	}

	validationToken := uuid.New()
	err = sendValidationMail(validationToken, user.Email)
	if err != nil {
		return "", 500, err
	}

	err = service.cache.PostCacheData(newUser.ID.Hex(), validationToken.String())
	if err != nil {
		log.Fatalf("failed to post validation data to redis: %s", err)
		return "", 500, err
	}

	return newUser.ID.Hex(), 200, nil
}

func sendValidationMail(validationToken uuid.UUID, email string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", smtpEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Verify your Twitter Clone account")

	bodyString := fmt.Sprintf("Your validation token for twitter account is:\n%s", validationToken)
	message.SetBody("text", bodyString)

	client := gomail.NewDialer(smtpServer, smtpServerPort, smtpEmail, smtpPassword)

	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send verification mail because of: %s", err)
		return err
	}

	return nil
}

func (service *AuthService) ValidateAccount(validation *domain.RegisterValidation) error {
	token, err := service.cache.GetCachedValue(validation.UserToken)
	if err != nil {
		log.Fatalf("failed to get value from redis: %s", err)
		return err
	}

	if validation.MailToken == token {
		return nil
	}

	return fmt.Errorf(errors.InvalidTokenError)
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
		UserID:   user.ID,
		Username: user.Username, //menjanje za userID
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

	//service.GetID(service.GetClaims(tokenString))

	return tokenString, nil
}

func responseToType(response io.ReadCloser, user *domain.User) error {
	responseBodyBytes, err := io.ReadAll(response)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBodyBytes, &user)
	if err != nil {
		return err
	}

	return nil
}

// handling token
func (service *AuthService) ValidateJWT(endpoint func(writer http.ResponseWriter, request *http.Request) http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Token"] != nil {
			token, err := jwt.Parse(request.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					writer.Write([]byte("not authorized"))
				}
				return jwtKey, nil

			})

			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				writer.Write([]byte("not authorized"))
			}

			if token.Valid {
				endpoint(writer, request)
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("not authorized"))
		}
	})
}

func (service *AuthService) GetClaims(tokenString string) jwt.MapClaims {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println(ok)
		}
		return token, nil
	})

	if err != nil {
		fmt.Println(err)
		return nil

	}

	return token.Claims.(jwt.MapClaims)
}

func (service *AuthService) GetID(claims jwt.MapClaims) string {

	userId := claims["UserID"]
	//fmt.Println(userId, claims["Username"].(string))
	return userId.(string)
}
