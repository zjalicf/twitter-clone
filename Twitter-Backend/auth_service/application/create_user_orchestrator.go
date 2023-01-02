package application

import (
	"auth_service/domain"
	events "github.com/tamararankovic/microservices_demo/common/saga/create_order"
	saga "github.com/tamararankovic/microservices_demo/common/saga/messaging"
)

type CreateUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserOrchestrator, error) {
	orchestrator := &CreateOrderOrchestrator{
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
	event := &events.CreateUserCommand{
		Type: events.UpdateInventory,
		User: user,
	}
	//for _, item := range user.Items {
	//	eventItem := events.OrderItem{
	//		Product: events.Product{
	//			Id:    item.Product.Id,
	//			Color: events.Color{Code: item.Product.Color.Code},
	//		},
	//		Quantity: item.Quantity,
	//	}
	//	event.Order.Items = append(event.Order.Items, eventItem)
	//}

	//upis u bazu treba
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
	case events.UsersUpdated:
		return events.UpdateGraph
	default:
		return events.UnknownCommand
	}
}
