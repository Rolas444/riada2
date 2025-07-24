package domain

import (
	"time"
)

// Sex define los valores posibles para el sexo de una persona.
type Sex string

const (
	Female Sex = "F"
	Male   Sex = "M"
)

// DocType define los valores posibles para un tipo de documento.
type DocType string

const (
	DNI      DocType = "DNI"
	CE       DocType = "CE"
	Passport DocType = "Passport"
)

// Person representa la información personal de un usuario.
type Person struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     *uint     `gorm:"index"` // Clave foránea del usuario que realizó la última modificación.
	User       *User     // GORM usará UserID como clave foránea por convención.
	Name       string    `gorm:"not null"`
	MiddleName string    `gorm:"not null"`
	LastName   string    `gorm:"not null"`
	Sex        Sex       `gorm:"type:varchar(1);not null"`
	Birthday   time.Time `gorm:"not null"`
	DocNumber  *string   `gorm:"type:varchar(20)"` // Permite nulos
	TypeDoc    *DocType  `gorm:"type:varchar(10)"` // Permite nulos
	Email      *string   `gorm:"type:varchar(100);uniqueIndex"`
	Photo      *string   `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para el modelo Person.
func (p *Person) TableName() string {
	return "persons"
}
