package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	macaron "gopkg.in/macaron.v1"
)

func InitializeRouter() {
	m := macaron.Classic()

	m.Use(macaron.Renderer()) //for html rendering

	m.Get("/", GetHome)

	//auth
	m.Post("/signin", SignIn)
	m.Post("/signup", SignUp)
	m.Get("/refresh-token", RefreshToken)

	//user
	m.Get("/user", GetUser)
	m.Put("/user/:user_id", UpdateUser)
	m.Delete("/user/:user_id", DeleteUser)

	//company
	m.Get("/company", GetCompany)
	m.Post("/company", NewCompany)
	m.Put("/company/:company_id", UpdateCompany)
	m.Delete("/company/:company_id", DeleteCompany)

	//job
	m.Post("/search", SearchJobs)
	m.Post("company/:company_id/job", NewJob)
	m.Put("company/:company_id/job/:job_id", UpdateJob)
	m.Delete("company/:company_id/job/:job_id", DeleteJob)

	//starting the server
	m.Run()
}

//homepage
func GetHome(ctx *macaron.Context) string {
	return "hello-world"
}
