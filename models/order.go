package models

import "go.mongodb.org/mongo-driver/v2/bson"

type OrderStatus int

const (
	Processing OrderStatus = iota
	InProgress
	Done
)

type Order struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductID bson.ObjectID `bson:"product_id" json:"product_id"`
	Price     float64       `bson:"price" json:"price"`
	Status    OrderStatus   `bson:"status" json:"status"`
}
