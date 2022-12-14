package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Receiver  string             `bson:"receiver" json:"receiver"`
	Requester string             `bson:"requester" json:"requester"`
	Status    Status             `bson:"status" json:"status,omitempty"`
}

type Status int

const (
	Pending Status = iota + 1
	Declined
	Accepted
)

func (status Status) String() string {
	return [...]string{"Pending", "Declined", "Accepted"}[status-1]
}

func (status Status) EnumIndex() int {
	return int(status)
}
