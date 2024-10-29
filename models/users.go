package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	Password      string             `json:"password,omitempty" bson:"password,omitempty"`
	NIF           int64              `json:"nif,omitempty" bson:"nif,omitempty"`
	Phone         int64              `json:"phone,omitempty" bson:"phone,omitempty"`
	Email         string             `json:"email,omitempty" bson:"email,omitempty"`
	Role          []Role             `json:"role,omitempty" bson:"role,omitempty"`
	ServiceTypes  []ServiceType      `json:"service_types,omitempty" bson:"service_types,omitempty"`
	Locality      string             `json:"locality,omitempty" bson:"locality,omitempty"`
	Rating        float64            `json:"rating,omitempty" bson:"rating,omitempty"`
	BlockServices bool               `json:"block_services,omitempty" bson:"block_services,omitempty"`
	IsActive      bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	RecoveryCode  int                `json:"recovery_code,omitempty" bson:"recovery_code,omitempty"`
}

type Role struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

type Fee struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	NIF      int64              `json:"nif,omitempty" bson:"nif,omitempty"`
	Value    float64            `json:"value,omitempty" bson:"value,omitempty"`
	JobsDone int64              `json:"jobs_done,omitempty" bson:"jobs_done,omitempty"`
	Paid     bool               `json:"paid,omitempty" bson:"paid,omitempty"`
	Date     time.Time          `json:"date,omitempty" bson:"date,omitempty"`
}
