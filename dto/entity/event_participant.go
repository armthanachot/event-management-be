package entity

type EventParticipant struct {
	ID            uint `json:"id" gorm:"primaryKey"`
	EventID       uint `json:"event_id"`
	UserID        uint `json:"user_id"`

	Event *Event `json:"event" gorm:"foreignKey:EventID"`
	User  *User  `json:"user" gorm:"foreignKey:UserID"`
}
