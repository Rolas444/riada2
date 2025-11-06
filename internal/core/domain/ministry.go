package domain

import "time"

// Ministry represents a ministry within the church.
type Ministry struct {
    ID          uint       `json:"id"`
    Name        string     `json:"name"`
    Description *string    `json:"description,omitempty"`
    Mission     *string    `json:"mission,omitempty"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}

// MinistryMember links a person to a ministry with an optional role.
type MinistryMember struct {
    ID         uint       `json:"id"`
    MinistryID uint       `json:"ministryID"`
    PersonID   uint       `json:"personID"`
    Role       *string    `json:"role,omitempty"`
    Status     string     `json:"status"`
    CreatedAt  time.Time  `json:"createdAt"`
    UpdatedAt  time.Time  `json:"updatedAt"`
}
