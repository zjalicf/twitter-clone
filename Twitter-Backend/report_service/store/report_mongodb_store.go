package store

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
	"report_service/domain"
)

const (
	DATABASE   = "report_mongo"
	COLLECTION = "reports"
)

type ReportMongoDBStore struct {
	reports *mongo.Collection
	tracer  trace.Tracer
}

func (store *ReportMongoDBStore) CreateEvent(ctx context.Context, event domain.Event) (*domain.Event, error) {
	//TODO implement me
	panic("implement me")
}

func NewReportMongoDBStore(client *mongo.Client, tracer trace.Tracer) domain.ReportStore {
	reports := client.Database(DATABASE).Collection(COLLECTION)
	return &ReportMongoDBStore{
		reports: reports,
		tracer:  tracer,
	}
}

func (store *ReportMongoDBStore) filter(filter interface{}) ([]*domain.Report, error) {
	cursor, err := store.reports.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *ReportMongoDBStore) filterOne(filter interface{}) (user *domain.Report, err error) {
	result := store.reports.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func decode(cursor *mongo.Cursor) (reports []*domain.Report, err error) {
	for cursor.Next(context.TODO()) {
		var report domain.Report
		err = cursor.Decode(&report)
		if err != nil {
			return
		}
		reports = append(reports, &report)
	}
	err = cursor.Err()
	return
}
