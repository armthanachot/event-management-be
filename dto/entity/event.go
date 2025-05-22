package entity

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string    `json:"name"`
	DateTime    time.Time `json:"date_time"`
	OrganizerId uint      `json:"organizer_id"`
	IsAvailable bool      `json:"is_available"`
	IsCancelled bool      `json:"is_cancelled"`

	Organizer *User `json:"organizer" gorm:"foreignKey:OrganizerId"`
}