package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Job struct {
	ID          uint       `gorm:"primary_key"`
	CompanyID   *uint      `json:"companyID"`
	CompanyName *string    `json:"companyName"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Field       *string    `json:"field"`
	Location    *string    `json:"location"'`
	Salary      *string    `json:"salary"'`
	Experience  *string    `json:"experience"`
	JobType     *string    `json:"type"`
	Requirement *string    `json:"requirement"`
	Deadline    *time.Time `json:"deadline"`
	Vacancy     *int       `json:"vacancy"`
	Tags        []Tag      `json:"tags" gorm:"many2many:job_tags;"`
	PostedAt    *time.Time `json:"postedAt"`
	Applicants  []User     `json:"applicants" gorm:"many2many:user_jobs;"`
}

func NewJob(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["job"] = nil

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

	// get job-info from req body
	job := ReqBody.Job

	// check nil entry
	nilEntry := checkNilEntryJob(job)
	if nilEntry != "" {
		resp["message"] = "empty" + nilEntry
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	job.CompanyID = &company.ID
	job.CompanyName = company.Name

	// insert into db
	if err = app.db.Create(job).Error; err != nil {
		resp["message"] = "error while inserting in db => " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	resp["job"] = *job
	resp["message"] = "successfully registered company"
	writeResp(w, http.StatusCreated, resp)
}

func AllJobs(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	vars := mux.Vars(r)
	tmp, err := strconv.Atoi(vars["id"])
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}
	id := uint(tmp)

	// check if company-id exists
	var company Company
	err = app.db.Where("id = ?", vars["id"]).First(&company).Error
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	log.Println("1111111111111")
	log.Println("id")

	// list all the jobs by company id
	var jobs []Job
	if err = app.db.Where("company_id = ?", id).Find(&jobs).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	resp["message"] = "success"
	resp["jobs"] = jobs
	writeResp(w, http.StatusOK, resp)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	vars := mux.Vars(r)

	// check if company-id exists
	var job Job
	err := app.db.Where("id = ?", vars["id"]).First(&job).Error
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if err = app.db.Preload("Applicants").First(&job).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	resp["message"] = "success"
	resp["job"] = job
	writeResp(w, http.StatusOK, resp)
}

func UpdateJob(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["job"] = nil

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

	// check if job-id exists
	err = app.db.Where("id = ?", vars["id"]).Error
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// get job-info from req body
	job := ReqBody.Job

	// check nil entry
	nilEntry := checkNilEntryJob(job)
	if nilEntry != "" {
		resp["message"] = "empty" + nilEntry
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if err = app.db.Save(job).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	resp["message"] = "successfully updated job"
	resp["job"] = job
	writeResp(w, http.StatusOK, resp)
}

func ApplyToJob(w http.ResponseWriter, r *http.Request) {
	log.Println("in apply")
	resp := make(map[string]interface{})

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

	// check if job-id exists
	var job Job
	err = app.db.Where("id = ?", vars["id"]).First(&job).Error
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// check if current time < deadline
	//	if time.Now().After(*job.Deadline) {
	//		resp["message"] = "deadline crossed"
	//		writeResp(w, http.StatusBadRequest, resp)
	//		return
	//	}

	// update user profile
	if err = app.db.Preload("AppliedAt").First(user).Error; err != nil {
		resp["message"] = "error updating user struct " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if err = app.db.Model(user).Association("AppliedAt").Append(job).Error; err != nil {
		resp["message"] = "11 " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	//	user.AppliedAt = append(user.AppliedAt, job)
	//	if err = app.db.Update(user).Error; err != nil {
	//		resp["message"] = "11 " + err.Error()
	//		writeResp(w, http.StatusBadRequest, resp)
	//		return
	//	}

	// update job
	if err = app.db.Preload("Applicants").First(&job).Error; err != nil {
		resp["message"] = "22 " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if err = app.db.Model(&job).Association("Applicants").Append(*user).Error; err != nil {
		resp["message"] = "33 " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	//job.Applicants = append(job.Applicants, *user)
	//if err = app.db.Update(&job).Error; err != nil {
	//	resp["message"] = "33 " + err.Error()
	//	writeResp(w, http.StatusBadRequest, resp)
	//	return
	//}

	resp["message"] = "success"
	resp["job"] = job
	resp["user"] = user.Response()
	writeResp(w, http.StatusOK, resp)
}

func Decline(w http.ResponseWriter, r *http.Request) {
	// remove user from job []applicants

	// remove
}

func checkNilEntryJob(job *Job) string {
	// check for nil entries
	// TODO:DRY

	if job == nil {
		return "job"
	}

	if job.Title == nil || *job.Title == "" {
		return "title"
	}
	if job.Description == nil || *job.Description == "" {
		return "Description"
	}
	if job.Field == nil || *job.Field == "" {
		return "Field"
	}
	if job.Location == nil || *job.Location == "" {
		return "Location"
	}
	if job.Experience == nil || *job.Experience == "" {
		return "Experience"
	}
	if job.JobType == nil || *job.JobType == "" {
		return "JobType"
	}
	return ""
}
