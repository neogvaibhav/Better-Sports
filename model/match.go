package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID    primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Date  time.Time            `json:"match_date,omitempty" bson:"match_date,omitempty"`
	Teams []primitive.ObjectID `json:"teams,omitempty" bson:"teams,omitempty"`
}
