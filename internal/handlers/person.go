package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/memoir/internal/data-process-layer"
	"github.com/thejixer/memoir/internal/models"
)

func (h *HandlerService) HandleCreatePerson(c echo.Context) error {

	me, err := GetMe(&c)

	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	body := models.CreatePersonDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	newPerson, err := h.dbStore.PersonRepo.Create(body.Name, body.Avatar, me.ID)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToPersonDto(newPerson))

}

func (h *HandlerService) HandleQueryMyPersons(c echo.Context) error {
	me, err := GetMe(&c)

	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	text := c.QueryParam("text")
	p := c.QueryParam("page")
	l := c.QueryParam("limit")

	var page int
	var limit int

	page, err = strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	limit, err = strconv.Atoi(l)
	if err != nil {
		limit = 10
	}

	persons, count, err := h.dbStore.PersonRepo.QueryMyPersons(text, me.ID, page, limit)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	result := dataprocesslayer.ConvertToLLPersonDto(persons, count)
	return c.JSON(http.StatusOK, result)

}

func (h *HandlerService) HandleGetSinglePerson(c echo.Context) error {

	me, err := GetMe(&c)

	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	id := c.Param("id")
	personId, err := strconv.Atoi(id)
	if err != nil {
		return WriteReponse(c, http.StatusNotFound, "not found")
	}
	person, err := h.dbStore.PersonRepo.FindById(personId)
	if err != nil {
		msg := err.Error()
		if msg == "not found" {
			return WriteReponse(c, http.StatusNotFound, "not found")
		}
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	if person.UserId != me.ID {
		return WriteReponse(c, http.StatusForbidden, "access denied")
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToPersonDto(person))
}
