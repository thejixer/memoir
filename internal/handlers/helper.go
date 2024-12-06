package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/thejixer/memoir/internal/models"
	"github.com/thejixer/memoir/internal/utils"
)

func WriteReponse(c echo.Context, s int, m string) error {
	return c.JSON(s, models.ResponseDTO{Msg: m, StatusCode: s})
}

func CreateUUID() string {

	env := os.Getenv("ENVIROMENT")
	if env == "DEV" || env == "TEST" {
		return "1111"
	}

	return uuid.New().String()
}

func FindSingleUser(h *HandlerService, id int) (*models.User, error) {

	var thisUser *models.User
	var err error
	thisUser = h.redisStore.GetUser(id)

	if thisUser != nil {
		return thisUser, nil
	}

	thisUser, err = h.dbStore.UserRepo.FindById(id)
	if err != nil {
		return nil, errors.New("not found")
	}

	h.redisStore.CacheUser(thisUser)

	return thisUser, nil
}

func generateMe(c *echo.Context, h *HandlerService) (*models.User, int, error) {
	tokenString, err := getToken(c)

	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("unathorized")
	}

	token, err := utils.VerifyToken(tokenString)

	if err != nil || !token.Valid {
		return nil, http.StatusUnauthorized, errors.New("unathorized")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["id"] == nil {
		return nil, http.StatusUnauthorized, errors.New("unathorized")
	}

	i := claims["id"].(string)
	userId, err := strconv.Atoi(i)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("this one's on us")
	}

	thisUser, err := FindSingleUser(h, userId)

	if err != nil || thisUser == nil {
		return nil, http.StatusUnauthorized, errors.New("unathorized")
	}

	if !thisUser.IsEmailVerified {
		return nil, http.StatusForbidden, errors.New("your email is not verified")
	}

	return thisUser, 0, nil

}

func getToken(c *echo.Context) (string, error) {
	req := (*c).Request()
	authSlice := req.Header["Auth"]

	if len(authSlice) == 0 {
		return "", fmt.Errorf("token does not exist")
	}

	s := strings.Split(authSlice[0], " ")

	if len(s) != 2 || s[0] != "ut" {
		return "", fmt.Errorf("bad token format")
	}

	return s[1], nil
}

func (h *HandlerService) Gaurd(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		me, code, err := generateMe(&c, h)
		if err != nil {
			return WriteReponse(c, code, err.Error())
		}

		cc := CustomContext{c, me}
		return next(cc)

	}

}
