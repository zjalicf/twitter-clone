package handlers

import (
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
	"tweet_service/application"
)

type CreateReportCommandHandler struct {
	reportService     *application.TweetService
	publisher         saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateEventCommandHandler(reportService *application.TweetService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateReportCommandHandler, error) {
	o := &CreateReportCommandHandler{
		reportService:     reportService,
		publisher:         publisher,
		commandSubscriber: subscriber,
	}
	//prijava za slusanje komandi
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// hendlovanje komandama
func (handler *CreateReportCommandHandler) handle(command *events.CreateEventCommand) {
	reply := events.CreateEventReply{Event: command.Event}
	switch command.Type {

	//case events.SendMessageToReportService:
	//	log.Println("Salje se poruka u tweet_service")
	//	reply.Type = events.MessageRecieved

	//case events.UpdateCassandra:
	//	log.Println("Uslo u mongo update")
	//	reply.Type = events.CassandraUpadated

	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.publisher.Publish(reply)
	}
}
