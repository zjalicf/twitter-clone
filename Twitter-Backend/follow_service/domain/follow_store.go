package domain

type FollowRequestStore interface {
	GetAll() ([]*FollowRequest, error)
	SaveUser(*User) error
	DeleteUser(id *string) error
	GetRequestsForUser(username string) ([]*FollowRequest, error)
	GetRequestByRequesterReceiver(requester, receiver *string) (*FollowRequest, error)
	GetFollowingsOfUser(username string) ([]*User, error)
	SaveRequest(*FollowRequest) error
	SaveFollow(request *FollowRequest) error
	AcceptRequest(id *string) (*FollowRequest, error)
	DeclineRequest(id *string) error
	FollowExist(followRequest *FollowRequest) (bool, error)
	UpdateRequest(request *FollowRequest) error
	SaveAd(ad *Ad) error
	CountFollowings(username string) (int, error)
	RecommendWithFollowings(username string) ([]string, error)
	RecommendationWithoutFollowings(username string, recommends []string) ([]string, error)
}
