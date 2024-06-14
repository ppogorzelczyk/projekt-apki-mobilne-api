package lists

import (
	"buymeagiftapi/internal/domain"
	"database/sql"
)

type listRepository interface {
	GetUserLists(userId int) ([]domain.List, error)
	Create(list domain.List) (domain.List, error)
	GetListById(listId int, userId int) (*domain.List, error)
	Delete(listId int, userId int) error
	Update(listId int, list domain.List, userId int) error
}

type service struct {
	repository listRepository
}

func NewService(db *sql.DB) *service {
	repo := NewRepository(db)
	return &service{repository: repo}
}

func (s *service) GetUserLists(userId int) ([]domain.List, error) {
	lists, err := s.repository.GetUserLists(userId)
	if err != nil {
		return nil, err
	}

	if lists == nil {
		return []domain.List{}, nil
	}

	return lists, nil
}

func (s *service) CreateList(list domain.List) (domain.List, error) {
	list, err := s.repository.Create(list)
	if err != nil {
		return domain.List{}, err
	}

	return list, nil
}

func (s *service) GetListById(listId int, userId int) (*domain.List, error) {
	list, err := s.repository.GetListById(listId, userId)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *service) DeleteList(listId int, userId int) error {
	err := s.repository.Delete(listId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateList(listId int, list domain.List, userId int) error {
	err := s.repository.Update(listId, list, userId)
	if err != nil {
		return err
	}

	return nil
}
