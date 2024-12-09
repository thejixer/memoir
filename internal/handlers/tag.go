package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/memoir/internal/data-process-layer"
	"github.com/thejixer/memoir/internal/models"
)

func (h *HandlerService) HandleCreateTag(c echo.Context) error {
	me, err := GetMe(&c)

	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	body := models.CreateTagDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	if !body.IsForMeeting && !body.IsForNote {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	newTag, err := h.dbStore.TagRepo.Create(body.Title, body.IsForNote, body.IsForMeeting, me.ID)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	return c.JSON(http.StatusOK, newTag)

}

func (h *HandlerService) HandleQueryNoteTags(c echo.Context) error {
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

	tags, count, err := h.dbStore.TagRepo.QueryNoteTags(text, me.ID, page, limit)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToLLTagDto(tags, count))

}

func (h *HandlerService) HandleQueryMeetingTags(c echo.Context) error {
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

	tags, count, err := h.dbStore.TagRepo.QueryMeetingTags(text, me.ID, page, limit)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToLLTagDto(tags, count))

}
