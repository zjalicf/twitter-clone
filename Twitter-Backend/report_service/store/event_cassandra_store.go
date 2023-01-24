package store

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"report_service/domain"
)

const (
	DATABASE_CASSANDRA = "events"
	COLLECTION_EVENT   = "events"
)

type EventCassandraStore struct {
	session *gocql.Session
	logger  *log.Logger
	tracer  trace.Tracer
}

func New(logger *log.Logger, tracer trace.Tracer) (*EventCassandraStore, error) {
	db := os.Getenv("EVENT_DB")
	log.Println(db)

	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		log.Println("puklo na 31")
		return nil, err
	}

	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, DATABASE_CASSANDRA, 1)).Exec()
	if err != nil {
		log.Println("puklo na 42")
		logger.Println(err)
	}
	session.Close()

	cluster.Keyspace = DATABASE_CASSANDRA
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()

	if err != nil {
		log.Println("puklo na 52")
		logger.Println(err)
		return nil, err
	}

	return &EventCassandraStore{
		session: session,
		logger:  logger,
		tracer:  tracer,
	}, nil
}

func (store *EventCassandraStore) CloseSession() {
	store.session.Close()
}

func (store *EventCassandraStore) CreateTables() {
	err := store.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id UUID, event_type text, timestamp time, PRIMARY KEY ((id), timestamp))`, COLLECTION_EVENT)).Exec()

	if err != nil {
		store.logger.Printf("CASSANDRA CREATE TABLE ERR: %s", err.Error())
	}
}

func (store *EventCassandraStore) CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	ctx, span := store.tracer.Start(ctx, "EventStore.CreateEvent")
	defer span.End()

	log.Println("Uslo u create event u cassandra store linije 83")

	tweetID, err := gocql.ParseUUID(event.TweetID)
	insert := fmt.Sprintf("INSERT INTO %s (id, event_type, timestamp) VALUES (?, ?, ?)", COLLECTION_EVENT)

	log.Println("Proslo upis u bazu")

	err = store.session.Query(
		insert, tweetID, event.Type, int64(event.Timestamp)).Exec()
	if err != nil {
		store.logger.Println(err)
		return nil, err
	}
	return event, nil
}
