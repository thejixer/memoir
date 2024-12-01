package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/thejixer/memoir/internal/models"
	"github.com/thejixer/memoir/pkg/encryption"
)

func (s *PostgresStore) createUserTable() error {

	query := `
	create table if not exists users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100) UNIQUE,
		isEmailVerified BOOLEAN,
		password VARCHAR,
		createdAt TIMESTAMPTZ
	)`

	_, err := s.db.Exec(query)

	return err
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(name, email, password string, isEmailVerified bool) (*models.User, error) {

	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		Name:            name,
		Email:           email,
		Password:        hashedPassword,
		IsEmailVerified: isEmailVerified,
		CreatedAt:       time.Now().UTC(),
	}

	query := `
	INSERT INTO USERS (name, email, isEmailVerified, password, createdAt)
	VALUES ($1, LOWER($2), $3, $4, $5) RETURNING id`
	lastInsertId := 0

	insertErr := r.db.QueryRow(
		query,
		newUser.Name,
		newUser.Email,
		newUser.IsEmailVerified,
		newUser.Password,
		newUser.CreatedAt,
	).Scan(&lastInsertId)

	if insertErr != nil {
		return nil, insertErr
	}

	newUser.ID = lastInsertId

	return newUser, nil
}

func (r *UserRepo) FindById(id int) (*models.User, error) {
	rows, err := r.db.Query("SELECT * FROM USERS WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanIntoUsers(rows)
	}

	return nil, errors.New("not found")
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	rows, err := r.db.Query("SELECT * FROM USERS WHERE email = LOWER($1)", email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanIntoUsers(rows)
	}

	return nil, errors.New("not found")
}

func ScanIntoUsers(rows *sql.Rows) (*models.User, error) {
	u := new(models.User)
	if err := rows.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.IsEmailVerified,
		&u.Password,
		&u.CreatedAt,
	); err != nil {
		return nil, err
	}
	return u, nil
}
