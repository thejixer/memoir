package dataprocesslayer

import "github.com/thejixer/memoir/internal/models"

func ConvertToTagDto(t *models.Tag) models.TagDto {
	return models.TagDto{
		ID:    t.ID,
		Title: t.Title,
	}
}

func ConvertToTagDtoArray(a []*models.Tag) []models.TagDto {
	tags := make([]models.TagDto, 0)
	for _, s := range a {
		tags = append(tags, ConvertToTagDto(s))
	}
	return tags
}

func ConvertToLLTagDto(a []*models.Tag, count int) models.LL_TagDto {

	tags := make([]models.TagDto, 0)

	for _, s := range a {
		tags = append(tags, ConvertToTagDto(s))
	}

	return models.LL_TagDto{
		Total:  count,
		Result: tags,
	}

}
