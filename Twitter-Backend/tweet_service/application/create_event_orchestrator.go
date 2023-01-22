package application

import (
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
	"log"
	"tweet_service/domain"
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

func (o *CreateEventOrchestrator) Start(event *domain.Event) error {

	var eventOut events.Event
	eventOut.Type = event.Type
	eventOut.TweetID = event.TweetID
	eventOut.Timestamp = event.Timestamp

	eventPush := events.CreateEventCommand{
		Event: eventOut,
		Type:  events.UpdateMongo,
	}
	log.Println("PUBLISH EVENT UPDATE MONGO")

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
		return events.UpdateCassandra
	case events.CassandraUpadated:
		log.Println("CASSANDRA UPDATED")
		return events.UnknownCommand
	case events.UnknownReply:
		return events.UnknownCommand
	default:
		return events.UnknownCommand
	}
}
