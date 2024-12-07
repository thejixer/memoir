package database

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/thejixer/memoir/internal/models"
)

func (s *PostgresStore) createPersonTable() error {

	query := `
	create table if not exists persons (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		avatar VARCHAR,
 		userId integer REFERENCES users (id),
		createdAt TIMESTAMPTZ
	)`

	_, err := s.db.Exec(query)

	return err
}

type PersonRepo struct {
	db *sql.DB
}

func NewPersonRepo(db *sql.DB) *PersonRepo {
	return &PersonRepo{
		db: db,
	}
}

func (r *PersonRepo) Create(name, avatar string, userId int) (*models.Person, error) {

	newPerson := &models.Person{
		Name:      name,
		Avatar:    avatar,
		UserId:    userId,
		CreatedAt: time.Now().UTC(),
	}

	query := `
	INSERT INTO PERSONS (name, avatar, userId, createdAt)
	VALUES ($1, $2, $3, $4) RETURNING id`
	lastInsertId := 0

	insertErr := r.db.QueryRow(
		query,
		newPerson.Name,
		newPerson.Avatar,
		newPerson.UserId,
		newPerson.CreatedAt,
	).Scan(&lastInsertId)

	if insertErr != nil {
		return nil, insertErr
	}

	newPerson.ID = lastInsertId

	return newPerson, nil
}

func (r *PersonRepo) QueryMyPersons(text string, userId, page, limit int) ([]*models.Person, int, error) {
	offset := page * limit
	query := `SELECT * FROM persons 
		WHERE LOWER(PERSONS.name) LIKE $1 AND userId = $2
		ORDER BY id
		OFFSET $3 ROWS
		FETCH NEXT $4 ROWS ONLY`
	str := "%" + strings.ToLower(text) + "%"
	rows, err := r.db.Query(query, str, userId, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	persons := []*models.Person{}
	for rows.Next() {
		u, err := ScanIntoPersons(rows)
		if err != nil {
			return nil, 0, err
		}
		persons = append(persons, u)
	}
	var count int
	r.db.QueryRow(`
		SELECT count(id) 
		FROM PERSONS
		WHERE LOWER(PERSONS.name) LIKE $1 AND userId = $2
	`, str, userId).Scan(&count)

	return persons, count, nil

}

func (r *PersonRepo) FindById(id int) (*models.Person, error) {
	rows, err := r.db.Query("SELECT * FROM PERSONS WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanIntoPersons(rows)
	}

	return nil, errors.New("not found")
}

func ScanIntoPersons(rows *sql.Rows) (*models.Person, error) {
	u := new(models.Person)
	if err := rows.Scan(
		&u.ID,
		&u.Name,
		&u.Avatar,
		&u.UserId,
		&u.CreatedAt,
	); err != nil {
		return nil, err
	}
	return u, nil
}
