package domain

import "time"

// MembershipState define el tipo para el estado de la membresía.
type MembershipState string

const (
	Active   MembershipState = "A" // Activa
	Inactive MembershipState = "I" // Inactivo
)

// Membership representa la entidad de membresía en el sistema.
// Tiene una relación 1 a 1 con Person.
type Membership struct {
	ID               uint            `gorm:"primaryKey;autoIncrement"`
	PersonID         uint            `gorm:"uniqueIndex;not null"`        // FK a Person, único para relación 1:1
	StartedAt        *time.Time      `gorm:"nullable"`                    // Fecha de inicio de membresía
	MembershipSigned bool            `gorm:"default:false"`               // Si la membresía está firmada
	State            MembershipState `gorm:"type:varchar(1);default:'A'"` // Estado: A=Activa, O=Inactiva, S=Suspendida
	Transferred      bool            `gorm:"default:false"`               // Si fue transferida
	NameLastChurch   *string         `gorm:"nullable"`                    // Nombre de la iglesia anterior
	Baptized         bool            `gorm:"default:false"`               // Si está bautizado
	BaptismDate      *time.Time      `gorm:"nullable"`                    // Fecha de bautismo
	Person           Person          `gorm:"foreignKey:PersonID"`         // Relación 1:1 con Person
	CreatedAt        time.Time       `gorm:"autoCreateTime"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime"`
}
