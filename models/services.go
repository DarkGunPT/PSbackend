package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceType struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Price float64            `json:"priceHour,omitempty" bson:"priceHour,omitempty"`
}

type Appointment struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ServiceName string             `json:"service_name" bson:"service_name,omitempty"`
	Provider    User               `json:"provider,omitempty" bson:"provider,omitempty"`
	Client      User               `json:"client,omitempty" bson:"client,omitempty"`
	Status      string             `json:"status,omitempty" bson:"status,omitempty"`
	Start       time.Time          `json:"start,omitempty" bson:"start,omitempty"`
	End         time.Time          `json:"end,omitempty" bson:"end,omitempty"`
	Phone       int                `json:"phone" bson:"phone"`
	NIF         int                `json:"nif" bson:"nif"`
	Locality    string             `json:"locality" bson:"locality"`
	Notes       string             `json:"notes" bson:"notes"`
	PriceHour   float64            `json:"priceHour,omitempty" bson:"priceHour,omitempty"`
	TotalPrice  float64            `json:"totalPrice,omitempty" bson:"totalPrice,omitempty"`
}

type Review struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EmployeeID string             `json:"employee_id,omitempty" bson:"employee_id,omitempty"`
	EmployerID string             `json:"employer_id,omitempty" bson:"employer_id,omitempty"`
	Date       time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	Comment    string             `json:"comment,omitempty" bson:"comment,omitempty"`
}
