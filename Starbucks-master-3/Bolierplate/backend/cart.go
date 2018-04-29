package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)



//-----------------------------------------------------------Code Goes Here------------------------------------------------------------------


// GetCartItems retrieves all the cart orders

func (oc OrderController) GetCartItems(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r)
	vars := mux.Vars(r)
	id := vars["username"]

	var cart Cart
	if error := oc.session.DB("test").C("Cart").Find(bson.M{"Username": id}).One(&cart); error != nil {
		fmt.Println(error)
		w.WriteHeader(400)
		data := `{"status":"error","message":"Cart doesn't exist for this user!"}`
		json.NewEncoder(w).Encode(data)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&cart)
	}

}
// Delete Items deletes the order with specified order id
func (oc OrderController) DeleteItems(w http.ResponseWriter, r *http.Request) {

	var item1 Item
	fmt.Println(r)
	id := r.FormValue("id")
	data := r.FormValue("data")

	json.Unmarshal([]byte(data), &item1)
	fmt.Printf("%+v\n", item1)

	orderId := id
	itemname := item1.Name
	qty := item1.Quantity
	price := item1.Price

	fmt.Println("--------%%%%%%%%%%%%%%%%%%------------", price)

	fmt.Println("--------%%%%%%%%%%%%%%%%%%------------", itemname)

	//q1,_ := strconv.Atoi(qty)
	qty -= 1

	fmt.Println("--------%%%%%%%%%%%%%%%%%%------------", qty)
	// Fetch order
	if err := oc.session.DB("test").C("Cart").Update(bson.M{"_id": orderId, "Items.Name": itemname}, bson.M{"$set": bson.M{"Items.$.Quantity": qty}}); err != nil {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		data := `{"status":"success","message":"Order quantity has been decremented"}`
		json.NewEncoder(w).Encode(data)
	}
}
// Delete Cart deletes the entire item
func (oc OrderController) DeleteCart(w http.ResponseWriter, r *http.Request) {

	var item1 Item
	fmt.Println(r)
	orderid := r.FormValue("id")
	data := r.FormValue("data")

	json.Unmarshal([]byte(data), &item1)
	fmt.Printf("%+v\n", item1)

	itemname := item1.Name

	fmt.Println("--------%%%%%%%%In DELETE%%%%%%%%%%------------", orderid, data, itemname)

	// Delete order
	if err := oc.session.DB("test").C("Cart").Update(bson.M{"_id": orderid}, bson.M{"$pull": bson.M{"Items": bson.M{"Name": itemname}}}); err != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		data := `{"status":"success","message":"Order Items has been deleted"}`
		json.NewEncoder(w).Encode(data)
	}

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
