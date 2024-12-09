package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/memoir/internal/data-process-layer"
	"github.com/thejixer/memoir/internal/models"
)

func (h *HandlerService) HandleCreateNote(c echo.Context) error {
	me, err := GetMe(&c)

	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	body := models.CreateNoteDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	theseTags, err := h.dbStore.TagRepo.GetTagsById(body.TagIds)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}
	for _, item := range theseTags {
		if item.UserId != me.ID || !item.IsForNote {
			return WriteReponse(c, http.StatusBadRequest, "bad input")
		}
	}

	note, err := h.dbStore.NoteRepo.CreatePersonNote(body.Title, body.Content, body.TargetId, me.ID, body.TagIds)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	tags := dataprocesslayer.ConvertToTagDtoArray(theseTags)

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToNoteDto(note, tags))

}
