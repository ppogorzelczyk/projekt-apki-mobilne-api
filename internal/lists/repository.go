package lists

import (
	"buymeagiftapi/internal/domain"
	"database/sql"
	"errors"
	"log/slog"
	"sort"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetUserLists(userId int) ([]domain.List, error) {
	rows, err := r.db.Query("SELECT id, owner_id, title, description, event_date, created_at, updated_at FROM lists WHERE owner_id = ?", userId)
	if err != nil {
		slog.Error("Error getting user lists: %s", err)
		return nil, err
	}

	defer rows.Close()
	var lists []domain.List
	for rows.Next() {
		var list domain.List
		err := rows.Scan(
			&list.Id,
			&list.OwnerId,
			&list.Title,
			&list.Description,
			&list.EventDate,
			&list.CreatedAt,
			&list.UpdatedAt,
		)
		if err != nil {
			slog.Error("Error scanning user lists: %s", err)
			return nil, err
		}
		lists = append(lists, list)
	}

	return lists, nil
}

func (r *repository) Create(list domain.List) (domain.List, error) {
	now := time.Now()
	list.CreatedAt = now
	list.UpdatedAt = now

	var tx *sql.Tx
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		slog.Error("Error beginning transaction: %s", err)
		return domain.List{}, err
	}

	result, err := tx.Exec("INSERT INTO lists (owner_id, title, description, event_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?) RETURNING id",
		list.OwnerId,
		list.Title,
		list.Description,
		list.EventDate,
		list.CreatedAt,
		list.UpdatedAt,
	)

	if err != nil {
		slog.Error("Error creating list: %s", err)
		tx.Rollback()
		return domain.List{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("Error getting last insert id: %s", err)
		tx.Rollback()
		return domain.List{}, err
	}

	list.Id = int(id)

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		slog.Error("Error committing transaction: %s", err)
		return domain.List{}, err
	}

	return list, nil
}

func (r *repository) GetListById(listId int, userId int) (*domain.List, error) {
	var list domain.List
	queryList := `
    SELECT id, owner_id, title, description, event_date, created_at, updated_at
    FROM lists
    WHERE id = ? AND (owner_id = ? OR EXISTS(SELECT * FROM list_assignments la WHERE la.user_id = ?) );
    `

	var tx *sql.Tx
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		slog.Error("Error beginning transaction: %s", err)
		return nil, err
	}

	err = tx.QueryRow(queryList, listId, userId, userId).Scan(
		&list.Id,
		&list.OwnerId,
		&list.Title,
		&list.Description,
		&list.EventDate,
		&list.CreatedAt,
		&list.UpdatedAt,
	)

	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("List not found")
			return nil, nil
		}

		slog.Error("Error getting list by id: %s", err)
		return nil, err
	}

	queryListItems := `
    SELECT li.id, li.list_id, li.title, li.description, li.price, li.link, li.photo, li.created_at, li.updated_at, u.id, u.username, u.email, u.phone, u.first_name, u.last_name, u.created_at, lia.amount
    FROM list_items li
    LEFT JOIN list_item_assignments lia ON li.id = lia.list_item_id
    LEFT JOIN users u ON lia.user_id = u.id
    WHERE li.list_id = ?;
    `

	rows, err := tx.Query(queryListItems, listId)
	if err != nil {
		tx.Rollback()
		slog.Error("Error getting list items: %s", err)
		return nil, err
	}

	defer rows.Close()

	listItemsMap := make(map[int]*domain.ListItem)
	for rows.Next() {
		var listItem domain.ListItem
		var user struct {
			Id        *int
			Username  *string
			Email     *string
			Phone     *string
			FirstName *string
			LastName  *string
			CreatedAt *time.Time
		}

		var amount *float64

		err := rows.Scan(
			&listItem.Id,
			&listItem.ListId,
			&listItem.Title,
			&listItem.Description,
			&listItem.Price,
			&listItem.Link,
			&listItem.Photo,
			&listItem.CreatedAt,
			&listItem.UpdatedAt,
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&amount,
		)
		if err != nil {
			tx.Rollback()
			slog.Error("Error scanning list items: %s", err)
			return nil, err
		}

		if _, ok := listItemsMap[listItem.Id]; !ok {
			listItemsMap[listItem.Id] = &listItem
		}

		if user.Id != nil {
			user := domain.Assignee{
				User: domain.User{
					Id:        *user.Id,
					Username:  *user.Username,
					Email:     *user.Email,
					Phone:     user.Phone,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					CreatedAt: *user.CreatedAt,
				},
				Amount: *amount,
			}
			listItemsMap[listItem.Id].Assignees = append(listItemsMap[listItem.Id].Assignees, user)
		}
	}

	listItems := make([]domain.ListItem, 0, len(listItemsMap))
	for _, listItem := range listItemsMap {
		listItems = append(listItems, *listItem)
	}

	//sort list by item id
	sort.Slice(listItems, func(i, j int) bool {
		return listItems[i].Id < listItems[j].Id
	})

	list.Items = listItems

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		slog.Error("Error committing transaction: %s", err)
		return nil, err
	}

	return &list, nil
}

func (r *repository) Delete(listId, userId int) error {
	var tx *sql.Tx
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		slog.Error("Error beginning transaction: %s", err)
		return err
	}

	_, err = tx.Exec("DELETE FROM lists WHERE id = ? AND owner_id = ?", listId, userId)
	if err != nil {
		slog.Error("Error deleting list: %s", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		slog.Error("Error committing transaction: %s", err)
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repository) Update(listId int, list domain.List, userId int) error {
	var tx *sql.Tx
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		slog.Error("Error beginning transaction: %s", err)
		return err
	}

	_, err = tx.Exec("UPDATE lists SET title = ?, description = ?, event_date = ?, updated_at = ? WHERE id = ? AND owner_id = ?",
		list.Title,
		list.Description,
		list.EventDate,
		time.Now(),
		listId,
		userId,
	)

	if err != nil {
		slog.Error("Error updating list: %s", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Error committing transaction: %s", err)
		tx.Rollback()
		return err
	}

	return nil
}
