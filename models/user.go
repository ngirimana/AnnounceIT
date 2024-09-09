package models

import (
	"errors"

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
	Flagged     bool   `json:"flagged"`
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

func (u *User) Authenticate() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("credentials invalid")
	}

	isPasswordValid := helpers.CheckPassword(u.Password, retrievedPassword)
	if !isPasswordValid {
		return errors.New("invalid credentials")
	}
	return nil

}

func GetUser(email string) (*User, error) {
	query := "SELECT id, first_name, last_name, email, phone_number, address, is_admin FROM users WHERE email = ?"
	var user User
	err := db.DB.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Address, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (u *User) FlagUser() error {
	query := "UPDATE users SET flagged = true WHERE id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.ID)
	return err
}
