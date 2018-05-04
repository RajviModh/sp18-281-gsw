package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"encoding/json"
)



//-----------------------------------------------------------Code Goes Here------------------------------------------------------------------


// GetCartItems retrieves all the cart orders
func (oc OrderController) GetCartItems(w http.ResponseWriter, r *http.Request) {

	var cart []Cart
	iter := oc.session.DB("test").C("Cart").Find(nil).Iter()
	result := Cart{}
	for iter.Next(&result) {
		cart = append(cart, result)
	}

	for _, cart := range cart {
		fmt.Println("--****************************- ", cart.Username, cart.Items)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&cart)

}

//-----------------------------------------------------------Function Goes Here----------------------------------------------------------------

func changeStatusToPlaced(orderId string, oc OrderController) {
	fmt.Println("placed")
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "PLACED"}})

}

func changeStatusToPaid(orderId string, oc OrderController) {
	fmt.Println("paid")
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "PAID"}})

}
