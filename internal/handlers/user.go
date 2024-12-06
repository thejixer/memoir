package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/memoir/internal/data-process-layer"
	"github.com/thejixer/memoir/internal/models"
	"github.com/thejixer/memoir/internal/utils"
	"github.com/thejixer/memoir/pkg/encryption"
)

type CustomContext struct {
	echo.Context
	User *models.User
}

func (h *HandlerService) HandleSignup(c echo.Context) error {

	body := models.SignUpDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	thisUser, _ := h.dbStore.UserRepo.FindByEmail(body.Email)

	if thisUser != nil {
		return WriteReponse(c, http.StatusBadRequest, "this email already exists in the database")
	}

	var err error
	thisUser, err = h.dbStore.UserRepo.Create(body.Name, body.Email, body.Password, false)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, err.Error())
	}

	verificationCode := CreateUUID()

	redisErr := h.redisStore.SetEmailVerificationCode(thisUser.Email, verificationCode)
	if redisErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	// ### to do ###
	// email the code to the user's email

	return WriteReponse(c, http.StatusOK, "please check your email to verify your email")
}

func (h *HandlerService) HandleRequestVerificationEmail(c echo.Context) error {
	body := models.RequestVerificationEmailDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "lack of data")
	}

	thisUser, err := h.dbStore.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "no user found")
	}

	if thisUser.IsEmailVerified {
		return WriteReponse(c, http.StatusBadRequest, "your email has already been verified")
	}

	verificationCode := CreateUUID()
	redisErr := h.redisStore.SetEmailVerificationCode(thisUser.Email, verificationCode)
	if redisErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	// ### to do ###
	// email the code to the user's email

	return WriteReponse(c, http.StatusOK, "please check your email to verify your email")
}

func (h *HandlerService) HandleEmailVerification(c echo.Context) error {
	email := c.QueryParam("email")
	verificationCode := c.QueryParam("code")

	if email == "" || verificationCode == "" {
		return WriteReponse(c, http.StatusBadRequest, "insufficient data")
	}

	thisUser, err := h.dbStore.UserRepo.FindByEmail(email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "user not found")
	}

	if val, err := h.redisStore.GetEmailVerificationCode(thisUser.Email); err != nil || val != verificationCode {
		return WriteReponse(c, http.StatusBadRequest, "code doesnt match")
	}

	updateErr := h.dbStore.UserRepo.VerifyEmail(email)
	h.redisStore.DeleteEmailVerificationCode(email)
	if updateErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	tokenString, err := utils.SignToken(thisUser.ID)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, models.TokenDTO{Token: tokenString})
}

func (h *HandlerService) HandleLogin(c echo.Context) error {
	body := models.LoginDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "lack of data")
	}

	thisUser, err := h.dbStore.UserRepo.FindByEmail(body.Email)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "bad requesst")
	}

	if !thisUser.IsEmailVerified {
		return WriteReponse(c, http.StatusUnauthorized, "your email is not verified")
	}

	if match := encryption.CheckPasswordHash(body.Password, thisUser.Password); !match {
		return WriteReponse(c, http.StatusBadRequest, "password doesnt match")
	}

	tokenString, err := utils.SignToken(thisUser.ID)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, models.TokenDTO{Token: tokenString})

}

func GetMe(c *echo.Context) (*models.User, error) {
	me := (*c).(CustomContext).User
	if me == nil {
		return nil, errors.New("unathorized")
	}
	return me, nil
}

func (h *HandlerService) HandleMe(c echo.Context) error {

	me, err := GetMe(&c)
	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	user := dataprocesslayer.ConvertToUserDto(me)

	return c.JSON(http.StatusOK, user)
}

func (h *HandlerService) HandleRequestChangePassword(c echo.Context) error {
	body := models.RequestChangePasswordDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "please provide a valid email")
	}

	thisUser, err := h.dbStore.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "no user found")
	}

	if !thisUser.IsEmailVerified {
		return WriteReponse(c, http.StatusBadRequest, "this option is for those who have validated their emails")
	}

	code := CreateUUID()

	redisErr := h.redisStore.SetPasswordChangeRequest(thisUser.Email, code)
	if redisErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	// ### to do
	// email the code

	return WriteReponse(c, http.StatusOK, "check your email")

}

func (h *HandlerService) HandleVerifyChangePasswordRequest(c echo.Context) error {
	email := c.QueryParam("email")
	code := c.QueryParam("code")

	if email == "" || code == "" {
		return WriteReponse(c, http.StatusBadRequest, "insufficient data")
	}

	thisUser, err := h.dbStore.UserRepo.FindByEmail(email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "user not found")
	}

	if val, err := h.redisStore.GetPasswordChangeRequest(thisUser.Email); err != nil || val != code {
		return WriteReponse(c, http.StatusBadRequest, "code doesnt match")
	}

	h.redisStore.DeletePasswordChangeRequest(thisUser.Email)
	redisErr := h.redisStore.CreatePasswordChangePermission(thisUser.Email, code)
	if redisErr != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	responsetext := fmt.Sprintf("you can change your password at %v/auth/change-password", os.Getenv("DOMAIN"))
	return WriteReponse(c, http.StatusOK, responsetext)
}

func (h *HandlerService) HandleChangePassword(c echo.Context) error {
	body := models.ChangePasswordDTO{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "please provide a valid email")
	}

	thisUser, err := h.dbStore.UserRepo.FindByEmail(body.Email)

	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "no user found")
	}

	if val, err := h.redisStore.GetPasswordChangePermission(thisUser.Email); err != nil || val != body.Code {
		return WriteReponse(c, http.StatusForbidden, "access denied")
	}

	if err := h.dbStore.UserRepo.UpdatePassword(body.Email, body.Password); err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	h.redisStore.DelPasswordChangePermission(body.Email)

	return WriteReponse(c, http.StatusOK, "password changed successfully")
}
