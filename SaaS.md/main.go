package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
    "gopkg.in/mgo.v2"
    "time"
    "os"
    "gopkg.in/mgo.v2/bson"



)

// MongoDB Config
var mongodb_server = "mongodb"
var mongodb_database = "cmpe281"
var mongodb_collection = "starbucks"


// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func connectToMongo() bool {
    ret := false
    fmt.Println("enter main - connecting to mongo")

    // tried doing this - doesn't work as intended
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Detected panic")
            var ok bool
            err, ok := r.(error)
            if !ok {
                fmt.Printf("pkg:  %v,  error: %s", r, err)
            }
        }
    }()

    maxWait := time.Duration(5 * time.Second)
    session, sessionErr := mgo.DialWithTimeout("localhost:27017", maxWait)
    if sessionErr == nil {
        session.SetMode(mgo.Monotonic, true)
        coll := session.DB("MyDB").C("MyCollection")
        if ( coll != nil ) {
            fmt.Println("Got a collection object")
            ret = true
        }
    } else { // never gets here
        fmt.Println("Unable to connect to local mongo instance!")
    }
    return ret
}

type Device struct {
    Id           bson.ObjectId `bson:"_id" json:"_id,omitempty"`
    Name       string        `bson:"Name" json:"name"`
    Password    string        `bson:"Password" json:"password"`
    CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
    ModifiedAt   time.Time     `bson:"modifiedAt" json:"modifiedAt"`

}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	fmt.Println("Name")

	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {

	name := request.FormValue("name")
	pass := request.FormValue("password")
	fmt.Println("Name",name)
	fmt.Println("Pass", pass)
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials ..
    doc := Device{
        Id:           bson.NewObjectId(),
        Name:       name,
        Password:   pass,
        CreatedAt:    time.Now(),
        ModifiedAt:   time.Now(),
    }
     if ( connectToMongo() ) {
                fmt.Println("Connected")
                session, err := mgo.Dial("localhost")
                	if err != nil {
                		fmt.Printf("dial fail %v\n", err)
                		os.Exit(1)
                	}
                	defer session.Close()

                	//error check on every access
                	session.SetSafe(&mgo.Safe{})

                	//get collection
                	collection := session.DB(mongodb_database).C(mongodb_collection)
                	 err = collection.Insert(doc)
                            if err != nil {
                                panic(err)
                            }
            }


		setSession(name, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// index page

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// internal page

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// server main method

var router = mux.NewRouter()

func main() {
	fmt.Println("Name")




	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
