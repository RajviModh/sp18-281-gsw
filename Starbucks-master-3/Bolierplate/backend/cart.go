package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type Cart struct {
	CartId   string  `json:"id" bson:"_id"`
	Items    []Item1 `json:"items" bson:"Items"`
	Username string  `json:"username" bson:"Username"`
}


//-----------------------------------------------------------Code Goes Here------------------------------------------------------------------



//-----------------------------------------------------------Function Goes Here----------------------------------------------------------------

func changeStatusToPlaced(orderId string, oc OrderController) {
	fmt.Println("placed")
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "PLACED"}})

}

func changeStatusToPaid(orderId string, oc OrderController) {
	fmt.Println("paid")
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "PAID"}})

}
