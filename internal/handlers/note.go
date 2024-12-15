package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	dataprocesslayer "github.com/thejixer/memoir/internal/data-process-layer"
	"github.com/thejixer/memoir/internal/models"
)

func (h *HandlerService) HandleCreatePersonNote(c echo.Context) error {
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

	thisPerson, err := h.dbStore.PersonRepo.FindById(body.TargetId)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}
	if thisPerson.UserId != me.ID {
		return WriteReponse(c, http.StatusForbidden, "forbidden")
	}

	note, err := h.dbStore.NoteRepo.CreatePersonNote(body.Title, body.Content, body.TargetId, me.ID, body.TagIds)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	tags := dataprocesslayer.ConvertToTagDtoArray(theseTags)

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToNoteDto(note, tags))

}

func (h *HandlerService) HandleGetNotesByPersonId(c echo.Context) error {
	me, err := GetMe(&c)

	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	i := c.Param("personId")
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
	personId, err := strconv.Atoi(i)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	thisPerson, err := h.dbStore.PersonRepo.FindById(personId)
	if err != nil {
		msg := err.Error()
		if msg == "not found" {
			return WriteReponse(c, http.StatusNotFound, "not found")
		}
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")

	}
	if thisPerson.UserId != me.ID {
		return WriteReponse(c, http.StatusForbidden, "forbidden")
	}

	theseNotes, count, err := h.dbStore.NoteRepo.GetNotesByPersonId(personId, me.ID, page, limit)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	var wg sync.WaitGroup
	ch := make(chan []models.TagDto, len(theseNotes))
	for i := range theseNotes {
		wg.Add(1)
		go func(noteID int, idx int) {
			defer wg.Done()
			h.dbStore.TagRepo.FetchTagsForNote(noteID, ch)
		}(theseNotes[i].ID, i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := 0; i < len(theseNotes); i++ {
		tags := <-ch
		theseNotes[i].Tags = tags
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToLLNoteDto(theseNotes, count))

}

func (h *HandlerService) HandleCreateMeetingNote(c echo.Context) error {
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

	if len(body.TagIds) == 0 {
		return WriteReponse(c, http.StatusBadRequest, "you need to provide tags")
	}

	theseTags, err := h.dbStore.TagRepo.GetTagsById(body.TagIds)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}
	for _, item := range theseTags {
		if item.UserId != me.ID || !item.IsForMeeting {
			return WriteReponse(c, http.StatusBadRequest, "bad input")
		}
	}

	thisMeeting, err := h.dbStore.MeetingRepo.FindById(body.TargetId)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}
	if thisMeeting.UserId != me.ID {
		return WriteReponse(c, http.StatusForbidden, "forbidden")
	}

	note, err := h.dbStore.NoteRepo.CreateMeetingNote(body.Title, body.Content, body.TargetId, me.ID, body.TagIds)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	tags := dataprocesslayer.ConvertToTagDtoArray(theseTags)

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToNoteDto(note, tags))

}

func (h *HandlerService) HandleGetNotesByMeetingId(c echo.Context) error {
	me, err := GetMe(&c)

	if err != nil {
		return WriteReponse(c, http.StatusUnauthorized, "unathorized")
	}

	i := c.Param("meetingId")
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
	meetingId, err := strconv.Atoi(i)
	if err != nil {
		return WriteReponse(c, http.StatusBadRequest, "bad input")
	}

	thisMeeting, err := h.dbStore.MeetingRepo.FindById(meetingId)
	if err != nil {
		msg := err.Error()
		if msg == "not found" {
			return WriteReponse(c, http.StatusNotFound, "not found")
		}
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")

	}
	if thisMeeting.UserId != me.ID {
		return WriteReponse(c, http.StatusForbidden, "forbidden")
	}

	theseNotes, count, err := h.dbStore.NoteRepo.GetNotesByMeetingId(meetingId, me.ID, page, limit)
	if err != nil {
		return WriteReponse(c, http.StatusInternalServerError, "this one's on us")
	}

	var wg sync.WaitGroup
	ch := make(chan []models.TagDto, len(theseNotes))
	for i := range theseNotes {
		wg.Add(1)
		go func(noteID int, idx int) {
			defer wg.Done()
			h.dbStore.TagRepo.FetchTagsForNote(noteID, ch)
		}(theseNotes[i].ID, i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := 0; i < len(theseNotes); i++ {
		tags := <-ch
		theseNotes[i].Tags = tags
	}

	return c.JSON(http.StatusOK, dataprocesslayer.ConvertToLLNoteDto(theseNotes, count))
}
