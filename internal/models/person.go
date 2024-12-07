package models

import "time"

type PersonRepository interface {
	Create(name, avatar string, userId int) (*Person, error)
	QueryMyPersons(text string, userId, page, limit int) ([]*Person, int, error)
	FindById(id int) (*Person, error)
}

type Person struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	UserId    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type PersonDto struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
}

type LL_PersonDto struct {
	Total  int         `json:"total"`
	Result []PersonDto `json:"result"`
}
