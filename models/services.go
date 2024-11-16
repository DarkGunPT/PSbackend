package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Services struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EmployeeID  string             `json:"employee_id,omitempty" bson:"employee_id,omitempty"`
	ServiceType ServiceType        `json:"service_type,omitempty" bson:"service_type,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Value       float64            `json:"value,omitempty" bson:"value,omitempty"`
	Appointment []Appointment      `json:"appointments,omitempty" bson:"appointments,omitempty"`
}

type ServiceType struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

type Appointment struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EmployeeID string             `json:"employee_id,omitempty" bson:"employee_id,omitempty"`
	EmployerID string             `json:"employer_id,omitempty" bson:"employer_id,omitempty"`
	Status     string             `json:"status,omitempty" bson:"status,omitempty"`
	Start      time.Time          `json:"start,omitempty" bson:"start,omitempty"`
	End        time.Time          `json:"end,omitempty" bson:"end,omitempty"`
}

type Review struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EmployeeID string             `json:"employee_id,omitempty" bson:"employee_id,omitempty"`
	EmployerID string             `json:"employer_id,omitempty" bson:"employer_id,omitempty"`
	Date       time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	Comment    string             `json:"comment,omitempty" bson:"comment,omitempty"`
}
