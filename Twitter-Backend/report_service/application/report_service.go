package application

import (
	"context"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
	"go.opentelemetry.io/otel/trace"
	"log"
	"report_service/domain"
	"time"
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

	thisT := time.Unix(int64(event.Timestamp), 0)
	dailyUnix := time.Date(thisT.Year(), thisT.Month(), thisT.Day(), 0, 0, 0, 0, time.UTC).Unix() //dailyRepDate
	monthlyUnix := time.Date(thisT.Year(), thisT.Month(), 1, 0, 0, 0, 0, time.UTC).Unix()         //monthlyRepDate
	_, err := service.reportStore.CreateReport(ctx, event, monthlyUnix, dailyUnix)
	if err != nil {
		log.Printf("Error in report_service CreateReport()", err.Error())
		return
	} else {
		log.Println("Succesfull updated report")
	}

}

func (service *ReportService) GetReportForAd(ctx context.Context, tweetID string, reportType string) (*domain.Report, error) {

	result, err := service.reportStore.GetReportForAd(ctx, tweetID, reportType)
	if err != nil {
		log.Printf("Error in ReportService GetReportForAd: %s", err.Error())
		return nil, err
	}
	return result, nil
}

//func (service *ReportService) (){}
