package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"shop/apps/orders/internal/application/usecases"
	"shop/apps/orders/internal/domain/entities"
)

type OrderController struct {
	orderService *application.OrdersUsecase
}

func NewOrderController(orderService *application.OrdersUsecase) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order entities.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newOrder, err := c.orderService.CreateOrder(r.Context(), order.Name, order.Price, order.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}

func (c *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := c.orderService.GetOrder(r.Context(), id)
	if err != nil {
		if err == application.ErrOrderNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (c *OrderController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order entities.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.ID = id
	if err := c.orderService.UpdateOrder(r.Context(), &order); err != nil {
		if err == application.ErrOrderNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (c *OrderController) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if err := c.orderService.DeleteOrder(r.Context(), id); err != nil {
		if err == application.ErrOrderNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
