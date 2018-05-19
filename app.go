package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// this is the base structure of the apiserver
// it contains two main components, router and database
type App struct {
	Router *mux.Router
	db     *gorm.DB
}

// it connects with db and defines all the routes
func (a *App) Initialize(host, port, dbname, username, password string) {
	// first connect with db
	conn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v", host, port, username, dbname, password)

	var err error
	a.db, err = gorm.Open("postgres", conn)

	a.createTables()

	if err != nil {
		log.Fatal(err)
	}

	// define all the routes
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/", Home).Methods("GET")
	a.Router.HandleFunc("/signin", Signin).Methods("POST")
	a.Router.HandleFunc("/signup", Signup).Methods("POST")
}

// it listens on the given address
func (a *App) Run(addr string) {
	http.ListenAndServe(addr, a.Router)
}

// create all necessary tables
func (a *App) createTables() {
	a.db.AutoMigrate(&User{})
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
