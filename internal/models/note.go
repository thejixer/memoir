package models

import "time"

type NoteRepository interface {
	CreatePersonNote(title, content string, personId, userId int, tagIds []int) (*Note, error)
	GetNotesByPersonId(persondId, userId, page, limit int) ([]*NoteDto, int, error)
}

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	PersonId  int       `json:"personId"`
	UserId    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type NoteDto struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []TagDto  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}

type LL_NoteDto struct {
	Total  int       `json:"total"`
	Result []NoteDto `json:"result"`
}
