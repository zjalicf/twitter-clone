package application

import (
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
	"log"
)

type CreateEventOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreateEventOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*CreateEventOrchestrator, error) {
	orchestrator := &CreateEventOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := orchestrator.replySubscriber.Subscribe(orchestrator.handle)
	if err != nil {
		return nil, err
	}
	return orchestrator, nil
}

func (o *CreateEventOrchestrator) Start(event events.Event) error {

	eventPush := events.CreateEventCommand{
		Event: event,
		Type:  events.UpdateCassandra,
	}

	return o.commandPublisher.Publish(eventPush)
}

func (o *CreateEventOrchestrator) handle(reply *events.CreateEventReply) {
	command := events.CreateEventCommand{Event: reply.Event}
	command.Type = o.nextCommandType(*reply)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *CreateEventOrchestrator) nextCommandType(reply events.CreateEventReply) events.CreateEventCommandType {
	switch reply.Type {
	case events.MongoUpdated:
		log.Println("MONGO UPDATED")
		return events.UnknownCommand
	case events.CassandraUpadated:
		log.Println("CASSANDRA UPDATED")
		return events.UpdateMongo
	case events.UnknownReply:
		return events.UnknownCommand
	default:
		return events.UnknownCommand
	}
}
