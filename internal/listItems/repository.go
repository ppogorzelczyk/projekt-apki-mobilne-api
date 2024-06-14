package listitems

import (
	"buymeagiftapi/internal/domain"
	"database/sql"
	"log/slog"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(listItem domain.ListItem, userId int) (*domain.ListItem, error) {
	now := time.Now()
	listItem.CreatedAt = now
	listItem.UpdatedAt = now

	var tx *sql.Tx
	tx, err := r.db.Begin()

	if err != nil {
		slog.Error("Error beginning transaction: %s", err.Error())
		return nil, err
	}

	var listOwnerId int
	err = tx.QueryRow("SELECT owner_id FROM lists WHERE id = ?", listItem.ListId).Scan(&listOwnerId)
	if err != nil {
		slog.Error("Error getting list owner id: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	if listOwnerId != userId {
		slog.Error("User is not the owner of the list")
		tx.Rollback()
		return nil, nil
	}

	result, err := tx.Exec("INSERT INTO list_items (list_id, title, description, price, link, photo, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id",
		listItem.ListId,
		listItem.Title,
		listItem.Description,
		listItem.Price,
		listItem.Link,
		listItem.Photo,
		listItem.CreatedAt,
		listItem.UpdatedAt,
	)

	if err != nil {
		slog.Error("Error creating list item: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("Error getting last insert id: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	listItem.Id = int(id)

	err = tx.Commit()
	if err != nil {
		slog.Error("Error committing transaction: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	return &listItem, nil
}

func (r *repository) Assign(listId, itemId, userId int, amount float64) error {
	var tx *sql.Tx
	tx, err := r.db.Begin()

	if err != nil {
		slog.Error("Error beginning transaction: %s", err.Error())
		return err
	}

	var id int
	err = tx.QueryRow("SELECT user_id FROM list_assignments WHERE list_id = ? AND user_id = ?", listId, userId).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("User is not allowed to assign items to this list")
			tx.Rollback()
			return err
		}

		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO list_item_assignments ( list_item_id, user_id, amount) VALUES (?, ?, ?)",
		itemId,
		userId,
		amount,
	)

	if err != nil {
		slog.Error("Error assigning item to user: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		slog.Error("Error committing transaction: %s", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repository) Delete(listId, itemId, userId int) error {
	var tx *sql.Tx
	tx, err := r.db.Begin()

	if err != nil {
		slog.Error("Error beginning transaction: %s", err.Error())
		tx.Rollback()
		return err
	}

	var listOwnerId int
	err = tx.QueryRow("SELECT owner_id FROM lists WHERE id = ?", listId).Scan(&listOwnerId)
	if err != nil {
		slog.Error("Error getting list owner id: %s", err.Error())
		tx.Rollback()
		return err
	}

	if listOwnerId != userId {
		slog.Error("User is not the owner of the list")
		tx.Rollback()
		return nil
	}

	_, err = tx.Exec("DELETE FROM list_items WHERE id = ?", itemId)
	if err != nil {
		slog.Error("Error deleting list item: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		slog.Error("Error committing transaction: %s", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repository) Update(listId, itemId int, listItem domain.ListItem, userId int) error {
	var tx *sql.Tx
	tx, err := r.db.Begin()

	if err != nil {
		slog.Error("Error beginning transaction: %s", err.Error())
		tx.Rollback()
		return err
	}

	var listOwnerId int
	err = tx.QueryRow("SELECT owner_id FROM lists WHERE id = ?", listId).Scan(&listOwnerId)
	if err != nil {
		slog.Error("Error getting list owner id: %s", err.Error())
		tx.Rollback()
		return err
	}

	if listOwnerId != userId {
		slog.Error("User is not the owner of the list")
		tx.Rollback()
		return nil
	}

	_, err = tx.Exec("UPDATE list_items SET title = ?, description = ?, price = ?, link = ?, photo = ?, updated_at = ? WHERE id = ?",
		listItem.Title,
		listItem.Description,
		listItem.Price,
		listItem.Link,
		listItem.Photo,
		time.Now(),
		itemId,
	)

	if err != nil {
		slog.Error("Error updating list item: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Error committing transaction: %s", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
