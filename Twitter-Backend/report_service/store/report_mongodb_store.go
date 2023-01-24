package store

import (
	"context"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
	"log"
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

func (store *ReportMongoDBStore) CreateReport(ctx context.Context, event *events.Event) (*events.Event, error) {
	ctx, span := store.tracer.Start(ctx, "ReportMongoDBStore.CreateReport")
	defer span.End()

	one, err := store.filterOne(bson.M{"tweet_id": event.TweetID})
	if err != nil {
		log.Println(err.Error())
		report := domain.Report{
			ID:          primitive.NewObjectID(),
			TweetID:     event.TweetID,
			LikeCount:   0,
			UnlikeCount: 0,
			ViewCount:   0,
		}

		if event.Type == "Liked" {
			report.LikeCount++
		} else if event.Type == "Unliked" {
			report.UnlikeCount++
		} else {
			report.ViewCount++
		}

		_, err = store.reports.InsertOne(ctx, report)
		if err != nil {
			log.Printf("Error in ReportMongoStore, line 43: %s", err.Error())
			return nil, err
		}
	} else {

		//like update

		if event.Type == "Liked" {
			one.LikeCount = one.LikeCount + 1
			_, err = store.reports.UpdateOne(context.TODO(), bson.M{"tweet_id": event.TweetID}, bson.M{"$set": one})
			if err != nil {
				log.Printf("Error in report_mongodb CreateReport() Like: %s", err.Error())
				return nil, err
			}
		} else if event.Type == "Unliked" {

			//unline update

			one.UnlikeCount = one.UnlikeCount + 1
			_, err = store.reports.UpdateOne(context.TODO(), bson.M{"tweet_id": event.TweetID}, bson.M{"$set": one})
			if err != nil {
				log.Printf("Error in report_mongodb CreateReport() Unlike: %s", err.Error())
				return nil, err
			}
		} else {
			//view update
		}
	}

	return event, nil
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
		log.Printf("Error in report_mongodb filter() Unlike: %s", err.Error())
		return nil, err
	}
	return decode(cursor)
}

func (store *ReportMongoDBStore) filterOne(filter interface{}) (user *domain.Report, err error) {
	result := store.reports.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func (store *ReportMongoDBStore) UpdateOne(filter interface{}) (user *domain.Report, err error) {
	result := store.reports.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func decode(cursor *mongo.Cursor) (reports []*domain.Report, err error) {
	for cursor.Next(context.TODO()) {
		var report domain.Report
		err = cursor.Decode(&report)
		if err != nil {
			log.Printf("Error in report_mongodb decode(): %s", err.Error())
			return
		}
		reports = append(reports, &report)
	}
	err = cursor.Err()
	return
}
