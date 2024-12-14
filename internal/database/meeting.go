package database

import (
	"database/sql"
	"time"

	"github.com/thejixer/memoir/internal/models"
)

func (s *PostgresStore) createMeetingTable() error {

	query := `
		CREATE TABLE IF NOT EXISTS meetings (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100),
			userId INTEGER REFERENCES users (id),
			createdAt TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) createAttendanceTable() error {

	query := `
		CREATE TABLE IF NOT EXISTS attendance (
			person_id INTEGER REFERENCES persons (id) ON DELETE CASCADE,
			meeting_id INTEGER REFERENCES meetings (id) ON DELETE CASCADE,
			PRIMARY KEY (meeting_id, person_id)
		);`

	_, err := s.db.Exec(query)

	return err
}

type MeetingRepo struct {
	db *sql.DB
}

func NewMeetingRepo(db *sql.DB) *MeetingRepo {
	return &MeetingRepo{
		db: db,
	}
}

func (r *MeetingRepo) Create(title string, userId int, personsIds []int) (*models.Meeting, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	newMeeting := &models.Meeting{
		ID:        0,
		Title:     title,
		UserId:    userId,
		CreatedAt: time.Now().UTC(),
	}

	var meetingId int
	err = tx.QueryRow(`
		INSERT INTO meetings (title, userId) 
		VALUES ($1, $2) RETURNING id`,
		title, userId,
	).Scan(&meetingId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newMeeting.ID = meetingId
	for _, personId := range personsIds {
		_, err := tx.Exec(`
			INSERT INTO attendance (meeting_id, person_id) 
			VALUES ($1, $2)`,
			meetingId, personId,
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

	return newMeeting, nil

}
