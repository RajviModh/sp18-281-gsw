package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"encoding/json"
	"net/http"
	time2 "time"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Order struct {
	OrderId     string     `json:"id" bson:"_id"`
	UserName    string     `json:"username" bson:"username"`
	Location    string     `json:"location" bson:"location"`
	Items       []Item     `json:"items" bson:"items"`
	Status      string     `json:"status" bson:"status"`
	Message     string     `json:"message" bson:"message"`
	Links       Links      `json:"links" bson:"links"`
	TotalAmount int    `json:"totalAmount" bson:"totalAmount"`
	OrderDate   time2.Time `json: "orderDate" bson: "orderDate"`
	PaymentDate time2.Time `json: "paymentDate" bson: "paymentDate"`
}

func (oc OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside createorder	")
	var o Order
	// Populate the order data from request.body to Order object
	json.NewDecoder(r.Body).Decode(&o)

	fmt.Println(o.Items[0].Name)
	fmt.Println(o.Items[1].Name)
	fmt.Println(o.Items[0])
	// Add an Id, using uuid for
	var orderId uuid.UUID
	orderId = uuid.NewV4()
	o.OrderId = orderId.String()

	o.Status = "PLACED"
	o.Message = "Order has been placed"
	o.TotalAmount = 0

	for index := 0; index < len(o.Items); index += 1 {
		o.TotalAmount += o.Items[index].Price * o.Items[index].Quantity
	}

	o.OrderDate = time2.Now()

	// Write the user to mongo
	oc.session.DB("test").C("Order").Insert(&o)

	// Write content-type, statuscode, payload
	fmt.Println("New Order Created, Order ID:", o.OrderId)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(o)
}

func (oc OrderController) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	fmt.Println(username)
	var orders []Order
	iter := oc.session.DB("test").C("Order").Find(bson.M{"username": username}).Iter()
	result := Order{}
	for iter.Next(&result) {
		orders = append(orders, result)
	}

	for _, order := range orders {
		//fmt.Println(order.OrderId)
		//fmt.Println(order.Items[0])
		fmt.Println("--- ", order.OrderId)
		fmt.Println("---", order.Location)
		//fmt.Println("------------",order.Items)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&orders)

}

func (oc OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderId := vars["id"]

	o := Order{}
	fmt.Println("------Order Id------", orderId)
	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&o); err != nil {
		w.WriteHeader(404)

		data := `{"status":"error","message":"Order not found"}`
		json.NewEncoder(w).Encode(data)
		return
	}
	fmt.Println(o)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(o)
}

func (oc OrderController) OrderPayment(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("inside createorder	")
	vars := mux.Vars(r)
	orderId := vars["id"]
	order := Order{}
	fmt.Println("pay order")
	fmt.Println(orderId)

	//call get order method
	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&order); err != nil {
		w.WriteHeader(404)
		return
	}

	json.NewDecoder(r.Body).Decode(&order)
	//order.Status :="PAID"
	fmt.Println(order)
	username := order.UserName //change to the value fetched from session
	totalAmount := order.TotalAmount

	fmt.Println(username)
	var user User
	if error := oc.session.DB("test").C("User").Find(bson.M{"username": username}).One(&user); error != nil {
		fmt.Println(error)
		w.WriteHeader(400)
		data := `{"status":"error","message":"User doesn't exist anymore'"}`
		json.NewEncoder(w).Encode(data)
		return
	}

	fmt.Println(user.FirstName)
	fmt.Println(user)
	/*	user.FirstName = "Anurag"
		user.LastName = "Panchal"*/
	credits := user.Credits
	fmt.Println(credits)
	credits -= totalAmount
	fmt.Println(credits)
	if credits < 0 {
		w.WriteHeader(400)
		data := `{"status":"error","message":"Not enough credits"}`
		json.NewEncoder(w).Encode(data)
		return
	}
	error := oc.session.DB("test").C("User").Update(bson.M{"username": username}, bson.M{"$set": bson.M{"credits": credits}})
	fmt.Println(error)
	oc.session.DB("test").C("User").Find(bson.M{"username": username}).One(&user)
	credits = user.Credits
	fmt.Println("Updated credits:", credits)
	/*if order.Status == "PAID" || order.Status == "PREPARING" || order.Status == "SERVED" || order.Status == "COLLECTED" {
		w.WriteHeader(400)
		data := `{"status":"error","message":"Order payment rejected "}`
		json.NewEncoder(w).Encode(data)
		return
	}*/

	//code to update status to paid goes here

	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "PAID", "message": "Payment Accepted", "paymentDate": time2.Now()}})
	//oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$unset": bson.M{"Links.Payment": ""}})
	fmt.Println("Order Status Updated: ", order.Status)

	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&order); err != nil {
		w.WriteHeader(404)
		data := `{"status":"error","message":"Order not found"}`
		json.NewEncoder(w).Encode(data)
		return
	}
	// to stop displaying payment after clicking on order pay(since payment set to omit empty)
	//order.Links.Payment = ""
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(order)
}

