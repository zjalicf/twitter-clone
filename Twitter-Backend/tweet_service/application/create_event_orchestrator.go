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
		log.Printf("Error in create orchestrator NewCreateEventOrchestrator(): %s", err.Error())
		return nil, err
	}
	return orchestrator, nil
}

func (o *CreateEventOrchestrator) Start(event events.Event) error {

	log.Printf("Starting orchestrator with event : %s", event)

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
		log.Println(command.Type)
	}
	return
}

func (o *CreateEventOrchestrator) nextCommandType(reply events.CreateEventReply) events.CreateEventCommandType {

	log.Println(reply)

	switch reply.Type {
	case events.CassandraUpadated:
		log.Println("CassandraUpdated")
		return events.UpdateMongo

	case events.MongoUpdated:
		log.Println("MongoUpdated")
		return events.UnknownCommand

	default:
		log.Println("unknown command nextCommandType orchestrator")
		return events.UnknownCommand
	}
}
