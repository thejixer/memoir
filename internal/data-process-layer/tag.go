package dataprocesslayer

import "github.com/thejixer/memoir/internal/models"

func ConvertToLLTagDto(a []*models.Tag, count int) models.LL_Tag {

	tags := make([]models.Tag, 0)

	for _, s := range a {

		tag := models.Tag{
			ID:           s.ID,
			Title:        s.Title,
			IsForNote:    s.IsForNote,
			IsForMeeting: s.IsForMeeting,
			UserId:       s.UserId,
		}

		tags = append(tags, tag)
	}

	return models.LL_Tag{
		Total:  count,
		Result: tags,
	}

}
