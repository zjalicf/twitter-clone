package application

import (
	"auth_service/authorization"
	"auth_service/domain"
	"auth_service/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	jwtKey          = []byte("SecretYouShouldHide")
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
	responseUser, _ := http.DefaultClient.Do(userServiceRequest)
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
	err = service.cache.PostCacheData(newUser.ID.Hex(), validationToken.String())
	if err != nil {
		log.Fatalf("failed to post validation data to redis: %s", err)
		return "", 500, err
	}

	err = sendValidationMail(validationToken, user.Email)
	if err != nil {
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

func (service *AuthService) VerifyAccount(validation *domain.RegisterRecoverVerification) error {
	token, err := service.cache.GetCachedValue(validation.UserToken)
	if err != nil {
		log.Println(errors.ExpiredTokenError)
		return fmt.Errorf(errors.ExpiredTokenError)
	}

	if validation.MailToken == token {
		err := service.cache.DelCachedValue(validation.UserToken)
		if err != nil {
			log.Printf("error in deleting cached value: %s", err)
			return err
		}
		return nil
	}

	return fmt.Errorf(errors.InvalidTokenError)
}

func (service *AuthService) ResendVerificationToken(request *domain.ResendVerificationRequest) error {

	if len(request.UserMail) == 0 {
		log.Println(errors.InvalidResendMailError)
		return fmt.Errorf(errors.InvalidResendMailError)
	}

	tokenUUID, _ := uuid.NewUUID()

	err := service.cache.PostCacheData(request.UserToken, tokenUUID.String())
	if err != nil {
		return err
	}

	err = sendValidationMail(tokenUUID, request.UserMail)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) SendRecoveryPasswordToken(email string) (string, int, error) {

	userServiceEndpoint := fmt.Sprintf("http://%s:%s/mailExist/%s", userServiceHost, userServicePort, email)
	userServiceRequest, _ := http.NewRequest("GET", userServiceEndpoint, nil)
	response, _ := http.DefaultClient.Do(userServiceRequest)
	if response.StatusCode != 200 {
		if response.StatusCode == 404 {
			return "", 404, fmt.Errorf(errors.NotFoundMailError)
		}
	}

	buf := new(strings.Builder)
	_, _ = io.Copy(buf, response.Body)
	userID := buf.String()

	recoverUUID, _ := uuid.NewUUID()
	err := sendRecoverPasswordMail(recoverUUID, email)
	if err != nil {
		return "", 500, err
	}

	err = service.cache.PostCacheData(userID, recoverUUID.String())
	if err != nil {
		return "", 500, err
	}

	return userID, 200, nil
}

func (service *AuthService) CheckRecoveryPasswordToken(request *domain.RegisterRecoverVerification) error {

	if len(request.UserToken) == 0 {
		return fmt.Errorf(errors.InvalidUserTokenError)
	}

	token, err := service.cache.GetCachedValue(request.UserToken)
	if err != nil {
		return fmt.Errorf(errors.InvalidTokenError)
	}

	if request.MailToken != token {
		return fmt.Errorf(errors.InvalidTokenError)
	}

	_ = service.cache.DelCachedValue(request.UserToken)
	return nil
}

func sendRecoverPasswordMail(validationToken uuid.UUID, email string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", smtpEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Recover password on your Twitter Clone account")

	bodyString := fmt.Sprintf("Your recover password token is:\n%s", validationToken)
	message.SetBody("text", bodyString)

	client := gomail.NewDialer(smtpServer, smtpServerPort, smtpEmail, smtpPassword)

	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send verification mail because of: %s", err)
		return err
	}

	return nil
}

func (service *AuthService) RecoverPassword(recoverPassword *domain.RecoverPasswordRequest) error {
	if recoverPassword.NewPassword != recoverPassword.RepeatedNew {
		return fmt.Errorf(errors.NotMatchingPasswordsError)
	}

	primitiveID, err := primitive.ObjectIDFromHex(recoverPassword.UserID)
	if err != nil {
		return err
	}
	credentials := service.store.GetOneUserByID(primitiveID)

	pass := []byte(recoverPassword.NewPassword)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	credentials.Password = string(hash)

	err = service.store.ChangePassword(credentials)
	if err != nil {
		return err
	}

	return nil
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
		return "not_same", err
	}

	tokenString, err := GenerateJWT(user)

	if err != nil {
		return "", err
	}

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

func GenerateJWT(user *domain.User) (string, error) {

	key := []byte(os.Getenv("SECRET_KEY"))
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		log.Println(err)
	}

	builder := jwt.NewBuilder(signer)

	claims := &domain.Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.UserType,
		ExpiresAt: time.Now().Add(time.Minute * 60),
	}

	token, err := builder.Build(claims)
	if err != nil {
		log.Println(err)
	}

	return token.String(), nil
}

func (service *AuthService) ChangePassword(password domain.PasswordChange, token string) error {

	parsedToken := authorization.GetToken(token)
	claims := authorization.GetMapClaims(parsedToken.Bytes())

	username := claims["username"]

	fmt.Println(username)

	user, err := service.store.GetOneUser(username)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Old password: %s", user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.OldPassword))
	if err != nil {
		return err
	}

	var validNew bool = false
	if password.NewPassword == password.NewPasswordConfirm {
		validNew = true
	}

	if validNew {
		newEncryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Printf("New password: %s", newEncryptedPassword)

		user.Password = string(newEncryptedPassword)

		err = service.store.ChangePassword(user)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("new password not match")
	}

	return nil
}
