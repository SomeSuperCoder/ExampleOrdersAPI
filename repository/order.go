package repository

import (
	"context"

	"github.com/SomeSuperCoder/OrdersAPI/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepo struct {
	database *mongo.Database
	orders   *mongo.Collection
}

func NewOrderRepo(db *mongo.Database) *OrderRepo {
	return &OrderRepo{
		database: db,
		orders:   db.Collection("orders"),
	}
}

// Implementation
func (r *OrderRepo) CreateOrder(ctx context.Context, order models.Order) (*models.Order, error) {
	res, err := r.orders.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	order.ID = res.InsertedID.(bson.ObjectID)

	return &order, nil
}

func (r *OrderRepo) UpdateOrder(ctx context.Context, id bson.ObjectID, update any) error {
	_, err := r.orders.UpdateByID(ctx, id, bson.M{
		"$set": update,
	})
	return err
}

func (r *OrderRepo) GetOrder(ctx context.Context, id bson.ObjectID) (*models.Order, error) {
	var order models.Order

	err := r.orders.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&order)

	return &order, err
}
