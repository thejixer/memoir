package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thejixer/memoir/internal/models"
)

func (h *HandlerService) HandleCreateMeeting(c echo.Context) error {
	me, err := GetMe(&c)

	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	body := models.CreateMeetingDto{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	thesePersons, err := h.dbStore.PersonRepo.GetPersonsByIds(body.PersonIds)
	if err != nil {
		fmt.Println("err here : ", err)
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	for _, p := range thesePersons {
		if p.UserId != me.ID {
			return WriteReponse(c, http.StatusForbidden, "forbidden")
		}
	}

	thisMeeting, err := h.dbStore.MeetingRepo.Create(body.Title, me.ID, body.PersonIds)
	if err != nil {
		fmt.Println("1")
		fmt.Println("err here : ", err)
		return WriteReponse(c, http.StatusInternalServerError, "this is on us, please try again")
	}

	u := models.MeetingDto{
		ID:        thisMeeting.ID,
		Title:     thisMeeting.Title,
		CreatedAt: thisMeeting.CreatedAt,
	}

	return c.JSON(http.StatusOK, u)
}
