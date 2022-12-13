package domain

type FollowRequestStore interface {
	GetAll() ([]FollowRequest, error)
	SendRequest() (*FollowRequest, error)
	HandleRequest()
	DeclineRequest(*FollowRequest) (int, error)
}
