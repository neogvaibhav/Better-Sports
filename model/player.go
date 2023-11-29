package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	PlayerProfile *PlayerProfile     `json:"player_profile,omitempty" bson:"player_profile,omitempty"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type PlayerProfile struct {
	Grade    string `json:"grade,omitempty" bson:"grade,omitempty"`
	Position string `json:"position,omitempty" bson:"position,omitempty"`
	Stats    *Stats `json:"stats,omitempty" bson:"stats,omitempty"`
}

type Stats struct {
	Goals   int `json:"goals" bson:"goals"`
	Assists int `json:"assists" bson:"assists"`
	Fouls   int `json:"fouls" bson:"fouls"`
}
