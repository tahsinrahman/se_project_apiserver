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
			Req:  `{"name":"a", "description": "a"}`,
			Resp: `{"company":null,"message":"user not logged in","user":null}`,
			Code: 401,
		},
		// successfully create a company
		UserTest{
			Req:   `{"name":"a", "description": "a"}`,
			Resp:  `{"company":{"HR":null,"ID":1,"admin":[{"ID":1,"email":"a","name":"a","username":"a"}],"description":"a","name":"a"},"message":"successfully registered company","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  201,
		},
		// empty name
		UserTest{
			Req:   `{"name":"", "description": "a"}`,
			Resp:  `{"company":null,"message":"empty name","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  400,
		},
		// empty description
		UserTest{
			Req:   `{"name":"a", "description": ""}`,
			Resp:  `{"company":null,"message":"empty description","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  400,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/company", strings.NewReader(test.Req))

		if err != nil {
			log.Println("error in creating req => ", err.Error())
			continue
		}

		//write token to header
		req.Header["Authorization"] = []string{"Bearer " + test.Token}

		//debugReq(req)

		response := executeRequest(req)

		resp, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error in reading resp => ", err.Error())
			continue
		}

		log.Println("request => ", test.Req)
		log.Println("response => ", string(resp))

		checkResponseCode(t, test.Code, response.Code, test.Resp, string(resp[:len(resp)-1]))
	}
}

func TestUpdateCompany(t *testing.T) {
	//droptable()

	//app.db.Create(&User{Name: "a", Username: "a", Password: "a", Email: "a"})
	app.db.Create(&User{Name: "b", Username: "b", Password: "b", Email: "b"})
	//app.db.Create(&Company{Name: "a", Description: "a"})

	token, err := GenerateToken("a")
	if err != nil {
		log.Println("error while generating token => ", err.Error())
	}
	token2, err := GenerateToken("b")
	if err != nil {
		log.Println("error while generating token => ", err.Error())
	}

	var tests = []CompanyTest{
		// first check if user logged in or not
		// so create request without a token first
		CompanyTest{
			Req:  `{"name":"a", "description": "a"}`,
			Resp: `{"company":null,"message":"user not logged in","user":null}`,
			Code: 401,
			ID:   "1",
		},
		// invalid id
		CompanyTest{
			Req:   `{"name":"a", "description": "a", "admin":[]}`,
			Resp:  `{"company":null,"message":"record not found","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Code:  400,
			Token: *token,
			ID:    "10",
		},

		// current user is not admin
		CompanyTest{
			Req:   `{"name":"a", "description": "a", "admin":["ID":1,"email":"a","name":"a","username":"a"]}`,
			Resp:  `{"company":null,"message":"user is not admin","user":{"ID":2,"email":"b","name":"b","username":"b"}}`,
			Code:  401,
			Token: *token2,
			ID:    "1",
		},

		// empty name
		CompanyTest{
			Req:   `{"name":"", "description": "a"}`,
			Resp:  `{"company":null,"message":"empty name","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  400,
			ID:    "1",
		},
		// empty description
		CompanyTest{
			Req:   `{"name":"a", "description": ""}`,
			Resp:  `{"company":null,"message":"empty description","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  400,
			ID:    "1",
		},

		// empty admin
		CompanyTest{
			Req:   `{"name":"a", "description": "a"}`,
			Resp:  `{"company":null,"message":"empty admin list","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  400,
			ID:    "1",
		},

		// invalid admin
		CompanyTest{
			Req:   `{"name":"a", "description": "a", "admin":[{ "ID":1,"email":"c","name":"a","username":"a"}]}`,
			Resp:  `{"company":null,"message":"admin not found =\u003e record not found","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  400,
			ID:    "1",
		},

		// invalid HR
		CompanyTest{
			Req:   `{"name":"a", "description": "a", "admin":[{ "ID":1,"email":"a","name":"a","username":"a"}], "HR":[{ "ID":1,"email":"c","name":"a","username":"a"}]}`,
			Resp:  `{"company":null,"message":"hr account not found =\u003e record not found","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  400,
			ID:    "1",
		},

		// ok
		CompanyTest{
			Req:   `{"name":"b", "description": "b", "admin":[{ "ID":2,"email":"b","name":"b","username":"b"}], "HR":[{ "ID":1,"email":"a","name":"a","username":"a"}]}`,
			Resp:  `{"company":{"HR":[{"ID":1,"email":"a","name":"a","username":"a"}],"ID":1,"admin":[{"ID":2,"email":"b","name":"b","username":"b"}],"description":"b","name":"b"},"message":"successfully updated","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  200,
			ID:    "1",
		},

		// prev admin
		CompanyTest{
			Req:   `{"name":"a", "description": "a", "admin":[{ "ID":2,"email":"b","name":"b","username":"b"}], "HR":[{ "ID":1,"email":"a","name":"a","username":"a"}]}`,
			Resp:  `{"company":null,"message":"user is not admin","user":{"ID":1,"email":"a","name":"a","username":"a"}}`,
			Token: *token,
			Code:  401,
			ID:    "1",
		},

		// new admin
		CompanyTest{
			Req:   `{"name":"a", "description": "a", "admin":[{ "ID":1,"email":"a","name":"a","username":"a"}], "HR":[{ "ID":1,"email":"a","name":"a","username":"a"}]}`,
			Resp:  `{"company":{"HR":null,"ID":1,"admin":[{"ID":1,"email":"a","name":"a","username":"a"}],"description":"a","name":"a"},"message":"successfully updated","user":{"ID":2,"email":"b","name":"b","username":"b"}}`,
			Token: *token2,
			Code:  200,
			ID:    "1",
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("PUT", "/"+test.ID+"/company", strings.NewReader(test.Req))

		if err != nil {
			log.Println("error in creating req => ", err.Error())
			continue
		}

		//write token to header
		req.Header["Authorization"] = []string{"Bearer " + test.Token}

		//debugReq(req)

		response := executeRequest(req)

		resp, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error in reading resp => ", err.Error())
			continue
		}

		log.Println("request => ", test.Req)
		log.Println("response => ", string(resp))

		checkResponseCode(t, test.Code, response.Code, test.Resp, string(resp[:len(resp)-1]))
	}
}
