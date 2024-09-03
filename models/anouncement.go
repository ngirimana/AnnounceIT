package models

import (
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
