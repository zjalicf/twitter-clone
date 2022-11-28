package application

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"user_service/domain"
	"user_service/errors"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

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

func (service *UserService) Login(user *domain.User) (string, error) {

	returnedUser, err := service.store.GetOneUser(user.Username)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(returnedUser)

	passError := bcrypt.CompareHashAndPassword([]byte(returnedUser.Password), []byte(user.Password))

	if passError != nil {
		fmt.Println(passError)
		return "", fmt.Errorf("Wrong password")
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &domain.Claims{
		UserID: user.ID,
		Role:   user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		fmt.Println(err) // key is invalid
		return "", err
	}

	return tokenString, nil
}

func (service *UserService) Register(user *domain.User) (*domain.User, error) {

	validatedUser, err := validateUserType(user)
	if err != nil {
		return nil, fmt.Errorf(errors.ValidationError)
	}

	retUser, err := service.store.Post(validatedUser)
	if err != nil {
		return nil, fmt.Errorf(errors.DatabaseError)
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
	if len(user.CompanyName) >= 3 &&
		len(user.Website) >= 3 &&
		len(user.Email) >= 3 &&
		len(user.Username) >= 3 &&
		len(user.Password) >= 8 {
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
		len(user.Username) >= 3 &&
		len(user.Password) >= 8 {
		return true
	}

	return false
}
