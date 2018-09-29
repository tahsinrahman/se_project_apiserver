package main

type UserJobStatus struct {
	ID     uint `gorm:"primary_key"`
	UserID uint
	JobID  uint
	Status string
}
