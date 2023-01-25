package domain

import (
	"context"
)

type EventStore interface {
	CreateEvent(context.Context, *Event) (*Event, error)
}
