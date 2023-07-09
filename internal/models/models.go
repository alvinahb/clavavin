package models

import (
	"time"
)

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Wine is the wine model
type Wine struct {
	ID              int
	Name            string
	Domain          string
	Year            string
	AppellationType string
	AppellationName string
	Location        string
	Color           string
	Culture         string
	Varieties       string
	Robe            string
	Nose            string
	Taste           string
	Dishes          string
	Season          string
}

// Opinion is the opinion model
type Opinion struct {
	ID      int
	Wine    Wine
	Rate    int
	Comment string
}

// MailData holds an email message
type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}
