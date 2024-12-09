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

	// 			meetingId INTEGER REFERENCES meeting (id),
	query := `
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100),
			content TEXT,
			type valid_note_types,
			personId INTEGER REFERENCES persons (id),
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
