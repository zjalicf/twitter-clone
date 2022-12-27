package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type FollowRequestStore interface {
	GetAll() ([]*FollowRequest, error)
	GetRequestsForUser(username string) ([]*FollowRequest, error)
	GetFollowingsOfUser(username string) ([]*FollowRequest, error)
	SaveRequest(*FollowRequest) (*FollowRequest, error)
	AcceptRequest(id primitive.ObjectID) error
	DeclineRequest(id primitive.ObjectID) error
	HandleRequest()
	FollowExist(followRequest *FollowRequest) (bool, error)
}
