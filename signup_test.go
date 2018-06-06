package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestSignUp(t *testing.T) {
	droptable()

	var users = []UserTest{
		// ok
		UserTest{
			Req:  `{"name":"a", "username": "a", "email": "a", "password": "a"}`,
			Resp: `{"message":"successfully registered","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Code: 201,
		},
		// same username
		UserTest{
			Req:  `{"name":"a", "username": "a", "email": "b", "password": "b"}`,
			Resp: `{"message":"pq: duplicate key value violates unique constraint \"users_username_key\"","user":null}`,
			Code: 400,
		},
		// same email
		UserTest{
			Req:  `{"name":"b", "username": "b", "email": "a", "password": "b"}`,
			Resp: `{"message":"pq: duplicate key value violates unique constraint \"users_email_key\"","user":null}`,
			Code: 400,
		},
		// empty name
		UserTest{
			Req:  `{"name":"", "username": "c", "email": "c", "password": "c"}`,
			Resp: `{"message":"empty name","user":null}`,
			Code: 400,
		},
		// empty username
		UserTest{
			Req:  `{"name":"c", "username": "", "email": "c", "password": "c"}`,
			Resp: `{"message":"empty username","user":null}`,
			Code: 400,
		},
		// empty email
		UserTest{
			Req:  `{"name":"c", "username": "c", "email": "", "password": "c"}`,
			Resp: `{"message":"empty email","user":null}`,
			Code: 400,
		},
		// empty password
		UserTest{
			Req:  `{"name":"c", "username": "c", "email": "c", "password": ""}`,
			Resp: `{"message":"empty password","user":null}`,
			Code: 400,
		},
	}

	for _, test := range users {
		req, err := http.NewRequest("POST", "/signup", strings.NewReader(test.Req))
		if err != nil {
			log.Println("error in creating req => ", err.Error())
			continue
		}
		//debugReq(req)

		runTest(req, test, t)
	}
}

func TestSignIn(t *testing.T) {
	var tests []UserTest = []UserTest{
		// ok
		UserTest{
			Req:  `{"username": "a", "password": "a"}`,
			Resp: `{"message":"successfully logged in", "user":{"ID":1,"name":"a","username":"a","email":"a"}}`,
			Code: 200,
		},
		// empty username
		UserTest{
			Req:  `{"username": "", "password": "a"}`,
			Resp: `{"message":"empty username","user":null}`,
			Code: 400,
		},
		// empty password
		UserTest{
			Req:  `{"username": "a", "password": ""}`,
			Resp: `{"message":"empty password","user":null}`,
			Code: 400,
		},
		// invalid user
		UserTest{
			Req:  `{"username": "b", "password": "b"}`,
			Resp: `{"message":"username not found","user":null}`,
			Code: 400,
		},
		// password mismatch
		UserTest{
			Req:  `{"username": "a", "password": "b"}`,
			Resp: `{"message":"invalid password","user":null}`,
			Code: 400,
		},
	}
	for _, test := range tests {
		req, err := http.NewRequest("POST", "/signin", strings.NewReader(test.Req))
		if err != nil {
			log.Println("error in creating req => ", err.Error())
			continue
		}
		//debugReq(req)

		response := executeRequest(req)

		resp, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error in reading resp => ", err.Error())
			continue
		}

		log.Println("request => ", test.Req)
		log.Println("response => ", string(resp))

		checkSignIn(t, test.Code, response.Code, test.Resp, string(resp[:len(resp)-1]))
	}
}

type signinTest struct {
	User    User
	Message string
}

func checkSignIn(t *testing.T, expectedCode int, actualCode int, expectedResp string, actualResp string) {
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, actualCode)
	}
	var expected, actual signinTest

	err := json.Unmarshal([]byte(expectedResp), &expected)

	if err != nil {
		log.Println("error unmarshal expected response => ", err.Error())
		return
	}

	err = json.Unmarshal([]byte(actualResp), &actual)
	if err != nil {
		log.Println("error unmarshal actual response => ", err.Error())
		return
	}

	if expected != actual {
		t.Errorf("Expected Response %v. Got Response %v\n", expected, actual)
	}
}
