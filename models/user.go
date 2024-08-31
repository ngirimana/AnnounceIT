package models

import (
	"github.com/ngirimana/AnnounceIT/db"
	"github.com/ngirimana/AnnounceIT/helpers"
)

type User struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	IsAdmin     bool   `json:"is_admin"`
}

func (u *User) Save() error {

	query := "INSERT INTO users (first_name, last_name, email, password, phone_number, address, is_admin) VALUES (?, ?, ?, ?, ?, ?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {

		return err
	}

	// Always close the statement after it is done being used
	defer stmt.Close()

	HashedPassword, err := helpers.HashPassword(u.Password)
	if err != nil {

		return err
	}

	newUser, err := stmt.Exec(u.FirstName, u.LastName, u.Email, HashedPassword, u.PhoneNumber, u.Address, u.IsAdmin)
	if err != nil {
		return err
	}

	u.ID, err = newUser.LastInsertId()
	return err
}
