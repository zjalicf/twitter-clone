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
	DATABASE           = "report_mongo"
	COLLECTION_DAILY   = "daily_reports"
	COLLECTION_MONTHLY = "monthly_reports"
)

type ReportMongoDBStore struct {
	dailyReports   *mongo.Collection
	monthlyReports *mongo.Collection
	tracer         trace.Tracer
}

func (store *ReportMongoDBStore) GetReportForAd(ctx context.Context, tweetID string, reportType string) (*domain.Report, error) {

	if reportType == "daily" {

		result, err := store.filterOneDaily(bson.M{"tweet_id": tweetID})
		if err != nil {
			log.Printf("Error in ReportMongoDB filterOneDaily: &s", err.Error())
			return nil, err
		}
		return result, nil

	} else if reportType == "monthly" {

		result, err := store.filterOneMonthly(bson.M{"tweet_id": tweetID})
		if err != nil {
			log.Printf("Error in ReportMongoDB filterOneMonthly: &s", err.Error())
			return nil, err
		}
		return result, nil

	}
	return nil, nil
}

func (store *ReportMongoDBStore) CreateReport(ctx context.Context, event *events.Event,
	monthlyUnix, dailyUnix int64) (*events.Event, error) {
	ctx, span := store.tracer.Start(ctx, "ReportMongoDBStore.CreateReport")
	defer span.End()

	oneDaily, err := store.filterOneDaily(bson.M{"tweet_id": event.TweetID, "timestamp": dailyUnix})
	if err != nil {
		log.Println(err.Error())
		report := domain.Report{
			ID:          primitive.NewObjectID(),
			TweetID:     event.TweetID,
			Timestamp:   dailyUnix,
			LikeCount:   0,
			UnlikeCount: 0,
			ViewCount:   0,
			TimeSpent:   0,
		}

		if event.Type == "Liked" {
			report.LikeCount++
		} else if event.Type == "Unliked" {
			report.UnlikeCount++
		} else {
			report.ViewCount++
		}

		_, err = store.dailyReports.InsertOne(ctx, report)
		if err != nil {
			log.Printf("Error in ReportMongoStore, daily_reporst: %s", err.Error())
			return nil, err
		}
	} else {

		//like update

		if event.Type == "Liked" {
			oneDaily.LikeCount = oneDaily.LikeCount + 1
			_, err = store.dailyReports.UpdateOne(context.TODO(), bson.M{"tweet_id": event.TweetID}, bson.M{"$set": oneDaily})
			if err != nil {
				log.Printf("Error in report_mongodb CreateReport() Like daily: %s", err.Error())
				return nil, err
			}
		} else if event.Type == "Unliked" {

			//unline update

			oneDaily.UnlikeCount = oneDaily.UnlikeCount + 1
			_, err = store.dailyReports.UpdateOne(context.TODO(), bson.M{"tweet_id": event.TweetID}, bson.M{"$set": oneDaily})
			if err != nil {
				log.Printf("Error in report_mongodb CreateReport() Unlike daily: %s", err.Error())
				return nil, err
			}
		} else {
			//view update
		}

	}
	//monthly
	oneMonthly, err := store.filterOneMonthly(bson.M{"tweet_id": event.TweetID, "timestamp": monthlyUnix})
	if err != nil {
		log.Println(err.Error())
		report := domain.Report{
			ID:          primitive.NewObjectID(),
			TweetID:     event.TweetID,
			Timestamp:   monthlyUnix,
			LikeCount:   0,
			UnlikeCount: 0,
			ViewCount:   0,
			TimeSpent:   0,
		}

		if event.Type == "Liked" {
			report.LikeCount++
		} else if event.Type == "Unliked" {
			report.UnlikeCount++
		} else {
			report.ViewCount++
		}
		_, err = store.monthlyReports.InsertOne(ctx, report)
		if err != nil {
			log.Printf("Error in ReportMongoStore, monthly_reporst: %s", err.Error())
			return nil, err
		}
	} else {

		//like update

		if event.Type == "Liked" {
			oneMonthly.LikeCount = oneMonthly.LikeCount + 1
			_, err = store.monthlyReports.UpdateOne(context.TODO(), bson.M{"tweet_id": event.TweetID}, bson.M{"$set": oneMonthly})
			if err != nil {
				log.Printf("Error in report_mongodb CreateReport() Like monthly: %s", err.Error())
				return nil, err
			}
		} else if event.Type == "Unliked" {

			//unline update

			oneMonthly.UnlikeCount = oneMonthly.UnlikeCount + 1
			_, err = store.monthlyReports.UpdateOne(context.TODO(), bson.M{"tweet_id": event.TweetID}, bson.M{"$set": oneMonthly})
			if err != nil {
				log.Printf("Error in report_mongodb CreateReport() Unlike monthly: %s", err.Error())
				return nil, err
			}
		} else {
			//view update
		}

	}

	return event, nil
}

func NewReportMongoDBStore(client *mongo.Client, tracer trace.Tracer) domain.ReportStore {
	dailyReports := client.Database(DATABASE).Collection(COLLECTION_DAILY)
	monthlyReports := client.Database(DATABASE).Collection(COLLECTION_MONTHLY)

	return &ReportMongoDBStore{
		dailyReports:   dailyReports,
		monthlyReports: monthlyReports,
		tracer:         tracer,
	}
}

func (store *ReportMongoDBStore) filterDaily(filter interface{}) ([]*domain.Report, error) {
	cursor, err := store.dailyReports.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		log.Printf("Error in report_mongodb filter() Unlike daily: %s", err.Error())
		return nil, err
	}
	return decode(cursor)
}

func (store *ReportMongoDBStore) filterOneDaily(filter interface{}) (user *domain.Report, err error) {
	result := store.dailyReports.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func (store *ReportMongoDBStore) UpdateOneDaily(filter interface{}) (user *domain.Report, err error) {
	result := store.dailyReports.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func (store *ReportMongoDBStore) filterMonthly(filter interface{}) ([]*domain.Report, error) {
	cursor, err := store.monthlyReports.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		log.Printf("Error in report_mongodb filter() Unlike monthly: %s", err.Error())
		return nil, err
	}
	return decode(cursor)
}

func (store *ReportMongoDBStore) filterOneMonthly(filter interface{}) (user *domain.Report, err error) {
	result := store.monthlyReports.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func (store *ReportMongoDBStore) UpdateOneMonthly(filter interface{}) (user *domain.Report, err error) {
	result := store.monthlyReports.FindOne(context.TODO(), filter)
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
