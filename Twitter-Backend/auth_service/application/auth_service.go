package application

type AuthService struct {
	store domain.AuthStore	
}

func NewAuthService(store domain.AuthService) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (service *TweetService) Get(id primitive.ObjectID) (*domain.Tweet, error) {
	return service.store.Get(id)
}
