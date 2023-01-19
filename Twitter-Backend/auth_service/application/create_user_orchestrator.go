package application

import (
	"auth_service/domain"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_user"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
)

type CreateUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserOrchestrator, error) {
	orchestrator := &CreateUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := orchestrator.replySubscriber.Subscribe(orchestrator.handle)
	if err != nil {
		return nil, err
	}
	return orchestrator, nil
}

func (o *CreateUserOrchestrator) Start(user *domain.User) error {

	var gender events.Gender
	var userType events.UserType

	if user.Gender == "Male" {
		gender = "Male"
	} else {
		gender = "Female"
	}

	if user.UserType == "Regular" {
		userType = "Regular"
	} else {
		userType = "Business"
	}

	user1 := events.User{
		ID:          user.ID,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Gender:      gender,
		Age:         user.Age,
		Residence:   user.Residence,
		Email:       user.Email,
		Username:    user.Username,
		Password:    user.Password,
		UserType:    userType,
		Visibility:  user.Visibility,
		CompanyName: user.CompanyName,
		Website:     user.Website,
	}

	event := &events.CreateUserCommand{
		User: user1,
		Type: events.UpdateAuth,
	}
	return o.commandPublisher.Publish(event)
}

func (o *CreateUserOrchestrator) handle(reply *events.CreateUserReply) {
	command := events.CreateUserCommand{User: reply.User}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *CreateUserOrchestrator) nextCommandType(reply events.CreateUserReplyType) events.CreateUserCommandType {
	switch reply {
	case events.AuthUpdated:
		return events.UpdateUsers
	case events.UsersUpdated:
		return events.UpdateGraph
	default:
		return events.UnknownCommand
	}
}
