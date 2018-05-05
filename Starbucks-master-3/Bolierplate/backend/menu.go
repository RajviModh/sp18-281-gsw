package main

import (
	"fmt"
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

// Get Item Information

/*func (oc OrderController) GetItemInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var uuid string = params["item"]

	fmt.Println(uuid)
	result := Description{}
	if err := oc.session.DB("test").C("ItemDesc").Find(bson.M{"Name": uuid}).One(&result); err != nil {
			w.WriteHeader(404)
			data := `{"status":"error","message":"Order not found"}`
			json.NewEncoder(w).Encode(data)
			return
		}

		// Write content-type, statuscode, payload
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
}*/

// Add to cart

func (oc OrderController) AddToCart(w http.ResponseWriter, r *http.Request) {

	cart := Cart{}
	//username := r.FormValue("username")
	//username := "Rajvi"
	name := r.FormValue("name")
	price, _ := strconv.Atoi(r.FormValue("price"))
	calories, _ := strconv.Atoi(r.FormValue("calories"))
	quantity := 1
	
	//fmt.Println(name)
	//fmt.Println(price)
	//fmt.Println(calories)
	item := Item1{name,calories,price, quantity}
	if error := oc.session.DB("test").C("Cart").Find(bson.M{"Username": username}).One(&cart); error != nil {
		fmt.Println("errors:", error)
	
	var items []Item1
		
	items = append(items, item)
	cart = Cart{items, username}
	error = oc.session.DB("test").C("Cart").Insert(&cart)
	fmt.Println("errors:", error)
	} else {
		if error := oc.session.DB("test").C("Cart").Find(bson.M{"Items.Name": name}).One(&cart); error == nil {
			//means there's a cart with this item
			fmt.Println("means there's a cart with this item")
			fmt.Println(cart)

			for i := 0; i < len(cart.Items); i++ {
				if(cart.Items[i].Name == name){
					cart.Items[i].Quantity += 1
					break
				}
			}

			if error := oc.session.DB("test").C("Cart").Update(bson.M{"Username": username}, bson.M{"$set": bson.M{"Items": cart.Items}}); error != nil {
				fmt.Println(error)
			}

		} else {
			//means there isn't an item
			fmt.Println("means there isn't an item")
			cart.Items = append (cart.Items, item)
			if error := oc.session.DB("test").C("Cart").Update(bson.M{"Username": username}, bson.M{"$set": bson.M{"Items": cart.Items}}); error != nil {
				fmt.Println(error)
			}
		}

	}

	fmt.Println("Added")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	data := `{"status":"success","message":"Change successful"}`
	json.NewEncoder(w).Encode(data)
}

