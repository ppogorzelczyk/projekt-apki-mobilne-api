package listitems

import (
	"buymeagiftapi/internal/domain"
	"database/sql"
)

type listItemRepository interface {
	Create(listItem domain.ListItem, userId int) (*domain.ListItem, error)
	Assign(listId, itemId, userId int, amount float64) error
	Delete(listId, itemId, userId int) error
	Update(listId, itemId int, listItem domain.ListItem, userId int) error
}

type service struct {
	repository listItemRepository
}

func NewService(db *sql.DB) *service {
	repo := NewRepository(db)
	return &service{repository: repo}
}

func (s *service) Create(item domain.ListItem, userId int) (*domain.ListItem, error) {
	listItem, err := s.repository.Create(item, userId)
	if err != nil {
		return nil, err
	}

	return listItem, nil
}

func (s *service) Assign(listId, itemId, userId int, amount float64) error {
	err := s.repository.Assign(listId, itemId, userId, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(listId, itemId, userId int) error {
	err := s.repository.Delete(listId, itemId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Update(listId, itemId int, listItem domain.ListItem, userId int) error {
	err := s.repository.Update(listId, itemId, listItem, userId)
	if err != nil {
		return err
	}

	return nil
}
