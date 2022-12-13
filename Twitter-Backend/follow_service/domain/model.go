package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowRequest struct {
	ID        primitive.ObjectID `json:"id"`
	Receiver  primitive.ObjectID `json:"receiver_id"`
	Requester primitive.ObjectID `json:"requester_id"`
	Status    Status             `json:"status"`
}

type Status int

const (
	Pending Status = iota
	Declined
	Accepted
)
