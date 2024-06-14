package users

import (
	"buymeagiftapi/internal/domain"
	"database/sql"
)

type userRepository interface {
	GetUserByEmail(email string) (*domain.User, error)
	Create(user domain.User) (*domain.User, error)
	DeleteUser(userId int) error
}

type service struct {
	repository userRepository
}

func NewService(db *sql.DB) *service {
	repo := NewRepository(db)
	return &service{repository: repo}
}

func (s *service) GetUserByEmail(email string) (*domain.User, error) {
	return s.repository.GetUserByEmail(email)
}

func (s *service) CreateUser(user domain.User) (*domain.User, error) {
	return s.repository.Create(user)
}

func (s *service) DeleteUser(userId int) error {
	return s.repository.DeleteUser(userId)
}
