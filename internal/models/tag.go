package models

type TagRepository interface {
	Create(title string, isForNote, isForMeeting bool, userId int) (*Tag, error)
	QueryNoteTags(text string, userId, page, limit int) ([]*Tag, int, error)
	QueryMeetingTags(text string, userId, page, limit int) ([]*Tag, int, error)
}

type Tag struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	IsForNote    bool   `json:"isForNote"`
	IsForMeeting bool   `json:"isForMeeting"`
	UserId       int    `json:"userId"`
}

type LL_Tag struct {
	Total  int   `json:"total"`
	Result []Tag `json:"result"`
}
