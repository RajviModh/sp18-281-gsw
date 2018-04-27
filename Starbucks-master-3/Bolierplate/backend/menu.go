package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"encoding/json"
)

type Item struct {
	Name string `json:"name" bson:"Name"`
	Calories    int    `json:"calories" bson:"Calories"`
	Price       int    `json:"price" bson:"Price"`
	Quantity    int    `json:"quantity" bson:"Quantity"`
	Fat int `json:"fat" bson:"Fat"`
	Cholestrol int `json:"cholestrol" bson:"Cholestrol"`
	Sodium int `json:"sodium" bson:"Sodium"`
	Protein int `json:"protein" bson:"Protein"`
	Caffeine int `json:"caffeine" bson:"Caffeine"`
	Description string `json:"description" bson:"Description"`
}

type Item1 struct {
	Name string `json:"name" bson:"Name"`
	Calories int `json:"calories" bson:"Calories"`
	Price    int `json:"price" bson:"Price"`
	Quantity int `json:"quantity" bson:"Quantity"`
}

type Description struct {
	Name string `json:"name" bson:"Name"`
	Desc string `json:"desc" bson:"Desc"`
}

// GetOrders retrieves all the orders

func (oc OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {

	var options []Item
	iter := oc.session.DB("test").C("Menu").Find(nil).Iter()
	result := Item{}
	for iter.Next(&result) {
		options = append(options, result)
	}
	fmt.Println(options)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&options)
}
