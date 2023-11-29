package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	ID      primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string               `json:"team_name,omitempty" bson:"team_name,omitempty"`
	Players []primitive.ObjectID `json:"players,omitempty" bson:"players,omitempty"`
}
