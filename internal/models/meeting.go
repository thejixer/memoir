package models

import "time"

type MeetingRepository interface {
	Create(title string, userId int, personsIds []int) (*Meeting, error)
	FindById(id int) (*Meeting, error)
}

type Meeting struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	UserId    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type MeetingDto struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}
