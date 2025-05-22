package model

import (
	"event-management-system/dto/entity"
	"time"

	"github.com/samber/lo"
)

type FindAllEventCriteria struct {
	Limit  int    `json:"limit" validate:"gte=0"`
	Offset int    `json:"offset" validate:"gte=0"`
	Search string `json:"search"`
}

type FindAllEventParticipantCriteria struct {
	Limit  int `json:"limit" validate:"gte=0"`
	Offset int `json:"offset" validate:"gte=0"`
}

type CreateEventParticipant struct {
	EventID      uint   `json:"event_id" validate:"required"`
	Participants []uint `json:"participants" validate:"required"`
}

type FindEventParticipantResp struct {
	EventID      uint          `json:"event_id"`
	EventName    string        `json:"event_name"`
	EventDate    time.Time     `json:"event_date"`
	Available    bool          `json:"available"`
	Cancelled    bool          `json:"cancelled"`
	Organizer    Participant   `json:"organizer"`
	Participants []Participant `json:"participants"`
}

type Participant struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (c *CreateEventParticipant) UniqueParticipants() CreateEventParticipant {
	return CreateEventParticipant{
		EventID:      c.EventID,
		Participants: lo.Uniq(c.Participants),
	}
}

func (c *CreateEventParticipant) TransformtoDBPayload() []*entity.EventParticipant {
	var eventParticipants []*entity.EventParticipant
	for _, participant := range c.Participants {
		eventParticipants = append(eventParticipants, &entity.EventParticipant{
			EventID: c.EventID,
			UserID:  participant,
		})
	}
	return eventParticipants
}
