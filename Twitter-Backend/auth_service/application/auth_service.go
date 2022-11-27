package application

import (
	"auth_service/domain"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

var (
	jwtKey = []byte(os.Getenv("SECRET_KEY"))
	//odakle povlazi GetEnv keys?
)

type AuthService struct {
	store domain.AuthStore
}

func NewAuthService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}


func (service *AuthService) Login(user *domain.User) (string, error) {

	user, err := service.store.GetOneUser(user.Username)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Password))

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
		fmt.Println(err) // key is invalid
		return "", err
	}

	service.GetID(service.GetClaims(tokenString))

	return tokenString, nil
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

	fmt.Println("187")
	fmt.Print(claims)
	fmt.Println("189")
	fmt.Println(claims["UserID"])
	userId := claims["UserID"]
	//fmt.Println(userId, claims["Username"].(string))
	return userId.(string)
}
