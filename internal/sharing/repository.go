package sharing

import (
	"buymeagiftapi/internal/domain"
	"database/sql"
	"errors"
	"log/slog"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

type ShareResult int8
type UnshareResult int8

const (
	ShareSuccess ShareResult = iota
	ShareUserNotFound
	ShareAlreadyShared
	ShareError
)

const (
	UnshareSuccess UnshareResult = iota
	UnshareNotShared
	UnshareUserNotOwner
	UnshareError
)

type SharedList struct {
	domain.List
	Owner domain.User `json:"owner"`
}

func (r *repository) Share(listId, ownerId int, email string) (ShareResult, error) {
	var tx, err = r.db.Begin()
	if err != nil {
		slog.Error("Error beginning transaction: %s", err)
		return ShareError, err
	}

	var userId int
	queryUser := `
    SELECT id
    FROM users
    WHERE email = ?;
    `

	err = tx.QueryRow(queryUser, email).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			slog.Error("User not found")
			return ShareUserNotFound, nil
		}

		slog.Error("Error getting user by email: %s", err)
		return ShareError, err
	}

	var shared bool
	queryShared := `
    SELECT EXISTS (
        SELECT 1
        FROM list_assignments
        WHERE list_id = ? AND user_id = ?
    );
    `
	err = tx.QueryRow(queryShared, listId, userId).Scan(&shared)
	if err != nil {
		tx.Rollback()
		slog.Error("Error checking if list is already shared: %s", err)
		return ShareError, err
	}

	if shared {
		tx.Rollback()
		slog.Error("List is already shared")
		return ShareAlreadyShared, nil
	}

	var queryShare = `
    INSERT INTO list_assignments (list_id, user_id)
    VALUES (?, ?);
    `
	_, err = tx.Exec(queryShare, listId, userId)
	if err != nil {
		tx.Rollback()
		slog.Error("Error sharing list: %s", err)
		return ShareError, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		slog.Error("Error committing transaction: %s", err)
		return ShareError, err
	}

	return ShareSuccess, nil
}

func (r *repository) GetSharedLists(userId int) ([]SharedList, error) {
	var lists []SharedList
	queryLists := `
    SELECT l.id, l.owner_id, l.title, l.description, l.event_date, l.created_at, l.updated_at, u.id, u.username, u.email, u.phone, u.first_name, u.last_name, u.created_at
    FROM lists l
    JOIN list_assignments la ON l.id = la.list_id
    JOIN users u on u.id = l.owner_id
    WHERE la.user_id = ?;
    `

	var tx, err = r.db.Begin()
	if err != nil {
		tx.Rollback()
		slog.Error("Error beginning transaction: %s", err)
		return nil, err
	}

	rows, err := tx.Query(queryLists, userId)
	if err != nil {
		slog.Error("Error getting shared lists: %s", err)
		tx.Rollback()
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var list SharedList
		err = rows.Scan(
			&list.Id,
			&list.OwnerId,
			&list.Title,
			&list.Description,
			&list.EventDate,
			&list.CreatedAt,
			&list.UpdatedAt,
			&list.Owner.Id,
			&list.Owner.Username,
			&list.Owner.Email,
			&list.Owner.Phone,
			&list.Owner.FirstName,
			&list.Owner.LastName,
			&list.Owner.CreatedAt,
		)
		if err != nil {
			slog.Error("Error scanning shared list: %s", err)
			tx.Rollback()
			return nil, err
		}

		lists = append(lists, list)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		slog.Error("Error committing transaction: %s", err)
		return nil, err
	}

	return lists, nil
}

func (r *repository) GetUsersThatLtstIsSharedWith(listId, userId int) ([]domain.User, error) {
	var users []domain.User
	queryUsers := `
    SELECT u.id, u.username, u.email, u.phone, u.first_name, u.last_name, u.created_at 
    FROM list_assignments la 
    JOIN users u ON u.id = la.user_id
    JOIN lists l on l.id = la.list_id
    WHERE la.list_id = ? AND l.owner_id = ?;
    `

	rows, err := r.db.Query(queryUsers, listId, userId)
	if err != nil {
		slog.Error("Error getting users that list is shared with: %s", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var user domain.User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
		)
		if err != nil {
			slog.Error("Error scanning user that list is shared with: %s", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *repository) Unshare(listId, ownerId, userId int) (UnshareResult, error) {
	var tx, err = r.db.Begin()
	if err != nil {
		slog.Error("Error beginning transaction: %s", err)
		return UnshareError, err
	}

	var owner bool
	queryOwner := `
    SELECT EXISTS (
        SELECT 1
        FROM lists
        WHERE id = ? AND owner_id = ?
    );
    `
	err = tx.QueryRow(queryOwner, listId, ownerId).Scan(&owner)
	if err != nil {
		tx.Rollback()
		slog.Error("Error checking if list owner is unsharing: %s", err)
		return UnshareError, err
	}

	if !owner {
		tx.Rollback()
		slog.Error("List owner is not unsharing")
		return UnshareUserNotOwner, nil
	}

	var shared bool
	queryShared := `
    SELECT EXISTS (
        SELECT 1
        FROM list_assignments
        WHERE list_id = ? AND user_id = ?
    );
    `

	err = tx.QueryRow(queryShared, listId, userId).Scan(&shared)

	if err != nil {
		tx.Rollback()
		slog.Error("Error checking if list is shared: %s", err)
		return UnshareError, err
	}

	if !shared {
		tx.Rollback()
		slog.Error("List is not shared")
		return UnshareNotShared, nil
	}

	var queryUnshare = `
    DELETE FROM list_assignments
    WHERE list_id = ? AND user_id = ?;
    `

	_, err = tx.Exec(queryUnshare, listId, userId)
	if err != nil {
		tx.Rollback()
		slog.Error("Error unsharing list: %s", err)
		return UnshareError, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		slog.Error("Error committing transaction: %s", err)
		return UnshareError, err
	}

	return UnshareSuccess, nil
}
