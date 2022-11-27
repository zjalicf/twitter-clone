package application

import (
	"auth_service/domain"
)

type AuthService struct {
	store domain.AuthStore
}

func NewAuthService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (service *AuthService) Login(credentials *domain.Credentials) error {
	return nil
}
