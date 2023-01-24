package application

import (
	"context"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
	"go.opentelemetry.io/otel/trace"
	"log"
	"report_service/domain"
)

type ReportService struct {
	eventStore  domain.EventStore
	reportStore domain.ReportStore
	tracer      trace.Tracer
}

func NewReportService(eventStore domain.EventStore, reportStore domain.ReportStore, tracer trace.Tracer) *ReportService {
	return &ReportService{
		eventStore:  eventStore,
		reportStore: reportStore,
		tracer:      tracer,
	}
}

func (service *ReportService) CreateEvent(event events.Event) {

	eventOut := domain.Event{
		TweetID:   event.TweetID,
		Type:      event.Type,
		Timestamp: event.Timestamp,
	}
	_, err := service.eventStore.CreateEvent(context.TODO(), &eventOut)
	if err != nil {
		log.Printf("Error in report_service CreateEvent()", err.Error())
		return
	}
}

func (service *ReportService) CreateReport(event *events.Event) {
	ctx, span := service.tracer.Start(context.TODO(), "ReportService.CreateReport")
	defer span.End()

	_, err := service.reportStore.CreateReport(ctx, event)
	if err != nil {
		log.Printf("Error in report_service CreateReport()", err.Error())
		return
	} else {
		log.Println("Succesfull updated report")
	}

}
