package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

// user table
type User struct {
	ID          uint             `gorm:"primary_key"`
	FirstName   *string          `json:"firstname" gorm:"not null"`
	LastName    *string          `json:"lastname"`
	Username    *string          `json:"username" gorm:"not null;unique"`
	Name        *string          `json:"name"`
	Password    *string          `json:"password" gorm:"not null"`
	Email       *string          `json:"email" gorm:"not null;unique"`
	CompanyID   uint             `json:"adminCompanyID"`
	Admin       *string          `json:"admin"`
	DateOfBirth *string          `json:dataOfBirth`
	Location    *string          `json:"location"`
	CurrentJob  *string          `json:"currentJob"`
	Designation *string          `json:"designation"`
	Street      *string          `json:"street"`
	State       *string          `json:"state"`
	Zip         *string          `json:"zip"`
	Country     *string          `json:"country"`
	Experience  *string          `json:"experience"`
	Description *string          `json:"description"`
	AppliedAt   []*UserJobStatus `json:"appliedAt"`
	CG          string           `json:"cgpa"`
	//Education   string    `json:education`
	//about-yourself
	// cv-pdf
	//photo
	// current job position
	// accepted_at
	// rejected_at
}

// get current user who is logged in
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["user"] = nil

	// checks if a user is already logged in or not
	// collect and varify token form request header
	ReqBody, err := checkLoggedIn(w, r)
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
	}

	log.Println(ReqBody.UserDB)
	user := ReqBody.UserDB
	if user != nil {
		if err := app.db.Preload("AppliedAt").First(user).Error; err != nil {
			log.Println("error here", err.Error())
		}
		log.Println("length ==========>", len(user.AppliedAt))

		resp["user"] = user.Response()
		writeResp(w, http.StatusOK, resp)
	} else {
		writeResp(w, http.StatusUnauthorized, resp)
	}
}

func createFile(path string, file *multipart.FileHeader) error {
	log.Println("filepath =========== ", path)
	newInputFile, err := os.Create(path)

	if err != nil {
		return err
	}
	defer newInputFile.Close()

	f, err := file.Open()
	if err != nil {
		return nil
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}

	_, err = newInputFile.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func TestFileUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("uploading")
	_, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("error uploading file", err.Error())
		return
	}
	if err = createFile("files/myfile", handler); err != nil {
		log.Println("error 2 uploading file", err.Error())
		return
	}
}

