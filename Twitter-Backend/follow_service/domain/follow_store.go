package domain

type FollowRequestStore interface {
	GetAll() ([]*FollowRequest, error)
	SaveUser(*User) error
	GetRequestsForUser(username string) ([]*FollowRequest, error)
	GetFollowingsOfUser(username string) ([]*User, error)
	SaveRequest(*FollowRequest) error
	SaveFollow(requestID *string) error
	AcceptRequest(id *string) error
	DeclineRequest(id *string) error
	//HandleRequest()
	FollowExist(followRequest *FollowRequest) (bool, error)
}
