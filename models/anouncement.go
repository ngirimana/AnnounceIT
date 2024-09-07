package models

import (
	"fmt"
	"time"

	"github.com/ngirimana/AnnounceIT/db"
)

type Status int

// Define enum constants using iota
const (
	Pending     Status = iota // 0
	Accepted    Status = iota // 1
	Declined    Status = iota // 2
	Active      Status = iota // 3
	Deactivated Status = iota // 4
)

func (s Status) String() string {
	return [...]string{"Pending", "Accepted", "Declined", "Active", "Deactivated"}[s]
}

type Announcement struct {
	ID         int64     `json:"id"`
	OwnerID    int64     `json:"owner_id"`
	Status     Status    `json:"status"`
	Text       string    `json:"text"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	CreateDate time.Time `json:"create_date"`
}

func (a *Announcement) Create() error {
	query := `INSERT INTO announcements (owner_id, status, text, start_date, end_date, create_date) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	a.CreateDate = time.Now()
	a.Status = Pending
	result, err := stmt.Exec(a.OwnerID, a.Status, a.Text, a.StartDate, a.EndDate, a.CreateDate)
	if err != nil {
		return err
	}

	a.ID, err = result.LastInsertId()
	return err
}

func GetAnnouncements() ([]Announcement, error) {
	query := `SELECT * FROM announcements`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	announcements := []Announcement{}
	for rows.Next() {
		var a Announcement
		err := rows.Scan(&a.ID, &a.OwnerID, &a.Status, &a.Text, &a.StartDate, &a.EndDate, &a.CreateDate)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}

	return announcements, nil
}

func GetAnnouncementByID(id int64) (*Announcement, error) {
	query := `SELECT * FROM announcements WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var a Announcement
	err := row.Scan(&a.ID, &a.OwnerID, &a.Status, &a.Text, &a.StartDate, &a.EndDate, &a.CreateDate)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *Announcement) Update() error {
	query := `UPDATE announcements SET status = ?, text = ?, start_date = ?, end_date = ? WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	fmt.Println(a.Text)

	_, err = stmt.Exec(a.Status, a.Text, a.StartDate, a.EndDate, a.ID)
	return err
}

func (a *Announcement) Delete() error {
	query := `DELETE FROM announcements WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.ID)
	return err
}
