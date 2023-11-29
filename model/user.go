package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	UserName string             `json:"username,omitempty" bson:"user_name,omitempty"`
	Passowrd string             `json:"password,omitempty" bson:"password,omitempty"`
}
