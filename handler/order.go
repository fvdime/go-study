package handler

import (
	"fmt"
	"net/http"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/fvdime/go-study/model"
	"github.com/fvdime/go-study/repository/order"
)

type Order struct {
	Repository *order.RedisRepository
}

// not working lol
func (h *Order) Create(w http.ResponseWriter, r *http.Request){
	var body struct{
		UserId uuid.UUID `json:"user_id"`
		Items []model.Item `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	// not good 4 production
	order := model.Order{
		OrderId: rand.Uint64(),
		UserId: body.UserId,
		Items: body.Items,
		CreatedAt: &now,
	}

	// inserting into the repository
	err := h.Repository.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to encode: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (h *Order) List(w http.ResponseWriter, r *http.Request){
	fmt.Println("list orders")
}

func (h *Order) GetById(w http.ResponseWriter, r *http.Request){
	fmt.Println("get an order by id")
}

func (h *Order) UpdateById(w http.ResponseWriter, r *http.Request){
	fmt.Println("update an order by id")
}

func (h *Order) DeleteById(w http.ResponseWriter, r *http.Request){
	fmt.Println("delete an order by id")
}