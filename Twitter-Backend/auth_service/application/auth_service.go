package application

import (
	"auth_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
)

var (
	user_service_port = os.Getenv("USER_SERVICE_PORT")
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
	//if isBusiness {
	//	user.UserType = domain.Business
	//} else {
	//	user.UserType = domain.Regular
	//}

	//body, err := json.Marshal(user)
	//if err != nil {
	//	return err
	//}

	//uServRequest, _ := http.NewRequest("POST", "http://localhost:"+user_service_port, bytes.NewReader(body))
	//client := &http.Client{}
	//userResponse, err := client.Do(uServRequest)

	//if err != nil {
	//	return err
	//}
	//
	//defer userResponse.Body.Close()

	credentials := domain.Credentials{Username: user.Username, Password: user.Password, UserType: user.UserType}
	credentials.ID = primitive.NewObjectID()

	return service.store.Register(&credentials)
}

func (service *AuthService) Login(credentials *domain.Credentials) error {
	return nil
}
