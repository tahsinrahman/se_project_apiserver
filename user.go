package main

import (
	"encoding/json"
)

// user table
type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
}

// company table
type Company struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	Admin       []User `json:"admin" gorm:"many2many:company_admin"`
	HR          []User `json:"HR" gorm:"many2many:company_hr"`
}

// response user without password field
func (user *User) Response() (resp map[string]interface{}) {
	js, _ := json.Marshal(user)
	json.Unmarshal(js, &resp)
	delete(resp, "password")
	return resp
}

func (company *Company) Response() (resp map[string]interface{}) {
	js, _ := json.Marshal(company)
	json.Unmarshal(js, &resp)

	if resp["admin"] != nil {
		for index, _ := range resp["admin"].([]interface{}) {
			delete(resp["admin"].([]interface{})[index].(map[string]interface{}), "password")
		}
	}
	if resp["HR"] != nil {
		for index, _ := range resp["HR"].([]interface{}) {
			delete(resp["HR"].([]interface{})[index].(map[string]interface{}), "password")
		}
	}

	return resp
}
