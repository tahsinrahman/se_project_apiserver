package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	app = App{}
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME_TEST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
	)

	code := m.Run()

	app.db.DropTable(User{})

	os.Exit(code)
}

type UserTest struct {
	user     User
	code     int
	response string
}

func TestSignUp(t *testing.T) {
	var users = []UserTest{
		// ok
		UserTest{
			user: User{
				Name:     "a",
				Username: "a",
				Email:    "a",
				Password: "a",
			},
			code:     201,
			response: `{"User":{"ID":1,"name":"a","username":"a","password":"a","email":"a"},"Message":"successfully registered"}`,
		},
		// same username
		UserTest{
			user: User{
				Name:     "a",
				Username: "a",
				Email:    "b",
				Password: "b",
			},
			code:     400,
			response: `{"User":null,"Message":"pq: duplicate key value violates unique constraint \"users_username_key\""}`,
		},
		// same email
		UserTest{
			user: User{
				Name:     "b",
				Username: "b",
				Email:    "a",
				Password: "b",
			},
			code:     400,
			response: `{"User":null,"Message":"pq: duplicate key value violates unique constraint \"users_email_key\""}`,
		},
		// empty name
		UserTest{
			user: User{
				Name:     "",
				Username: "c",
				Email:    "c",
				Password: "c",
			},
			code:     400,
			response: `{"User":null,"Message":"empty name"}`,
		},
		// empty username
		UserTest{
			user: User{
				Name:     "a",
				Username: "",
				Email:    "d",
				Password: "d",
			},
			code:     400,
			response: `{"User":null,"Message":"empty username"}`,
		},
		// empty password
		UserTest{
			user: User{
				Name:     "a",
				Username: "e",
				Email:    "e",
				Password: "",
			},
			code:     400,
			response: `{"User":null,"Message":"empty password"}`,
		},
		// empty email
		UserTest{
			user: User{
				Name:     "f",
				Username: "f",
				Email:    "",
				Password: "f",
			},
			code:     400,
			response: `{"User":null,"Message":"empty email"}`,
		},
	}

	for _, test := range users {
		js, err := json.Marshal(test.user)
		if err != nil {
			continue
		}
		payload := fmt.Sprintf(string(js[:len(js)]))

		req, _ := http.NewRequest("POST", "/signup", strings.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, test.code, response.Code, test.response, debugResp(response), test.user)
	}

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	return resp
}

func checkResponseCode(t *testing.T, expectedCode int, actualCode int, expectedResp string, actualResp string, user User) {
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, actualCode)
		t.Errorf("%v\n", user)
	}
	if expectedResp != actualResp {
		t.Errorf("\nExpected response %v\nGot response %v\n", expectedResp, actualResp)
		t.Errorf("%v\n", user)
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

func debugResp(resp *httptest.ResponseRecorder) string {
	respDump, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error ", err)
		return ""
	} else {
		return string(respDump[:len(respDump)-1])
	}
}
