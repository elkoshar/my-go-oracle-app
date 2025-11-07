package database

import (
	"context"
	"time"
)

type Event struct {
	Name string
	Time time.Time
}

type ContextKey string

const MongoEventContextKey ContextKey = "MysqlEvent"

type MysqlMonitor struct {
}

func StartMetrics(ctx context.Context, event Event) context.Context {
	event.Time = time.Now()
	return context.WithValue(ctx, MongoEventContextKey, event)
}

func getEventFromContext(ctx context.Context) *Event {
	val := ctx.Value(MongoEventContextKey)
	if val == nil {
		return nil
	}

	event := val.(Event)

	return &event
}

func NewCommandMonitor() MysqlMonitor {
	return MysqlMonitor{}
}

func (m *MysqlMonitor) Succeeded(ctx context.Context) {
	eventMysqlFromClient := getEventFromContext(ctx)
	if eventMysqlFromClient == nil {
		return
	}

}

func (m *MysqlMonitor) Failed(ctx context.Context) {

	eventMysqlFromClient := getEventFromContext(ctx)
	if eventMysqlFromClient == nil {
		return
	}

}
