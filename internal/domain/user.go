package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	Phone        *string   `json:"phone"`
	FirstName    *string   `json:"firstName"`
	LastName     *string   `json:"lastName"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return err
	}

	u.PasswordHash = string(bytes)
	return nil
}

func (u *User) CheckPassword(providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(providedPassword))

	return err == nil
}
