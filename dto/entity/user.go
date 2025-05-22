package entity

import (
	"event-management-system/pkg/constant"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string  `json:"name"`
	Email  string  `json:"email" gorm:"unique"`
	Credit float64 `json:"credit"`
	Role   string  `json:"role"`
	Events []Event `json:"events" gorm:"foreignKey:OrganizerId"`
}

func (e *User) CheckRole() bool {
	return lo.Contains([]string{constant.ORGANIZER, constant.PARTICIPANT}, e.Role)
}
