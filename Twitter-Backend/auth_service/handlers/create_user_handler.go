package handlers

import (
	"auth_service/application"
	events "github.com/tamararankovic/microservices_demo/common/saga/create_order"
	saga "github.com/tamararankovic/microservices_demo/common/saga/messaging"
)

type CreateUserCommandHandler struct {
	authService       *application.AuthService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(authService *application.AuthService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		authService:       authService,
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

	reply := events.CreateUserReply{User: command.User}

	switch command.Type {
	case events.UpdateAuth:
		_, _, err := handler.authService.Register(user)
		if err != nil {
			return
		}
		reply.Type = events.AuthUpdated
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
