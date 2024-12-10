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

func ConvertToLLNoteDto(a []*models.NoteDto, count int) models.LL_NoteDto {

	notes := make([]models.NoteDto, 0)

	for _, s := range a {
		x := models.NoteDto{
			ID:        s.ID,
			Title:     s.Title,
			Content:   s.Content,
			Tags:      s.Tags,
			CreatedAt: s.CreatedAt,
		}
		notes = append(notes, x)
	}

	return models.LL_NoteDto{
		Total:  count,
		Result: notes,
	}

}
