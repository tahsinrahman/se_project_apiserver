package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// user table
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
}

// handler for api/signup
// first checks if a user is already logged in or not
// if logged in redirect to home
// else check if username/email exits, if yes then return "username/email exists"
// register user and update db
func Signup(w http.ResponseWriter, r *http.Request) {
	// checks if a user is already logged in or not
	// if logged in redirect to home
	checkLoggedIn(w, r)

	// get user from request body
	user, err := getUserFromReq(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "empty username", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		http.Error(w, "empty password", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "empty Email", http.StatusBadRequest)
		return
	}
	if user.Name == "" {
		http.Error(w, "empty Name", http.StatusBadRequest)
		return
	}

	// insert into db
	err = app.db.Create(user).Error
	if err != nil {
		log.Println("1")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// write json response
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

// handler for api/signin
// first checks if a user has already logged in or not
// if logged in redirect to home
// else check if username exits, if not then return "user not found"
// check password, if do not match then return "invalid password"
// else login user and redirect to home
func Signin(w http.ResponseWriter, r *http.Request) {
	// check if user logged in
	checkLoggedIn(w, r)

	// check if username exists

	// check if password matches

	// if everything is ok, then generate a new token
}

// get jwt token from request header
// heder is like this "Authorization: Bearer TOKEN"
// so we need to split the authorization header to get the token
func checkTokenFromReq(r *http.Request) string {
	header := r.Header.Get("Authorization")
	token := strings.Split(header, " ")
	if len(token) > 1 {
		return token[1]
	}
	return ""
}

// varifies a jwt token
func verifyToken(reqToken string) bool {
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err == nil && token.Valid {
		return true
	}
	return false
}

// get user from request
func getUserFromReq(r *http.Request) (*User, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var user User
	err := decoder.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, err
}

// check if user logged in
// if logged in, redirect to home
func checkLoggedIn(w http.ResponseWriter, r *http.Request) {
	// get token from request header
	token := checkTokenFromReq(r)
	ok := verifyToken(token)

	// if token is valid then user is logged in
	//redirect to home
	if ok {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}
