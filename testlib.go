package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

type UserTest struct {
	Req   string
	Resp  string
	Code  int
	Token string
}

type CompanyTest struct {
	Req   string
	Resp  string
	Code  int
	Token string
	ID    string
}

func runTest(req *http.Request, test UserTest, t *testing.T) {
	response := executeRequest(req)

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("error in reading resp => ", err.Error())
		return
	}

	log.Println("request => ", test.Req)
	log.Println("response => ", string(resp))

	checkResponseCode(t, test.Code, response.Code, test.Resp, string(resp[:len(resp)-1]))
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
		t.Errorf("\nExp response %v\nGot response %v\n", expectedResp, actualResp)
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

func droptable() {
	app.db.DropTable(User{})
	app.db.DropTable(Company{})

	app.db.Exec("DROP TABLE company_admin;")
	app.db.Exec("DROP TABLE company_hr;")

	app.CreateTables()
}
