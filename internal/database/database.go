package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/thejixer/memoir/internal/models"
)

type PostgresStore struct {
	db         *sql.DB
	UserRepo   models.UserRepository
	PersonRepo models.PersonRepository
}

func NewPostgresStore() (*PostgresStore, error) {

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	conString := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=disable host=%v", dbUser, dbName, dbPassword, dbHost)
	db, err := sql.Open("postgres", conString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	userRepo := NewUserRepo(db)
	PersonRepo := NewPersonRepo(db)

	return &PostgresStore{
		db:         db,
		UserRepo:   userRepo,
		PersonRepo: PersonRepo,
	}, nil
}

func (s *PostgresStore) Init() error {

	if err := s.createUserTable(); err != nil {
		return err
	}

	if err := s.createPersonTable(); err != nil {
		return err
	}

	return nil

}
