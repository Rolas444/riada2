package domain

import (
	"time"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	ID           uint   `gorm:"primarykey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         Role   `gorm:"type:varchar(10);not null"`
	CreatedAt    time.Time
}
