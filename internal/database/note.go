package database

import (
	"database/sql"
	"time"

	"github.com/thejixer/memoir/internal/models"
)

func (s *PostgresStore) createNoteTagTable() error {

	query := `
		CREATE TABLE IF NOT EXISTS note_tags (
			note_id INTEGER REFERENCES notes (id) ON DELETE CASCADE,
			tag_id INTEGER REFERENCES tags (id) ON DELETE CASCADE,
			PRIMARY KEY (note_id, tag_id)
		);`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) createNoteTable() error {

	query := `
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100),
			content TEXT,
			type valid_note_types,
			personId INTEGER REFERENCES persons (id),
			meetingId INTEGER REFERENCES meetings (id),
			userId INTEGER REFERENCES users (id),
			createdAt TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);`

	_, err := s.db.Exec(query)

	return err
}

type NoteRepo struct {
	db *sql.DB
}

func NewNoteRepo(db *sql.DB) *NoteRepo {
	return &NoteRepo{
		db: db,
	}
}

func (r *NoteRepo) CreatePersonNote(title, content string, personId, userId int, tagIds []int) (*models.Note, error) {
	// not completed
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	newNote := &models.Note{
		ID:        0,
		Title:     title,
		Content:   content,
		PersonId:  personId,
		MeetingId: 0,
		UserId:    userId,
		CreatedAt: time.Now().UTC(),
	}

	var noteID int
	err = tx.QueryRow(`
		INSERT INTO notes (title, content, type, personId, userId) 
		VALUES ($1, $2, 'person', $3, $4) RETURNING id`,
		title, content, personId, userId,
	).Scan(&noteID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newNote.ID = noteID
	for _, tagID := range tagIds {
		_, err := tx.Exec(`
			INSERT INTO note_tags (note_id, tag_id) 
			VALUES ($1, $2)`,
			noteID, tagID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return newNote, nil

}

func (r *NoteRepo) GetNotesByPersonId(persondId, userId, page, limit int) ([]*models.NoteDto, int, error) {
	offset := page * limit
	query := `
		SELECT id, title, content, createdAt FROM notes 
		WHERE personId = $1 AND userId = $2
		ORDER BY id desc
		OFFSET $3 ROWS
		FETCH NEXT $4 ROWS ONLY
	`
	rows, err := r.db.Query(query, persondId, userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var notes []*models.NoteDto
	for rows.Next() {
		note, err := ScanIntoNoteDto(rows)
		if err != nil {
			return nil, 0, err
		}
		notes = append(notes, note)
	}
	var count int
	r.db.QueryRow(`
		SELECT count(id) 
		FROM notes
		WHERE personId = $1 AND userId = $2
	`, persondId, userId).Scan(&count)

	return notes, count, nil
}

func (r *NoteRepo) CreateMeetingNote(title, content string, meetingId, userId int, tagIds []int) (*models.Note, error) {
	// not completed
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	newNote := &models.Note{
		ID:        0,
		Title:     title,
		Content:   content,
		PersonId:  0,
		MeetingId: meetingId,
		UserId:    userId,
		CreatedAt: time.Now().UTC(),
	}

	var noteID int
	err = tx.QueryRow(`
		INSERT INTO notes (title, content, type, meetingId, userId) 
		VALUES ($1, $2, 'meeting', $3, $4) RETURNING id`,
		title, content, meetingId, userId,
	).Scan(&noteID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newNote.ID = noteID
	for _, tagID := range tagIds {
		_, err := tx.Exec(`
			INSERT INTO note_tags (note_id, tag_id) 
			VALUES ($1, $2)`,
			noteID, tagID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return newNote, nil

}

func (r *NoteRepo) GetNotesByMeetingId(meetingId, userId, page, limit int) ([]*models.NoteDto, int, error) {
	offset := page * limit
	query := `
		SELECT id, title, content, createdAt FROM notes 
		WHERE meetingId = $1 AND userId = $2 AND type ='meeting'
		ORDER BY id desc
		OFFSET $3 ROWS
		FETCH NEXT $4 ROWS ONLY
	`
	rows, err := r.db.Query(query, meetingId, userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var notes []*models.NoteDto
	for rows.Next() {
		note, err := ScanIntoNoteDto(rows)
		if err != nil {
			return nil, 0, err
		}
		notes = append(notes, note)
	}
	var count int
	r.db.QueryRow(`
		SELECT count(id) 
		FROM notes
		WHERE meetingId = $1 AND userId = $2 AND type ='meeting'
	`, meetingId, userId).Scan(&count)

	return notes, count, nil
}

func ScanIntoNoteDto(rows *sql.Rows) (*models.NoteDto, error) {
	u := new(models.NoteDto)
	if err := rows.Scan(
		&u.ID,
		&u.Title,
		&u.Content,
		&u.CreatedAt,
	); err != nil {
		return nil, err
	}
	return u, nil
}
