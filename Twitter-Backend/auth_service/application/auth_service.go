package application

import (
	"auth_service/domain"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type AuthService struct {
	store domain.AuthStore
}

func NewAuthService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
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
