package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/OrdersAPI/models"
	"github.com/SomeSuperCoder/OrdersAPI/repository"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var validate = validator.New()

type OrderHandler struct {
	Repo *repository.OrderRepo
}

// ===== Update order status =====
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	parsedId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid id provided", http.StatusBadRequest)
		return
	}

	order, err := h.Repo.GetOrder(r.Context(), parsedId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Falied to get: %v", err.Error()), http.StatusNotFound)
		return
	}

	marshaledOrder, err := json.Marshal(order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(marshaledOrder))
}

// ===== Create order =====
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ProductID bson.ObjectID `json:"product_id"`
		Price     float64       `json:"price" validate:"gt=0"`
	}

	// Parse
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Falied to parse JSON: %v", err.Error()), http.StatusBadRequest)
		return
	}

	// Validate
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON validation error: %v", err.Error()), http.StatusNotAcceptable)
		return
	}

	// Do work
	order := models.Order{
		ProductID: request.ProductID,
		Price:     request.Price,
		Status:    models.Processing,
	}

	createdOrder, err := h.Repo.CreateOrder(r.Context(), order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// Respond
	createdOrderMarshaled, err := json.Marshal(createdOrder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(createdOrderMarshaled))
}

// ===== Update order status =====
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var request struct {
		Status models.OrderStatus `json:"status" bson:"status,omitempty" validate:"omitempty,oneof=1 2 3"`
	}

	// Parse
	parsedId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid id provided: %v", err), http.StatusBadRequest)
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err.Error()), http.StatusBadRequest)
		return
	}

	// Validate
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON validation error: %v", err.Error()), http.StatusNotAcceptable)
		return
	}

	// Do work
	err = h.Repo.UpdateOrder(r.Context(), parsedId, request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// Respond
	fmt.Fprintln(w, "Order updated succesfully")
}
