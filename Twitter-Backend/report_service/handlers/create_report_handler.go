package handlers

import (
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
	saga "github.com/zjalicf/twitter-clone-common/common/saga/messaging"
	"report_service/application"
)

type CreateReportCommandHandler struct {
	reportService     *application.ReportService
	publisher         saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateEventCommandHandler(reportService *application.ReportService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateReportCommandHandler, error) {
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

//hendlovanje komandama
func (handler *CreateReportCommandHandler) handle(command *events.CreateEventCommand) {
	reply := events.CreateEventReply{Event: command.Event}
	switch command.Type {

	case events.UpdateMongo:
		reply.Type = events.MongoUpdated

	case events.UpdateCassandra:
		handler.reportService.CreateEvent(command.Event)
		reply.Type = events.CassandraUpadated

	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.publisher.Publish(reply)
	}
}
