package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// handler for api/create-company
// first check if user is logged in or not
func CreateCompany(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	// checks if a user is already logged in or not
	user := checkLoggedIn(w, r)
	resp["user"] = user.Response()
	resp["company"] = nil

	// not logged in
	if user == nil {
		resp["message"] = "user not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	// get company-infos from request body
	company, err := getCompanyFromReq(r)
	if err != nil {
		resp["message"] = "unable to get company info from req body => " + err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	// check for nil entries
	if company.Name == "" {
		resp["message"] = "empty name"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if company.Description == "" {
		resp["message"] = "empty description"
		writeResp(w, http.StatusBadRequest, resp)
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

// handler for api/update-company
func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	// checks if a user is already logged in or not
	user := checkLoggedIn(w, r)
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
	err := app.db.Where("id = ?", vars["id"]).First(&tmpCompany).Error
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// exits
	// now load admin and hr accounts
	app.db.Preload("Admin").First(&tmpCompany)
	app.db.Preload("HR").First(&tmpCompany)

	// check if current user is admin
	flag := false
	for _, admin := range tmpCompany.Admin {
		if admin == *user {
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
	if company.Name == "" {
		resp["message"] = "empty name"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if company.Description == "" {
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

	// check for invalid intry
	for index, HR := range company.HR {
		err = app.db.Where(HR).First(&company.HR[index]).Error
		if err != nil {
			resp["message"] = "hr account not found => " + err.Error()
			writeResp(w, http.StatusBadRequest, resp)
			return
		}
	}

	// check duplicate entry in admin and hr filed
	var admins []User
	existsAdmin := make(map[uint]bool)

	for _, admin := range company.Admin {
		if !existsAdmin[admin.ID] {
			existsAdmin[admin.ID] = true
			admins = append(admins, admin)
		}
	}
	company.Admin = admins

	var hrs []User
	existsHR := make(map[uint]bool)

	for _, hr := range company.HR {
		if !existsHR[hr.ID] && !existsAdmin[hr.ID] {
			existsHR[hr.ID] = true
			hrs = append(hrs, hr)
		}
	}
	company.HR = hrs

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
	if err = app.db.Model(company).Association("HR").Replace(company.HR).Error; err != nil {
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
