package store

import (
	"fmt"
	"github.com/gocql/gocql"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
)

const (
	DATABASE_CASSANDRA = "events"
)

type EventRepo struct {
	session *gocql.Session
	logger  *log.Logger
	tracer  trace.Tracer
}

func New(logger *log.Logger, tracer trace.Tracer) (*EventRepo, error) {
	db := os.Getenv("EVENT_DB")

	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, DATABASE_CASSANDRA, 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	cluster.Keyspace = DATABASE_CASSANDRA
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()

	if err != nil {
		logger.Println(err)
		return nil, err
	}

	return &EventRepo{
		session: session,
		logger:  logger,
		tracer:  tracer,
	}, nil
}

func (sr *EventRepo) CloseSession() {
	sr.session.Close()
}

func (sr *EventRepo) CreateTables() {
	err := sr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
					(id UUID, event_type text, timestamp time, PRIMARY KEY ((id), timestamp))`, //for now there is no clustering order!!
			DATABASE_CASSANDRA)).Exec()

	if err != nil {
		sr.logger.Printf("CASSANDRA CREATE TABLE ERR: %s", err.Error())
	}
}
