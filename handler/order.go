package handler

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/orders-api/model"
	"github.com/orders-api/repository/order"
)

type Order struct {
	Repo *order.RedisRepo
}


func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID `json:"customer_id"`
		LineItems []model.LineItem  `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Write([]byte("Please check and provide correct data!!!"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	order := model.Order{
		OrderID: rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems: body.LineItems,
		CreatedAt: &now,
	}

	err := o.Repo.Insert(r.Context(), order)
	if err != nil {
		w.Write([]byte("Cannot create Order!!!"))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	res, err := json.Marshal(order)
	if err != nil {
		w.Write([]byte("Json Marshal Error!!!"))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated) 
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.Write([]byte("Please provide correct page number!!!"))
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	const size = 50
	res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{Offset: cursor, Size: size})
	if err != nil {
		w.Write([]byte("Server Error. Try Again Later!!"))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	var response struct {
		Items []model.Order `json:"items"`
		Next uint64 `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = uint64(res.Cursor)

	data, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte("Json encode error!!"))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	_id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		w.Write([]byte("Please check and provide correct id!!!"))
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	res, err := o.Repo.FindByID(r.Context(), _id)
	if errors.Is(err, order.ErrNotExist) {
		w.Write([]byte("Order Not Found"))
		w.WriteHeader(http.StatusNotFound)
		return
	}else if err != nil {
		w.Write([]byte("Server Error. Try Again Later!!!"))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.Write([]byte("Server Error. Try Again Later!!!"))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.WriteHeader(http.StatusOK)
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Write([]byte("Please check your data and try again!!!"))
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	idParam := chi.URLParam(r, "id")
	_id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		w.Write([]byte("Please check and provide id correctly!!"))
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	theOrder, err := o.Repo.FindByID(r.Context(), _id)
	if errors.Is(err, order.ErrNotExist) {
		w.Write([]byte("Order does not exist"))
		w.WriteHeader(http.StatusNotFound)
		return 
	} else if err != nil {
		w.Write([]byte("Server Error. Try Again Later!!!"))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	var status string = body.Status
	now := time.Now().UTC()
	switch status {
	case "shipped":
		if theOrder.ShippedAt == nil {
			theOrder.ShippedAt = &now 
		} else {
			w.Write([]byte("Order is ALready shipped!!"))
			w.WriteHeader(http.StatusBadRequest)
			return 
		}
	case "completed":
		if theOrder.CompletedAt == nil && theOrder.ShippedAt != nil {
			theOrder.CompletedAt = &now 
		} else {
			w.Write([]byte("Order not shipped or delivered already!!"))
			w.WriteHeader(http.StatusBadRequest)
			return 
		}
	default:
		w.Write([]byte("Invalid status value provided. Possible choices are shipped and completed!!!"))
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	err = o.Repo.Update(r.Context(), theOrder)
	if err != nil {
		w.Write([]byte("Server Error. Try Again Later!!!"))
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.Write([]byte("Updated order successfully!!"))
	w.WriteHeader(http.StatusAccepted)
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	_id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		w.Write([]byte("Please check and provide id correctly"))
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	err = o.Repo.DeleteByID(r.Context(), _id)
	if errors.Is(err, order.ErrNotExist) {
		w.Write([]byte("Order Not Found"))
		w.WriteHeader(http.StatusNotFound)
		return 
	} else if err != nil {
		w.Write([]byte("Error Occurred while Deleting..."))
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	w.Write([]byte("Deleted Order Successfully"))
	w.WriteHeader(http.StatusOK)
}