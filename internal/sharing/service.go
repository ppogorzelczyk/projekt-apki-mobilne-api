package sharing

import (
	"buymeagiftapi/internal/domain"
	"database/sql"
)

type sharingRepository interface {
	Share(listId, userId int, email string) (ShareResult, error)
	GetSharedLists(userId int) ([]SharedList, error)
	GetUsersThatLtstIsSharedWith(listId, userId int) ([]domain.User, error)
	Unshare(listId, ownerId, userId int) (UnshareResult, error)
}

type service struct {
	repository sharingRepository
}

func NewService(db *sql.DB) *service {
	repo := NewRepository(db)
	return &service{repository: repo}
}

func (s *service) Share(listId, userId int, email string) (ShareResult, error) {
	res, err := s.repository.Share(listId, userId, email)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *service) GetSharedLists(userId int) ([]SharedList, error) {
	lists, err := s.repository.GetSharedLists(userId)
	if err != nil {
		return nil, err
	}

	if lists == nil {
		return []SharedList{}, nil
	}

	return lists, nil
}

func (s *service) GetUsersThatListIsShareWith(listId, userId int) ([]domain.User, error) {
	users, err := s.repository.GetUsersThatLtstIsSharedWith(listId, userId)
	if err != nil {
		return nil, err
	}

	if users == nil {
		return []domain.User{}, nil
	}

	return users, nil
}

func (s *service) Unshare(listId, ownerId, userId int) (UnshareResult, error) {
	return s.repository.Unshare(listId, ownerId, userId)
}
