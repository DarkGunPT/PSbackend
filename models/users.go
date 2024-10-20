package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	NIF          int64              `json:"nif,omitempty" bson:"nif,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone        int64              `json:"phone,omitempty" bson:"phone,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Password     string             `json:"password,omitempty" bson:"password,omitempty"`
	Role         string             `json:"role,omitempty" bson:"role,omitempty"`
	Locality     string             `json:"locality,omitempty" bson:"locality,omitempty"`
	Profession   string             `json:"profission,omitempty" bson:"profission,omitempty"`
	PricePerHour float64            `json:"price_hour,omitempty" bson:"price_hour,omitempty"`
	Rating       float64            `json:"rating,omitempty" bson:"rating,omitempty"`
	Fee          float64            `json:"fee,omitempty" bson:"fee,omitempty"`
	JobsDone     int64              `json:"jobs_done,omitempty" bson:"jobs_done,omitempty"`
	Block        bool               `json:"block,omitempty" bson:"block,omitempty"`
}
