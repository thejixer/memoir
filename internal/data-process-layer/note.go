package dataprocesslayer

import "github.com/thejixer/memoir/internal/models"

func ConvertToNoteDto(u *models.Note, t []models.TagDto) models.NoteDto {

	return models.NoteDto{
		ID:        u.ID,
		Title:     u.Title,
		Content:   u.Content,
		Tags:      t,
		CreatedAt: u.CreatedAt,
	}
}