func UploadCV(w http.ResponseWriter, r *http.Request) {
	log.Println("uploading cv")
	resp := make(map[string]interface{})

	_, handler, err := r.FormFile("file")
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	token := r.FormValue("token")
	log.Println("111111111122222222222", r.FormValue("token"))

	user := verifyToken(token)

	if user == nil {
		resp["message"] = "not logged in"
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	// not logged in
	if user == nil {
		resp["message"] = "not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	if err := createFile("files/cv_"+strconv.Itoa(int(user.ID)), handler); err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}
}

func UploadPP(w http.ResponseWriter, r *http.Request) {
	log.Println("uploading cv")
	resp := make(map[string]interface{})

	_, handler, err := r.FormFile("file")
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	token := r.FormValue("token")
	log.Println("111111111122222222222", r.FormValue("token"))

	user := verifyToken(token)

	if user == nil {
		resp["message"] = "not logged in"
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	// not logged in
	if user == nil {
		resp["message"] = "not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	if err := createFile("files/pp_"+strconv.Itoa(int(user.ID)), handler); err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	ReqBody, err := checkLoggedIn(w, r)
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	user := ReqBody.UserDB

	// not logged in
	if user == nil {
		resp["message"] = "not logged in"
		writeResp(w, http.StatusUnauthorized, resp)
		return
	}

	user = ReqBody.User
	if user.Password == nil {
		user.Password = ReqBody.UserDB.Password
	}

	//	file, header, err := r.FormFile("pic")
	//	if err != nil {
	//		resp["message"] = "error updating user" + err.Error()
	//		writeResp(w, http.StatusInternalServerError, resp)
	//	}
	//	defer file.Close()

	if err = app.db.Save(user).Error; err != nil {
		resp["message"] = "error updating user" + err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	//err = createFile("./", header)
	//if err != nil {
	//	resp["message"] = "error loading pic" + err.Error()
	//	writeResp(w, http.StatusInternalServerError, resp)
	//	return
	//}

	resp["message"] = "successfully updated user"
	resp["user"] = *user
	writeResp(w, http.StatusOK, resp)
}

// handler for api/signup
func Signup(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	ReqBody, err := checkLoggedIn(w, r)
	log.Println("signin", ReqBody)
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}
	user := ReqBody.UserDB

	// already logged in
	if user != nil {
		resp = map[string]interface{}{"user": *user, "message": "already logged in"}
		writeResp(w, http.StatusOK, resp)
		return
	}

	if err != nil {
		resp = map[string]interface{}{"user": nil, "message": err.Error()}
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	user = ReqBody.User
	if user == nil {
		resp["message"] = "empty user"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Username == nil {
		resp = map[string]interface{}{"user": nil, "message": "empty username"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Password == nil {
		resp = map[string]interface{}{"user": nil, "message": "empty password"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// check if username exists
	var userInDB User
	app.db.Where("username = ?", user.Username).First(&userInDB)

	if userInDB.Username != nil {
		resp = map[string]interface{}{"user": nil, "message": "username exists"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// check nil entry
	nilEntry := checkNilEntryUser(*user)
	if nilEntry != "" {
		resp["message"] = "empty " + nilEntry
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	name := *user.FirstName
	if user.LastName != nil {
		name += " " + *user.LastName
	}
	user.Name = &name
	//log.Println(*user.Name)

	log.Println(user.ID)
	log.Println(user.Username)
	log.Println(user.Email)

	// insert into db
	err = app.db.Create(user).Error
	if err != nil {
		resp = map[string]interface{}{"user": nil, "message": err.Error()}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	cmd := exec.Command("cp", "files/default", fmt.Sprintf("files/pp_"+strconv.Itoa(int(user.ID))))
	cmd.Run()

	resp = map[string]interface{}{"user": user.Response(), "message": "successfully registered"}
	writeResp(w, http.StatusCreated, resp)
}

func CheckShell() {
	cmd := exec.Command("cp", "files/default", "files/tmp_pp")
	err := cmd.Run()
	//stdout, err := cmd.Output()

	if err != nil {
		log.Println(err.Error())
		return
	}

	//	log.Println(string(stdout))
}

// handler for api/signin
// first checks if a user has already logged in or not
// if not check if username exits, if not then return "user not found"
// check password, if do not match then return "invalid password"
// else log-in user
func Signin(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	ReqBody, err := checkLoggedIn(w, r)
	log.Println("signin", ReqBody)
	if err != nil {
		resp["message"] = err.Error()
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}
	user := ReqBody.UserDB

	// already logged in
	if user != nil {
		resp = map[string]interface{}{"user": *user, "message": "already logged in"}
		writeResp(w, http.StatusOK, resp)
		return
	}

	if err != nil {
		resp = map[string]interface{}{"user": nil, "message": err.Error()}
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	user = ReqBody.User
	if user == nil {
		resp["message"] = "empty user"
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Username == nil {
		resp = map[string]interface{}{"user": nil, "message": "empty username"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Password == nil {
		resp = map[string]interface{}{"user": nil, "message": "empty password"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	//	if *user.Username == "" {
	//		resp = map[string]interface{}{"user": nil, "message": "empty username"}
	//		writeResp(w, http.StatusBadRequest, resp)
	//		return
	//	}
	//	if *user.Password == "" {
	//		resp = map[string]interface{}{"user": nil, "message": "empty password"}
	//		writeResp(w, http.StatusBadRequest, resp)
	//		return
	//	}

	// check if username exists
	log.Println("username", *user.Username)
	log.Println("password", *user.Password)
	var userInDB User
	if err = app.db.Where("username = ?", user.Username).First(&userInDB).Error; err != nil {
		resp = map[string]interface{}{"user": nil, "message": "empty password"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	if userInDB.Username == nil {
		resp = map[string]interface{}{"user": nil, "message": "username not found"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// check if password matches
	if *userInDB.Password != *user.Password {
		resp = map[string]interface{}{"user": nil, "message": "invalid password"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// everything is ok, now generate a new token
	token, err := GenerateToken(*user.Username)
	if err != nil {
		resp = map[string]interface{}{"user": nil, "message": "error generating token" + err.Error()}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// return token
	resp = map[string]interface{}{"user": userInDB.Response(), "token": *token, "message": "successfully logged in"}
	writeResp(w, http.StatusOK, resp)
}

func checkNilEntryUser(user User) string {
	if user.Username == nil {
		return "username"
	}
	if user.Password == nil {
		return "password"
	}
	if user.Email == nil {
		return "email"
	}
	if user.FirstName == nil {
		return "name"
	}
	log.Println(*user.Username)
	log.Println(*user.Password)
	if *user.Username == "" {
		return "username"
	}
	if *user.Password == "" {
		return "password"
	}
	if *user.Email == "" {
		return "email"
	}
	if *user.FirstName == "" {
		return "name"
	}
	return ""
}

// response user without password field
func (user *User) Response() (resp map[string]interface{}) {
	js, _ := json.Marshal(user)
	json.Unmarshal(js, &resp)
	delete(resp, "password")
	tmp := resp["appliedAt"]
	delete(resp, "appliedAt")

	switch t := tmp.(type) {
	case []interface{}:
		type JobStatus struct {
			Job    Job    `json:"job"`
			Status string `json:"status"`
		}

		var jobs []JobStatus
		for _, j := range t {
			jj := j.(map[string]interface{})
			log.Println(jj["JobID"])
			var job Job
			app.db.Where("id = ?", j.(map[string]interface{})["JobID"]).First(&job)
			jobs = append(jobs, JobStatus{job, j.(map[string]interface{})["Status"].(string)})
		}
		resp["appliedAt"] = jobs
	}

	return resp
}
