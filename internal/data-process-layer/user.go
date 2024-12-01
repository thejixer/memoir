package dataprocesslayer

import "github.com/thejixer/memoir/internal/models"

func ConvertToUserDto(u *models.User) models.UserDto {

	return models.UserDto{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func ConvertToLLUserDto(a []*models.User, count int) models.LL_UserDto {

	users := make([]models.UserDto, 0)

	for _, s := range a {
		users = append(users, ConvertToUserDto(s))
	}

	return models.LL_UserDto{
		Total:  count,
		Result: users,
	}

}
