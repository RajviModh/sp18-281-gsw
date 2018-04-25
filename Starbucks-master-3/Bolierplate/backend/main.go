package main

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	time2 "time"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)
// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))


func setSession(userName string, response http.ResponseWriter) {
	fmt.Println(userName)

	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:"session",
			Value: encoded,
			Path:"/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:"session",
		Value:"",
		Path:"/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

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

type Links struct {
	Payment string `json:"payment,omitempty"`
	Order   string `json:"order,omitempty"`
}

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

type Cart struct {
	CartId   string  `json:"id" bson:"_id"`
	Items    []Item1 `json:"items" bson:"Items"`
	Username string  `json:"username" bson:"Username"`
}


type User struct {
	UserName  string  `json:"username" bson:"username"`
	Password  string  `json:"password" bson:"password"`
	FirstName string  `json:"firstname" bson:"firstname"`
	LastName  string  `json:"lastname" bson:"lastname"`
	Credits   int `json:"credits" bson:"credits"`
	Location  string  `json:"location" bson:"location"`
}

// OrderController represents the controller for operating on the Order resource
type OrderController struct {
	session *mgo.Session
}

// NewOrderController provides a reference to a OrderController with provided mongo session
func NewOrderController(mgoSession *mgo.Session) *OrderController {
	return &OrderController{mgoSession}
}

type IndexController struct {
	session *mgo.Session
}

// NewOrderController provides a reference to a OrderController with provided mongo session

func NewIndexController(mgoSession *mgo.Session) *IndexController {
	return &IndexController{mgoSession}
}

func (sp IndexController) index(w http.ResponseWriter, r *http.Request) {
	var options []Item

	iter := sp.session.DB("cmpe281").C("User").Find(nil).Iter()
	fmt.Println("Inside iter")
	result := Item{}
	for iter.Next(&result) {
		options = append(options, result)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&options)

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

//ping resource function
func (oc OrderController) PingOrderResource(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pinging Order Resource")
}


//--------------------------------------------------------------Main Function----------------------------------------------------------------
func main() {

	r := mux.NewRouter()

	// Get a UserController instance
	oc := NewOrderController(getSession())

	// Get a UserController instance
	ic := NewIndexController(getSession())

	sp := NewSignUpController(getSession())

	lc := NewLoginController(getSession())

	lg := NewLogOutController(getSession())


	
	r.HandleFunc("/", ic.index).Methods("GET")
	r.HandleFunc("/signup", ic.index).Methods("GET")
	r.HandleFunc("/submitSignUp", sp.signup).Methods("POST")
	r.HandleFunc("/login", lc.login).Methods("POST")
	r.HandleFunc("/logout", lg.logout).Methods("GET")

	r.Methods("OPTIONS").HandlerFunc(IgnoreOption)

	r.HandleFunc("/ping", oc.PingOrderResource)

	fmt.Println("serving on port" + GetPort())

	http.ListenAndServe(GetPort(), r)

}

func IgnoreOption(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func getSession() (s *mgo.Session) {
	// Connect to local mongodb
	s, _ = mgo.Dial("127.0.0.1:27017")
	fmt.Println(s)
	return s
}
