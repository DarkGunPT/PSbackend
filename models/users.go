package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	Password      string             `json:"password" bson:"password"`
	NIF           int                `json:"nif" bson:"nif"`
	Phone         int                `json:"phone" bson:"phone"`
	Email         string             `json:"email" bson:"email"`
	Role          []Role             `json:"role" bson:"role"`
	ServiceTypes  []ServiceType      `json:"service_types" bson:"service_types"`
	Locality      string             `json:"locality" bson:"locality"`
	Rating        float64            `json:"rating" bson:"rating"`
	BlockServices bool               `json:"block_services" bson:"block_services"`
	IsActive      bool               `json:"is_active" bson:"is_active"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	RecoveryCode  int                `json:"recovery_code" bson:"recovery_code"`
	WorkStart     time.Time          `json:"workStart" bson:"workStart"`
	WorkEnd       time.Time          `json:"workEnd" bson:"workEnd"`
}

type Role struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	ServicesDone int                `json:"services_done" bson:"services_done"`
}

type Fee struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	NIF      int64              `json:"nif,omitempty" bson:"nif,omitempty"`
	Value    float64            `json:"value,omitempty" bson:"value,omitempty"`
	JobsDone int64              `json:"jobs_done,omitempty" bson:"jobs_done,omitempty"`
	Paid     bool               `json:"paid,omitempty" bson:"paid,omitempty"`
	Month    string             `json:"month" bson:"month,omitempty"`
	Year     string             `json:"year" bson:"year,omitempty"`
	Date     time.Time          `json:"date,omitempty" bson:"date,omitempty"`
}
