package main

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
	Admin       []User `gorm:"many2many:company_admin"`
	HR          []User `gorm:"many2many:company_hr"`
}
