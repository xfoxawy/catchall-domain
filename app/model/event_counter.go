package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DomainEventCounter struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Domain    string             `bson:"domain"`
	Delivered int                `bson:"delivered"`
	Bounced   int                `bson:"bounced"`
}
