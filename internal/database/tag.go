package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/thejixer/memoir/internal/models"
)

func (s *PostgresStore) createTagTable() error {

	query := `
	create table if not exists tags (
		id SERIAL PRIMARY KEY,
		title VARCHAR(50),
		isForNote BOOLEAN,
		isForMeeting BOOLEAN,
 		userId integer REFERENCES users (id)
	)`

	_, err := s.db.Exec(query)

	return err
}

type TagRepo struct {
	db *sql.DB
}

func NewTagRepo(db *sql.DB) *TagRepo {
	return &TagRepo{
		db: db,
	}
}

func (r *TagRepo) Create(title string, isForNote, isForMeeting bool, userId int) (*models.Tag, error) {
	newTag := &models.Tag{
		Title:        title,
		IsForNote:    isForNote,
		IsForMeeting: isForMeeting,
		UserId:       userId,
	}

	query := `
	INSERT INTO TAGS (title, isForNote, isForMeeting, userId)
	VALUES ($1, $2, $3, $4) RETURNING id`

	lastInsertId := 0
	insertErr := r.db.QueryRow(
		query,
		newTag.Title,
		newTag.IsForNote,
		newTag.IsForMeeting,
		newTag.UserId,
	).Scan(&lastInsertId)

	if insertErr != nil {
		return nil, insertErr
	}

	newTag.ID = lastInsertId

	return newTag, nil
}

func (r *TagRepo) QueryNoteTags(text string, userId, page, limit int) ([]*models.Tag, int, error) {
	offset := page * limit
	query := `SELECT * FROM TAGS 
		WHERE LOWER(TAGS.title) LIKE $1 AND userId = $2 AND isForNote = true
		ORDER BY id
		OFFSET $3 ROWS
		FETCH NEXT $4 ROWS ONLY`
	str := "%" + strings.ToLower(text) + "%"
	rows, err := r.db.Query(query, str, userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	tags := []*models.Tag{}
	for rows.Next() {
		u, err := ScanIntoTags(rows)
		if err != nil {
			return nil, 0, err
		}
		tags = append(tags, u)
	}
	var count int
	r.db.QueryRow(`
		SELECT count(id) 
		FROM TAGS
		WHERE LOWER(TAGS.title) LIKE $1 AND userId = $2 AND isForNote = true
	`, str, userId).Scan(&count)

	return tags, count, nil
}

func (r *TagRepo) QueryMeetingTags(text string, userId, page, limit int) ([]*models.Tag, int, error) {
	offset := page * limit
	query := `SELECT * FROM TAGS 
		WHERE LOWER(TAGS.title) LIKE $1 AND userId = $2 AND isForMeeting = true
		ORDER BY id
		OFFSET $3 ROWS
		FETCH NEXT $4 ROWS ONLY`
	str := "%" + strings.ToLower(text) + "%"
	rows, err := r.db.Query(query, str, userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	tags := []*models.Tag{}
	for rows.Next() {
		u, err := ScanIntoTags(rows)
		if err != nil {
			return nil, 0, err
		}
		tags = append(tags, u)
	}
	var count int
	r.db.QueryRow(`
		SELECT count(id) 
		FROM TAGS
		WHERE LOWER(TAGS.title) LIKE $1 AND userId = $2 AND isForMeeting = true
	`, str, userId).Scan(&count)

	return tags, count, nil
}

func (r *TagRepo) GetTagsById(tagIds []int) ([]*models.Tag, error) {
	// Create a placeholder string for the IN clause
	placeholders := make([]string, len(tagIds))
	args := make([]interface{}, len(tagIds))

	for i, id := range tagIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT * FROM tags WHERE id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag, err := ScanIntoTags(rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil

}

func (r *TagRepo) FetchTagsForNote(noteID int, ch chan<- []models.TagDto) {
	query := `
		SELECT t.id, t.title 
		FROM tags t
		INNER JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = $1`
	rows, err := r.db.Query(query, noteID)
	if err != nil {
		ch <- nil
		return
	}
	defer rows.Close()
	var tags []models.TagDto
	for rows.Next() {
		tag := new(models.TagDto)
		if err := rows.Scan(&tag.ID, &tag.Title); err != nil {
			ch <- nil
			return
		}
		tags = append(tags, *tag)
	}
	ch <- tags
}

func ScanIntoTags(rows *sql.Rows) (*models.Tag, error) {
	u := new(models.Tag)
	if err := rows.Scan(
		&u.ID,
		&u.Title,
		&u.IsForNote,
		&u.IsForMeeting,
		&u.UserId,
	); err != nil {
		return nil, err
	}
	return u, nil
}
