package handlers

import (
	"auth_service/application"
	"fmt"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_user"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
	"log"
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
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}

	//prijava za slusanje odgovora
	//err1 := o.replySubscriber.Subscribe(o.handleReplays)
	//if err1 != nil {
	//	return nil, err1
	//}
	return o, nil
}

//hendlovanje komandama
func (handler *CreateUserCommandHandler) handle(command *events.CreateUserCommand) {
	user := handler.authService.UserToDomain(command.User)
	reply := events.CreateUserReply{User: command.User}

	switch command.Type {

	case events.UpdateAuth:
		log.Println("Uslo u update auth")
		reply.Type = events.AuthUpdated

	case events.SendMail:
		err := handler.authService.SendMail(&user)
		if err != nil {
			log.Printf("Failed to send mail: %s", err.Error())
			reply.Type = events.MailFailed
		} else {
			reply.Type = events.MailSent

		}

	case events.RollbackAuth:
		//TODO
		reply.Type = events.UnknownReply
		fmt.Println("Rollback auth")

	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.publisher.Publish(reply)
	}
}
