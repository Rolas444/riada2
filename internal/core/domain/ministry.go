package domain

import "time"

// Ministry represents a ministry within the church.
type Ministry struct {
	ID          uint
	Name        string
	Description *string
	Mission     *string
	Status      string // one-letter status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MinistryMember links a person to a ministry with an optional role.
type MinistryMember struct {
	ID         uint
	MinistryID uint
	PersonID   uint
	Role       *string
	Status     string // one-letter status
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
