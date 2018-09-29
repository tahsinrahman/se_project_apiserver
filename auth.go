package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type ReqBody struct {
	Token     string   `json:"token"`
	User      *User    `json:"user"`
	Company   *Company `json:"company"`
	Job       *Job     `json:"job"`
	Search    *Search  `json:"search"`
	UserDB    *User
	CompanyDB *Company
	JobDB     *Company
}

type Search struct {
	Tag      string `json:"term"`
	Location string `json:"location"`
}

func writeResp(w http.ResponseWriter, code int, resp map[string]interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

// check if user logged in
// if logged in, return the user who's logged in
func checkLoggedIn(w http.ResponseWriter, r *http.Request) (*ReqBody, error) {
	// get token from request header
	ReqBody, err := getTknFromReq(r)
	if err != nil {
		return nil, err
	}

	// check if the token is valid or not
	// if valid return the user currently logged in
	user := verifyToken(ReqBody.Token)
	ReqBody.UserDB = user
	log.Println("checklogged in", ReqBody)
	return ReqBody, nil
}

// get user from request
func getTknFromReq(r *http.Request) (*ReqBody, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var token ReqBody
	err := decoder.Decode(&token)

	log.Println("get token from req", token, err)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// generate jwt token
func GenerateToken(username string) (*string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// set token claims
	claims["username"] = username
	//claims["exp"] = strconv.FormatInt(time.Now().Add(time.Hour*24).Unix(), 10)
	claims["exp"] = strconv.FormatInt(time.Now().Add(time.Hour*24).Unix(), 10)

	// sign token with secret
	tokenString, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

// get jwt token from request header
// heder is like this "Authorization: Bearer <TOKEN>"
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
	log.Println("verify ", reqToken)
	log.Println(reqToken)
	if reqToken == "" {
		log.Println("here")
		return nil
	}
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
		log.Println("verify username", username)
		app.db.Where("username = ?", username).First(&userInDB)

		// username not found
		if userInDB.Username == nil {
			log.Println("1 returning")
			return nil
		}

		// now check if current time is less than expiratin time
		unixTime, err := strconv.ParseInt(claims["exp"].(string), 10, 64)
		if err != nil {
			log.Println("2 returning")
			return nil
		}

		expirationTime := time.Unix(unixTime, 0)
		log.Println(expirationTime)
		if expirationTime.After(time.Now()) {
			return &userInDB
		}

		log.Println("3 returning")
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
