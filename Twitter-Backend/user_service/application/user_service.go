package application

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"user_service/domain"
	"user_service/errors"
)

var (
	smtpServer     = "smtp-mail.outlook.com"
	smtpServerPort = 587
	smtpEmail      = os.Getenv("SMTP_AUTH_MAIL")
	smtpPassword   = os.Getenv("SMTP_AUTH_PASSWORD")
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
		log.Println(errors.ValidationError)
		return nil, fmt.Errorf(errors.ValidationError)
	}

	retUser, err := service.store.Post(validatedUser)
	if err != nil {
		log.Println(errors.DatabaseError)
		return nil, fmt.Errorf(errors.DatabaseError)
	}

	err = sendValidationMail(validatedUser.Email)
	if err != nil {
		log.Println(err)
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
		len(user.Gender) != 0 && user.Age > 0 &&
		len(user.Residence) != 0 && len(user.Email) != 0 &&
		len(user.Username) != 0 && len(user.Password) != 0 {
		return true
	}

	return false
}

func sendValidationMail(email string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", smtpEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Verify your Twitter Clone account")
	validationID := uuid.New()
	bodyString := fmt.Sprintf("Your validation token for twitter account is: \n %s", validationID)
	message.SetBody("text", bodyString)

	client := gomail.NewDialer(smtpServer, smtpServerPort, smtpEmail, smtpPassword)

	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send verification mail because of: %s", err)
		return err
	}

	return nil
}
