package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Services struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EmployeeID string             `json:"employee_id,omitempty" bson:"employee_id,omitempty"`
	EmployerID string             `json:"employer_id,omitempty" bson:"employer_id,omitempty"`
	Value      float64            `json:"value,omitempty" bson:"value,omitempty"`
	Status     string             `json:"status,omitempty" bson:"status,omitempty"`
	Start      time.Time          `json:"start,omitempty" bson:"start,omitempty"`
	End        time.Time          `json:"end,omitempty" bson:"end,omitempty"`
}
