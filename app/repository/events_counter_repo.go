package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/xfoxawy/catchall-domain/app/model"
)

const EventsCounterCollection = "domain_events_counters"

type EventsCounterRepository struct {
	collection *mongo.Collection
}

func NewEventsCounterRepository(db *mongo.Database) *EventsCounterRepository {
	return &EventsCounterRepository{
		collection: db.Collection(EventsCounterCollection),
	}
}

func (r *EventsCounterRepository) FindOneByDomain(ctx context.Context, domain string) (model.DomainEventCounter, error) {
	var decounter model.DomainEventCounter
	err := r.collection.FindOne(ctx, bson.M{
		"domain": domain,
	}).Decode(&decounter)
	return decounter, err
}

func (r *EventsCounterRepository) IncrementCounter(ctx context.Context, domain string, evType string) (*mongo.UpdateResult, error) {
	result, err := r.collection.UpdateOne(ctx, bson.M{
		"domain": domain,
	}, bson.D{
		{
			"$inc", bson.D{
				{evType, 1},
			},
		},
	}, options.Update().SetUpsert(true))
	return result, err
}
