package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

// test creating a new company
func TestCreateCompany(t *testing.T) {
	droptable()
	app.db.Create(&User{Name: "a", Username: "a", Password: "a", Email: "a"})

	token, err := GenerateToken("a")
	if err != nil {
		log.Println("error while generating token => ", err.Error())
	}

	var tests = []UserTest{
		// first check if user logged in or not
		// so create request without a token first
		UserTest{
			req:  `{"name":"a", "description": "a"}`,
			resp: `{"company":null,"message":"user not logged in","user":null}`,
			code: 401,
		},
		// successfully create a company
		UserTest{
			req:   `{"name":"a", "description": "a"}`,
			resp:  `{"company":{"ID":1,"name":"a","description":"a","Admin":[{"ID":1,"name":"a","username":"a","password":"a","email":"a"}],"HR":null},"message":"successfully registered company","user":{"ID":1,"name":"a","username":"a","password":"a","email":"a"}}`,
			token: *token,
			code:  201,
		},
		// empty name
		UserTest{
			req:   `{"name":"", "description": "a"}`,
			resp:  `{"company":null,"message":"empty name","user":{"ID":1,"name":"a","username":"a","password":"a","email":"a"}}`,
			token: *token,
			code:  400,
		},
		// empty description
		UserTest{
			req:   `{"name":"a", "description": ""}`,
			resp:  `{"company":null,"message":"empty description","user":{"ID":1,"name":"a","username":"a","password":"a","email":"a"}}`,
			token: *token,
			code:  400,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/create-company", strings.NewReader(test.req))

		//write token to header
		req.Header["Authorization"] = []string{"Bearer " + test.token}
		log.Println(req.Header)

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

func TestUpdateCompany(t *testing.T) {
	droptable()
}
