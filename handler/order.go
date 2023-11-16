package handler

import (
	"fmt"
	"net/http"
)

type Order struct {
	router http.Handler
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request){
	fmt.Println("create an order")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request){
	fmt.Println("list orders")
}

func (o *Order) GetById(w http.ResponseWriter, r *http.Request){
	fmt.Println("get an order by id")
}

func (o *Order) UpdateById(w http.ResponseWriter, r *http.Request){
	fmt.Println("update an order by id")
}

func (o *Order) DeleteById(w http.ResponseWriter, r *http.Request){
	fmt.Println("delete an order by id")
}