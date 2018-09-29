package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// company table
type Company struct {
	ID          uint    `gorm:"primary_key"`
	Name        *string `json:"name" gorm:"not null"`
	Description *string `json:"description"`
	Email       *string `json:"email"`
	Admin       []User  `json:"admin" gorm:"foreignkey:CompanyID"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone"`
	Jobs        []Job   `json:"jobs" gorm:"foreignkey:CompanyID"`
	//Logo *string `json:"email"`
	//openJobs    []Jobs
}

// handler for api/create-company
// first check if user is logged in or not
func CreateCompany(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["user"] = nil
	resp["company"] = nil

	// checks if a user is already logged in or not
	// collect and varify token form request header
	ReqBody, err := checkLoggedIn(w, r)

	if err != nil {
		resp["message"] = "err" + err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	user := ReqBody.UserDB

	// not logged in
	if user == nil {
		resp["message"] = "user not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	company := ReqBody.Company

	if company == nil {
		resp["message"] = "empty company"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	// check for nil entries
	if company.Name == nil {
		resp["message"] = "empty name"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	/*
		if company.Description == nil {
			resp["message"] = "empty description"
			writeResp(w, http.StatusBadRequest, resp)
			return
		}
	*/

	// check if current user admin/hr of a company
	if user.Admin != nil {
		resp["message"] = "can't be in more than one company"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if err = UpdateUserAdmin(user, company); err != nil {
		resp["message"] = "error updating user profile"
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	// add current user to the admin list
	company.Admin = append(company.Admin, *user)

	// insert into db
	err = app.db.Create(company).Error
	if err != nil {
		resp["message"] = "error while inserting in db => " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	resp["company"] = company.Response()
	resp["message"] = "successfully registered company"

	writeResp(w, http.StatusCreated, resp)
}

func ShowCompany(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["user"] = nil
	resp["company"] = nil

	// checks if a user is already logged in or not
	// collect and varify token form request header
	ReqBody, err := checkLoggedIn(w, r)

	if err != nil {
		resp["message"] = "err" + err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	user := ReqBody.UserDB

	// not logged in
	if user == nil {
		resp["message"] = "user not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	vars := mux.Vars(r)
	tmp, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}
	id := uint(tmp)

	// user is not admin
	if user.CompanyID != id {
		resp["message"] = "user is not admin/hr"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	// check if company-id exists
	var company Company
	err = app.db.Where("id = ?", vars["id"]).First(&company).Error
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	var admins []User
	app.db.Model(&company).Related(&admins)

	company.Admin = admins

	resp["message"] = "successfull"
	resp["company"] = company.Response()
	resp["user"] = user.Response()
	writeResp(w, http.StatusOK, resp)
}

// api/company/{id}/admin POST
func AddAdmin(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["user"] = nil
	resp["company"] = nil

	// checks if a user is already logged in or not
	// collect and varify token form request header
	ReqBody, err := checkLoggedIn(w, r)

	if err != nil {
		resp["message"] = "err" + err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	user := ReqBody.UserDB

	// not logged in
	if user == nil {
		resp["message"] = "user not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	vars := mux.Vars(r)
	tmp, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}
	id := uint(tmp)

	// user is not admin
	if user.CompanyID != id {
		resp["message"] = "user is not admin/hr"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	// check if company-id exists
	var company Company
	err = app.db.Where("id = ?", vars["id"]).First(&company).Error
	if err != nil {
		resp["message"] = "company id doesn't exist"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	log.Println("11111111111111")

	user = ReqBody.User

	if user == nil {
		resp["message"] = "empty email"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if err = app.db.Where("email = ?", *user.Email).First(user).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	//update company
	var admins []User
	app.db.Model(&company).Related(&admins)
	admins = append(admins, *user)

	company.Admin = admins

	//update user profile
	if err = UpdateUserAdmin(user, &company); err != nil {
		resp["message"] = "error updating user profile"
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	resp["message"] = "successfull"
	resp["company"] = company.Response()
	resp["user"] = user.Response()
	writeResp(w, http.StatusOK, resp)
}

func DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["user"] = nil
	resp["company"] = nil

	// checks if a user is already logged in or not
	// collect and varify token form request header
	ReqBody, err := checkLoggedIn(w, r)

	if err != nil {
		resp["message"] = "err" + err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	user := ReqBody.UserDB

	// not logged in
	if user == nil {
		resp["message"] = "user not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	vars := mux.Vars(r)
	tmp, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}
	id := uint(tmp)

	// user is not admin
	if user.CompanyID != id {
		resp["message"] = "user is not admin/hr"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	// check if company-id exists
	var company Company
	err = app.db.Where("id = ?", vars["id"]).First(&company).Error
	if err != nil {
		resp["message"] = "company id doesn't exist"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	user = ReqBody.User

	if user == nil {
		resp["message"] = "empty user"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Email == nil {
		resp["message"] = "empty user"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if err = app.db.Where("email = ?", *user.Email).First(user).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if user.Admin != nil && *user.Admin != *company.Name {
		resp["message"] = "user is not currently admin"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	//update company
	var admins []User
	app.db.Model(&company).Related(&admins)

	if len(admins) == 1 {
		resp["message"] = "you're the only admin"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	indexToDelete := -1
	for id, admin := range admins {
		if *admin.Email == *user.Email {
			indexToDelete = id
			break
		}
	}

	if indexToDelete == -1 {
		resp["message"] = "user not in admin list"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	admins = append(admins[:indexToDelete], admins[indexToDelete+1:]...)

	company.Admin = admins

	user.CurrentJob = nil
	user.CompanyID = 0
	user.Admin = nil

	if err := app.db.Save(user).Error; err != nil {
		resp["message"] = "error updating user profile"
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	resp["selfDestruct"] = false
	if *user.Email == *ReqBody.UserDB.Email {
		resp["selfDestruct"] = true
	}

	resp["message"] = "successfull"
	resp["company"] = company.Response()
	resp["user"] = user.Response()
	writeResp(w, http.StatusOK, resp)
}

func UpdateUserAdmin(user *User, company *Company) error {
	//update user profile
	user.CurrentJob = company.Name
	user.CompanyID = company.ID
	user.Admin = company.Name

	if err := app.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// handler for api/update-company
func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	// checks if a user is already logged in or not
	ReqBody, err := checkLoggedIn(w, r)
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	user := ReqBody.User
	resp["user"] = user.Response()
	resp["company"] = nil

	// not logged in
	if user == nil {
		resp["message"] = "user not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	vars := mux.Vars(r)

	// check if company-id exists
	var tmpCompany Company
	err = app.db.Where("id = ?", vars["id"]).First(&tmpCompany).Error
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// exits
	// now load admin and hr accounts
	app.db.Preload("Admin").First(&tmpCompany)

	// check if current user is admin
	flag := false
	for _, admin := range tmpCompany.Admin {
		if admin.ID == user.ID {
			flag = true
			break
		}
	}

	// current user is not admin
	if !flag {
		resp["message"] = "user is not admin"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	// user is admin, now get company info from req body
	company, err := getCompanyFromReq(r)
	if err != nil {
		resp["message"] = "error while getting info from req body => " + err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	// check for empty fields
	if company.Name == nil {
		resp["message"] = "empty name"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if company.Description == nil {
		resp["message"] = "empty description"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	company.ID = tmpCompany.ID

	// check for invalid entry
	emptyAdmin := true
	for index, admin := range company.Admin {
		err = app.db.Where(admin).First(&company.Admin[index]).Error
		if err != nil {
			resp["message"] = "admin not found => " + err.Error()
			writeResp(w, http.StatusBadRequest, resp)
			return
		}
		emptyAdmin = false
	}

	if emptyAdmin {
		resp["message"] = "empty admin list"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	var admins []User
	existsAdmin := make(map[uint]bool)

	for _, admin := range company.Admin {
		if !existsAdmin[admin.ID] {
			existsAdmin[admin.ID] = true
			admins = append(admins, admin)
		}
	}
	company.Admin = admins

	err = app.db.Save(company).Error
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}
	if err = app.db.Model(company).Association("admin").Replace(company.Admin).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	resp["company"] = company.Response()
	resp["message"] = "successfully updated"
	writeResp(w, http.StatusOK, resp)
}

// get company-info from request
func getCompanyFromReq(r *http.Request) (*Company, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var company Company
	err := decoder.Decode(&company)

	if err != nil {
		return nil, err
	}

	return &company, err
}

func (company *Company) Response() (resp map[string]interface{}) {
	js, _ := json.Marshal(company)
	json.Unmarshal(js, &resp)

	if resp["admin"] != nil {
		for index, _ := range resp["admin"].([]interface{}) {
			delete(resp["admin"].([]interface{})[index].(map[string]interface{}), "password")
		}
	}
	return resp
}
