package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
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
	conn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", host, port, username, dbname, password)

	var err error
	a.db, err = gorm.Open("postgres", conn)

	a.CreateTables()

	if err != nil {
		log.Fatal(err)
	}

	// define all the routes
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/", Home).Methods("GET")
	a.Router.HandleFunc("/signin", Signin).Methods("POST") //signin
	a.Router.HandleFunc("/signup", Signup).Methods("POST") //signup

	a.Router.HandleFunc("/current-user", GetCurrentUser).Methods("POST")
	a.Router.HandleFunc("/user/update", UpdateUser).Methods("POST")

	//a.Router.HandleFunc("/user/{id}", ViewUser).Methods("GET") //get user profile
	//a.Router.HandleFunc("/user/{id}", UpdateUser).Methods("PUT") //update user profile
	//a.Router.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE") //delete user profile

	a.Router.HandleFunc("/company", CreateCompany).Methods("POST")                 //create new company
	a.Router.HandleFunc("/company/{id}", ShowCompany).Methods("POST")              //get company
	a.Router.HandleFunc("/company/{id}/admin", AddAdmin).Methods("POST")           //add admin
	a.Router.HandleFunc("/company/{id}/admin-delete", DeleteAdmin).Methods("POST") //delete admin
	//a.Router.HandleFunc("/{id}/company", UpdateCompany).Methods("PUT") //update company
	//a.Router.HandleFunc("/{id}/company", DeleteCompany).Methods("DELETE") //delete company

	a.Router.HandleFunc("/company/{id}/job", AllJobs).Methods("GET")       //show all jobs posted by company
	a.Router.HandleFunc("/company/{id}/job", NewJob).Methods("POST")       //post new job
	a.Router.HandleFunc("/job/{id}/sort/{sortopt}", GetJob).Methods("GET") //get a job
	a.Router.HandleFunc("/job/{id}/apply", ApplyToJob).Methods("POST")     //apply to a job
	a.Router.HandleFunc("/job/{id}/update", UpdateJob).Methods("POST")     //apply to a job
	a.Router.HandleFunc("/job/{id}/decline", DeclineUser).Methods("POST")  //apply to a job
	a.Router.HandleFunc("/job/{id}/accept", AcceptUser).Methods("POST")    //apply to a job
	//a.Router.HandleFunc("/company/{id}/job/{job_id}", DeleteJob).Methods("DELETE") //delete posted jobs

	a.Router.HandleFunc("/search", JobSearch).Methods("GET")  //query for job
	a.Router.HandleFunc("/search", JobSearch).Methods("POST") //query for job

	a.Router.HandleFunc("/upload", TestFileUpload).Methods("POST")
	a.Router.HandleFunc("/upload-cv", UploadCV).Methods("POST")
	a.Router.HandleFunc("/upload-pp", UploadPP).Methods("POST")
	a.Router.PathPrefix("/files/").Handler(
		http.StripPrefix("/files/", http.FileServer(http.Dir("/home/tahsin/go/src/github.com/tahsinrahman/se_project_apiserver/files"))),
	)

}

// it listens on the given address
func (a *App) Run(addr string) {
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	})
	handler := c.Handler(a.Router)
	http.ListenAndServe(addr, handler)
}

// create all necessary tables
func (a *App) CreateTables() {
	a.db.AutoMigrate(&User{})
	a.db.AutoMigrate(&Company{})
	a.db.AutoMigrate(&Job{})
	a.db.AutoMigrate(&Tag{})
	a.db.AutoMigrate(&UserJobStatus{})
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
