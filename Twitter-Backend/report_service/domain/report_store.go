package domain

import (
	"context"
)

type ReportStore interface {
	CreateEvent(context.Context, Event) (*Event, error)
}
