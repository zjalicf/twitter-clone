package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type FollowRequestStore interface {
	GetAll() ([]*FollowRequest, error)
	SaveRequest(*FollowRequest) (*FollowRequest, error)
	AcceptRequest(id primitive.ObjectID) error
	DeclineRequest(id primitive.ObjectID) error
	HandleRequest()
}
