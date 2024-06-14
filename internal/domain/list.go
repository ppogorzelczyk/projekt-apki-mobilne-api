package domain

import "time"

type List struct {
	Id          int        `json:"id"`
	OwnerId     int        `json:"ownerId"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	EventDate   *time.Time `json:"eventDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`

	Items []ListItem `json:"items"`
}
