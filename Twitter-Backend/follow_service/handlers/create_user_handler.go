package handlers

import (
	"fmt"
	"follow_service/application"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_user"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
)

type CreateUserCommandHandler struct {
	followService     *application.FollowService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(followService *application.FollowService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		followService:     followService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateUserCommandHandler) handle(command *events.CreateUserCommand) {
	user := handler.followService.UserToDomain(command.User)
	reply := events.CreateUserReply{User: command.User}

	switch command.Type {

	case events.UpdateGraph:
		err := handler.followService.CreateUser(&user)
		if err != nil {
			reply.Type = events.FollowFailed
		} else {
			reply.Type = events.GraphUpdated
		}

	case events.RollbackFollow:
		//TODO
		fmt.Println("Rollback follow")
		reply.Type = events.FollowFailed
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
