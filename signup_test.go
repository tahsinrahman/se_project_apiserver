package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
)

type UserTest struct {
	req  string
	resp string
	code int
}

func TestSignUp(t *testing.T) {
	var users = []UserTest{
		// ok
		UserTest{
			req:  `{"name":"a", "username": "a", "email": "a", "password": "a"}`,
			resp: `{"message":"successfully registered","user":{"ID":1,"name":"a","username":"a","password":"a","email":"a"}}`,
			code: 201,
		},
		// same username
		UserTest{
			req:  `{"name":"a", "username": "a", "email": "b", "password": "b"}`,
			resp: `{"message":"pq: duplicate key value violates unique constraint \"users_username_key\"","user":null}`,
			code: 400,
		},
		// same email
		UserTest{
			req:  `{"name":"b", "username": "b", "email": "a", "password": "b"}`,
			resp: `{"message":"pq: duplicate key value violates unique constraint \"users_email_key\"","user":null}`,
			code: 400,
		},
		// empty name
		UserTest{
			req:  `{"name":"", "username": "c", "email": "c", "password": "c"}`,
			resp: `{"message":"empty name","user":null}`,
			code: 400,
		},
		// empty username
		UserTest{
			req:  `{"name":"c", "username": "", "email": "c", "password": "c"}`,
			resp: `{"message":"empty username","user":null}`,
			code: 400,
		},
		// empty email
		UserTest{
			req:  `{"name":"c", "username": "c", "email": "", "password": "c"}`,
			resp: `{"message":"empty email","user":null}`,
			code: 400,
		},
		// empty password
		UserTest{
			req:  `{"name":"c", "username": "c", "email": "c", "password": ""}`,
			resp: `{"message":"empty password","user":null}`,
			code: 400,
		},
	}

	for _, test := range users {
		req, err := http.NewRequest("POST", "/signup", strings.NewReader(test.req))
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

		log.Println("request => ", test.req)
		log.Println("response => ", string(resp))

		checkResponseCode(t, test.code, response.Code, test.resp, string(resp[:len(resp)-1]))
	}

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	return resp
}

func checkResponseCode(t *testing.T, expectedCode int, actualCode int, expectedResp string, actualResp string) {
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, actualCode)
	}
	if expectedResp != actualResp {
		t.Errorf("\nExpected response %v\nGot response %v\n", expectedResp, actualResp)
	}
}

func debugReq(req *http.Request) {
	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println("error ", err)
	} else {
		fmt.Println(string(reqDump))
	}
}

func TestSignIn(t *testing.T) {
	var tests []UserTest = []UserTest{
		// ok
		UserTest{
			req:  `{"username": "a", "password": "a"}`,
			resp: `{"message":"successfully logged in", "user":{"ID":1,"name":"a","username":"a","password":"a","email":"a"}}`,
			code: 200,
		},
		// empty username
		UserTest{
			req:  `{"username": "", "password": "a"}`,
			resp: `{"message":"empty username","user":null}`,
			code: 400,
		},
		// empty password
		UserTest{
			req:  `{"username": "a", "password": ""}`,
			resp: `{"message":"empty password","user":null}`,
			code: 400,
		},
		// invalid user
		UserTest{
			req:  `{"username": "b", "password": "b"}`,
			resp: `{"message":"username not found","user":null}`,
			code: 400,
		},
		// password mismatch
		UserTest{
			req:  `{"username": "a", "password": "b"}`,
			resp: `{"message":"invalid password","user":null}`,
			code: 400,
		},
	}
	for _, test := range tests {
		req, err := http.NewRequest("POST", "/signin", strings.NewReader(test.req))
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

		log.Println("request => ", test.req)
		log.Println("response => ", string(resp))

		checkSignIn(t, test.code, response.Code, test.resp, string(resp[:len(resp)-1]))
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
