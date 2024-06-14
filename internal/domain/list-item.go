package domain

import "time"

type ListItem struct {
	Id          int       `json:"id"`
	ListId      int       `json:"listId"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	Link        *string   `json:"link"`
	Photo       *string   `json:"photo"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Assignees []Assignee `json:"assignees"`
}

type Assignee struct {
	User
	Amount float64 `json:"amount"`
}
