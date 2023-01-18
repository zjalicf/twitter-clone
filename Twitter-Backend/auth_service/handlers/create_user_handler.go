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
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(authService *application.AuthService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		authService:       authService,
		replyPublisher:    publisher,
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
	user := handler.authService.UserToDomain(command.User)
	reply := events.CreateUserReply{User: command.User}

	switch command.Type {
	case events.UpdateAuth:
		_, _, err := handler.authService.Register(nil, &user)
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

//hendlovanje odgovorima
func (handler *CreateUserCommandHandler) handleReplays(reply *events.CreateUserReply) {

	fmt.Println(reply)

	switch reply.Type {
	case events.UsersUpdated:

		//posto se update radi u dva servisa ,trebalo bi imati neki brojac koji ce da broji broj pristiglih poruka. Tipa prva iz user_Service a druga iz follow servisa
		//i samo u tom slucaju saga je uspesna

		//poslati mejl kada stigne ova poruka
		user := handler.authService.UserToDomain(reply.User)
		err := handler.authService.SendMail(&user)
		if err != nil {
			log.Printf("Failed to send mail: %s", err.Error())
			return
		}
	case events.GraphUpdated:
		fmt.Println("Napravljen NODE")

	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
