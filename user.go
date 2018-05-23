package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// user table
type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
}

// handler for api/signup
// first checks if a user is already logged in or not
// then check if username/email exits, if yes then return "username/email exists"
// register user and update db
func Signup(w http.ResponseWriter, r *http.Request) {
	// checks if a user is already logged in or not
	user := checkLoggedIn(w, r)

	var resp map[string]interface{}

	// already logged in
	if user != nil {
		resp = map[string]interface{}{"user": *user, "message": "already logged in"}
		writeResp(w, http.StatusOK, resp)
		return
	}

	// get user from request body
	var err error
	user, err = getUserFromReq(r)

	if err != nil {
		resp = map[string]interface{}{"user": nil, "message": err.Error()}
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	if user.Username == "" {
		resp = map[string]interface{}{"user": nil, "message": "empty username"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Password == "" {
		resp = map[string]interface{}{"user": nil, "message": "empty password"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Email == "" {
		resp = map[string]interface{}{"user": nil, "message": "empty email"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Name == "" {
		resp = map[string]interface{}{"user": nil, "message": "empty name"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// insert into db
	err = app.db.Create(user).Error
	if err != nil {
		resp = map[string]interface{}{"user": nil, "message": err.Error()}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	resp = map[string]interface{}{"user": user, "message": "successfully registered"}
	writeResp(w, http.StatusCreated, resp)
}

// handler for api/signin
// first checks if a user has already logged in or not
// if not check if username exits, if not then return "user not found"
// check password, if do not match then return "invalid password"
// else log-in user
func Signin(w http.ResponseWriter, r *http.Request) {
	// checks if a user is already logged in or not
	user := checkLoggedIn(w, r)

	var resp map[string]interface{}

	// already logged in
	if user != nil {
		resp = map[string]interface{}{"user": *user, "message": "already logged in"}
		writeResp(w, http.StatusOK, resp)
		return
	}

	// get user from request body
	var err error
	user, err = getUserFromReq(r)

	if err != nil {
		resp = map[string]interface{}{"user": nil, "message": err.Error()}
		writeResp(w, http.StatusInternalServerError, resp)
		return
	}

	if user.Username == "" {
		resp = map[string]interface{}{"user": nil, "message": "empty username"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}
	if user.Password == "" {
		resp = map[string]interface{}{"user": nil, "message": "empty password"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// check if username exists
	var userInDB User
	app.db.Where("username = ?", user.Username).First(&userInDB)

	if userInDB.Username == "" {
		resp = map[string]interface{}{"user": nil, "message": "username not found"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// check if password matches
	if userInDB.Password != user.Password {
		resp = map[string]interface{}{"user": nil, "message": "invalid password"}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// everything is ok, now generate a new token
	token, err := generateToken(r, user.Username)
	if err != nil {
		resp = map[string]interface{}{"user": nil, "message": "error generating token" + err.Error()}
		writeResp(w, http.StatusBadRequest, resp)
		return
	}

	// return token
	resp = map[string]interface{}{"user": userInDB, "token": *token, "message": "successfully logged in"}
	writeResp(w, http.StatusOK, resp)
}

func writeResp(w http.ResponseWriter, code int, resp map[string]interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

// check if user logged in
// if logged in, return the user who's logged in
func checkLoggedIn(w http.ResponseWriter, r *http.Request) *User {
	// get token from request header
	token := getTokenFromReq(r)

	// check if the token is valid or not
	// if valid return the user currently logged in
	user := verifyToken(token)
	return user
}

// generate jwt token
func generateToken(r *http.Request, username string) (*string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// set token claims
	claims["username"] = username
	claims["exp"] = strconv.FormatInt(time.Now().Add(time.Hour*24).Unix(), 10)

	// sign token with secret
	tokenString, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

// get jwt token from request header
// heder is like this "Authorization: Bearer TOKEN"
// so we need to split the authorization header to get the token
func getTokenFromReq(r *http.Request) string {
	header := r.Header.Get("Authorization")
	token := strings.Split(header, " ")
	if len(token) > 1 {
		return token[1]
	}
	return ""
}

// varifies a jwt token
func verifyToken(reqToken string) *User {
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check if claims are true
		username := claims["username"]

		// check if username exists
		var userInDB User
		app.db.Where("username = ?", username).First(&userInDB)

		// username not found
		if userInDB.Username == "" {
			return nil
		}

		// now check if current time is less than expiratin time
		unixTime, err := strconv.ParseInt(claims["exp"].(string), 10, 64)
		if err != nil {
			return nil
		}

		expirationTime := time.Unix(unixTime, 0)
		if expirationTime.After(time.Now()) {
			return &userInDB
		}

		return nil
	}

	return nil
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
