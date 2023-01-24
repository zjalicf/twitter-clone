package domain

import (
	"context"
	events "github.com/zjalicf/twitter-clone-common/common/saga/create_event"
)

type ReportStore interface {
	CreateReport(context.Context, *events.Event) (*events.Event, error)
}
