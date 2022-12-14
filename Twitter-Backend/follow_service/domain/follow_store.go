package domain

type FollowRequestStore interface {
	GetAll() ([]*FollowRequest, error)
	SaveRequest(*FollowRequest) (*FollowRequest, error)
	SendRequest() (*FollowRequest, error)
	DeclineRequest(*FollowRequest) (*FollowRequest, error)
	HandleRequest()
}
