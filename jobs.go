package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Job struct {
	ID          uint             `gorm:"primary_key"`
	CompanyID   *uint            `json:"companyID"`
	CompanyName *string          `json:"companyName"`
	Title       *string          `json:"title"`
	Description *string          `json:"description"`
	Field       *string          `json:"field"`
	Location    *string          `json:"location"'`
	Salary      *string          `json:"salary"'`
	Experience  *string          `json:"experience"`
	JobType     *string          `json:"type"`
	Requirement *string          `json:"requirement"`
	Deadline    *string          `json:"deadline"`
	Vacancy     *int             `json:"vacancy"`
	Tags        *string          `json:"tag"`
	TagList     []*Tag           `json:"-" gorm:"many2many:job_tags"`
	PostedAt    *time.Time       `json:"postedAt"`
	Applicants  []*UserJobStatus `json:"applicants"`
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
	//	nilEntry := checkNilEntryJob(job)
	//	if nilEntry != "" {
	//		resp["message"] = "empty" + nilEntry
	//		writeResp(w, http.StatusBadRequest, resp)
	//		return
	//	}

	job.CompanyID = &company.ID
	job.CompanyName = company.Name

	// insert into db
	if err = app.db.Create(job).Error; err != nil {
		resp["message"] = "error while inserting in db => " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if job.Tags != nil {
		tags := strings.Split(*job.Tags, ",")

		for _, tagName := range tags {
			if tagName[0] == ' ' {
				tagName = tagName[1:]
			}
			log.Println(tagName)
			tagName = strings.ToLower(tagName)

			var tag Tag
			err := app.db.Where("name = ?", tagName).First(&tag).Error
			if err != nil {
				log.Println("**************", err)
				tag.Name = &tagName
				app.db.Create(&tag)
			}
			log.Println(tag.ID, *tag.Name)

			app.db.Model(&job).Association("TagList").Append(&tag)
		}
	}

	resp["job"] = *job
	resp["message"] = "successfully posted job"
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

	log.Println(vars["sortopt"])
	log.Printf("%T\n", vars["sortopt"])

	tmp := vars["sortopt"]
	if tmp == "" {
		tmp = "cg"
	}

	resp["message"] = "success"
	resp["job"] = job.Response(tmp)
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

	app.db.Model(&job).Association("TagList").Clear()
	if job.Tags != nil {
		tags := strings.Split(*job.Tags, ",")

		for _, tagName := range tags {
			if tagName[0] == ' ' {
				tagName = tagName[1:]
			}
			tagName = strings.ToLower(tagName)
			log.Println(tagName)

			var tag Tag
			err := app.db.Where("name = ?", tagName).First(&tag).Error
			if err != nil {
				log.Println("**************", err)
				tag.Name = &tagName
				app.db.Create(&tag)
			}
			log.Println(tag.ID, *tag.Name)

			app.db.Model(&job).Association("TagList").Append(&tag)
		}
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

	// check if user already applied here
	if err = app.db.Preload("Applicants").Find(&job).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	for _, applicant := range job.Applicants {
		if applicant.UserID == user.ID {
			resp["message"] = "already applied"
			writeResp(w, http.StatusBadRequest, resp)
			return

		}
	}

	// check if current time < deadline
	if time.Now().After(finddate(*job.Deadline)) {
		resp["message"] = "deadline crossed"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	var userJobStatus []*UserJobStatus
	if err = app.db.Model(user).Related(&userJobStatus).Error; err != nil {
		resp["message"] = "error updating struct " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	userJobStatus = append(userJobStatus, &UserJobStatus{
		UserID: user.ID,
		JobID:  job.ID,
		Status: "pending",
	})

	user.AppliedAt = userJobStatus

	if err = app.db.Save(user).Error; err != nil {
		log.Println("444444444444")
		resp["message"] = "2 error updating struct " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	//if err := app.db.Model(user).Association("AppliedAt").Append(&UserJobStatus{
	//	User:   user,
	//	Job:    &job,
	//	Status: "pending",
	//}).Error; err != nil {
	//	resp["message"] = "error updating struct " + err.Error()
	//	writeResp(w, http.StatusBadRequest, resp)
	//	return
	//}

	/*
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
	*/

	resp["message"] = "success"
	resp["job"] = job
	resp["user"] = user.Response()
	writeResp(w, http.StatusOK, resp)
}

func DeclineUser(w http.ResponseWriter, r *http.Request) {
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
	var job Job
	if err = app.db.Where("id = ?", vars["id"]).First(&job).Error; err != nil {
		resp["message"] = "1 " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	/*
		if err = app.db.Where("id = ?", ReqBody.User.ID).First(user).Error; err != nil {
			resp["message"] = "2 " + err.Error()
			writeResp(w, http.StatusBadRequest, resp)
			return
		}
	*/

	var status UserJobStatus
	if err = app.db.Where("user_id = ? AND job_id = ? ", ReqBody.User.ID, job.ID).First(&status).Error; err != nil {
		resp["message"] = "2 " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	status.Status = "declined"
	if err = app.db.Model(&job).Association("Applicants").Replace(status).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if err = app.db.Where("id = ?", ReqBody.User.ID).First(ReqBody.User).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	SendEmail(*ReqBody.User.Email, "Jobengine-Notification", "Sorry, you've been rejected for job "+*job.Title)

	resp["message"] = "successfully declined"
	writeResp(w, http.StatusOK, resp)

}

func AcceptUser(w http.ResponseWriter, r *http.Request) {
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
	var job Job
	if err = app.db.Where("id = ?", vars["id"]).First(&job).Error; err != nil {
		resp["message"] = "1 " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	var status UserJobStatus
	if err = app.db.Where("user_id = ? AND job_id = ? ", ReqBody.User.ID, job.ID).First(&status).Error; err != nil {
		resp["message"] = "2 " + err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if err = app.db.Where("id = ?", ReqBody.User.ID).First(ReqBody.User).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	status.Status = "accepted"
	if err = app.db.Model(&job).Association("Applicants").Replace(status).Error; err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	SendEmail(*ReqBody.User.Email, "Jobengine-Notification", "Congratulations, you've been primariliy accepted for job "+*job.Title+" at "+*job.CompanyName+"! Take preperation for interview, we'll get back to you soon!")

	resp["message"] = "successfully accepted"
	writeResp(w, http.StatusOK, resp)
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

// response user without password field
func (user *Job) Response(sortopt string) (resp map[string]interface{}) {
	js, _ := json.Marshal(user)
	json.Unmarshal(js, &resp)
	tmp := resp["applicants"]
	delete(resp, "applicants")

	switch t := tmp.(type) {
	case []interface{}:
		type UserStatus struct {
			User   User   `json:"user"`
			Status string `json:"status"`
		}

		var users []UserStatus
		for _, j := range t {
			var user User
			app.db.Where("id = ?", j.(map[string]interface{})["UserID"]).First(&user)
			//if j.(map[string]interface{})["Status"]
			if j.(map[string]interface{})["Status"].(string) != "declined" {
				users = append(users, UserStatus{user, j.(map[string]interface{})["Status"].(string)})
			}
			//resp["appliedAt"] = append(resp["appliedAt"], job)
		}
		// sort by cgpa
		sort.Slice(users, func(i, j int) bool {
			if sortopt == "cgpa" {
				x, _ := strconv.ParseFloat(users[i].User.CG, 64)
				y, _ := strconv.ParseFloat(users[j].User.CG, 64)
				return x > y
			}
			x, _ := strconv.Atoi(*users[i].User.Experience)
			y, _ := strconv.Atoi(*users[j].User.Experience)
			return x > y
		})

		resp["applicants"] = users
	}

	return resp
}

func finddate(date string) time.Time {
	startDate := strings.Split(date, "/")

	year, _ := strconv.Atoi(startDate[0])
	month, _ := strconv.Atoi(startDate[1])
	day, _ := strconv.Atoi(startDate[2])
	deadline := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return deadline
}
