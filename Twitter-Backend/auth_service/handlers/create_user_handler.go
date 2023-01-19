package handlers

import (
	"auth_service/application"
	"fmt"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_user"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
)

type CreateUserCommandHandler struct {
	authService       *application.AuthService
	publisher         saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(authService *application.AuthService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		authService:       authService,
		publisher:         publisher,
		commandSubscriber: subscriber,
	}
	//prijava za slusanje komandi
	err := o.commandSubscriber.Subscribe(o.handleCommands)
	if err != nil {
		return nil, err
	}

	//prijava za slusanje odgovora
	err1 := o.commandSubscriber.Subscribe(o.handleReplays)
	if err1 != nil {
		return nil, err1
	}
	return o, nil
}

//hendlovanje komandama
func (handler *CreateUserCommandHandler) handleCommands(command *events.CreateUserCommand) {
	//user := handler.authService.UserToDomain(command.User)
	reply := events.CreateUserReply{User: command.User}

	switch command.Type {
	case events.UpdateAuth:
		fmt.Println("Stigla poruka u auth")
		//_, _, err := handler.authService.Register(nil, &user)
		//if err != nil {
		//	return
		//}

		reply.Type = events.AuthUpdated

	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.publisher.Publish(reply)
	}
}

//hendlovanje odgovorima
func (handler *CreateUserCommandHandler) handleReplays(reply *events.CreateUserReply) {

	switch reply.Type {
	case events.UsersUpdated:
		fmt.Println("Gotov user update")
		//user := handler.authService.UserToDomain(reply.User)
		//err := handler.authService.SendMail(&user)
		//if err != nil {
		//	log.Printf("Failed to send mail: %s", err.Error())
		//	return
		//}

	case events.GraphUpdated:
		fmt.Println("Gotova saga")

	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.publisher.Publish(reply)
	}
}
