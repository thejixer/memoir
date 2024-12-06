package dataprocesslayer

import "github.com/thejixer/memoir/internal/models"

func ConvertToPersonDto(u *models.Person) models.PersonDto {

	return models.PersonDto{
		ID:        u.ID,
		Name:      u.Name,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
	}
}

func ConvertToLLPersonDto(a []*models.Person, count int) models.LL_PersonDto {

	users := make([]models.PersonDto, 0)

	for _, s := range a {
		users = append(users, ConvertToPersonDto(s))
	}

	return models.LL_PersonDto{
		Total:  count,
		Result: users,
	}

}
