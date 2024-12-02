package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thejixer/memoir/internal/models"
)

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
