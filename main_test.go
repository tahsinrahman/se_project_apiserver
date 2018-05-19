package main

import (
	"encoding/json"
	"fmt"
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
	user User
	code int
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
			code: 201,
		},
		// same username
		UserTest{
			user: User{
				Name:     "a",
				Username: "a",
				Email:    "b",
				Password: "b",
			},
			code: 400,
		},
		// same email
		UserTest{
			user: User{
				Name:     "a",
				Username: "b",
				Email:    "a",
				Password: "b",
			},
			code: 400,
		},
		// empty name
		UserTest{
			user: User{
				Name:     "",
				Username: "c",
				Email:    "c",
				Password: "c",
			},
			code: 400,
		},
		// empty username
		UserTest{
			user: User{
				Name:     "a",
				Username: "",
				Email:    "d",
				Password: "d",
			},
			code: 400,
		},
		// empty password
		UserTest{
			user: User{
				Name:     "a",
				Username: "e",
				Email:    "e",
				Password: "",
			},
			code: 400,
		},
		// empty email
		UserTest{
			user: User{
				Name:     "f",
				Username: "f",
				Email:    "",
				Password: "f",
			},
			code: 400,
		},
	}

	for _, test := range users {
		//log.Println("/////////start///////////")

		js, err := json.Marshal(test.user)
		if err != nil {
			log.Println(err)
			continue
		}
		payload := fmt.Sprintf(string(js[:len(js)]))

		req, _ := http.NewRequest("POST", "/signup", strings.NewReader(payload))
		//debugReq(req)

		response := executeRequest(req)
		checkResponseCode(t, test.code, response.Code, test.user)
	}

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	return resp
}

func checkResponseCode(t *testing.T, expected, actual int, user User) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
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
