package domain

type FollowRequestStore interface {
	GetAll() ([]*FollowRequest, error)
	SaveRequest(*FollowRequest) (*FollowRequest, error)
	AcceptRequest() (*FollowRequest, error)
	DeclineRequest(*FollowRequest) (*FollowRequest, error)
	HandleRequest()
}
