package event

import (
	"event-management-system/dto/entity"
	"event-management-system/dto/model"

	"github.com/samber/lo"
)

type (
	Service interface {
		Create(payload *entity.Event) error
		CreateParticipant(payload []*entity.EventParticipant) error
		FindAll(criteria model.FindAllEventCriteria) ([]*entity.Event, int64, error)
		FindAllEventParticipant(criteria model.FindAllEventCriteria) ([]model.FindEventParticipantResp, int64, error)
		FindByID(id uint) (model.FindEventParticipantResp, error)
		Update(id uint, payload *entity.Event) error
		UpdateEventParticipant(event_id uint, payload []*entity.EventParticipant) error
		Delete(id uint) error
	}

	service struct {
		repository Repository
	}
)

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) Create(payload *entity.Event) error {
	err := s.repository.Create(payload)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateParticipant(payload []*entity.EventParticipant) error {
	err := s.repository.CreateParticipant(payload)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindAll(criteria model.FindAllEventCriteria) ([]*entity.Event, int64, error) {
	result, count, err := s.repository.FindAll(criteria)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}

func (s *service) FindAllEventParticipant(criteria model.FindAllEventCriteria) ([]model.FindEventParticipantResp, int64, error) {
	result, count, err := s.repository.FindAllParticipant(criteria)
	if err != nil {
		return nil, 0, err
	}
	resp := lo.Map(result, func(item *entity.Event, _ int) model.FindEventParticipantResp {
		return buildEventParticipant(item)
	})
	return resp, count, nil
}

func (s *service) FindByID(id uint) (model.FindEventParticipantResp, error) {
	result, err := s.repository.FindByID(id)
	if err != nil {
		return model.FindEventParticipantResp{}, err
	}

	return buildEventParticipant(result), nil
}

func (s *service) Update(id uint, payload *entity.Event) error {
	err := s.repository.Update(id, payload)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateEventParticipant(event_id uint, payload []*entity.EventParticipant) error {
	err := s.repository.UpdateEventParticipant(event_id, payload)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func buildEventParticipant(event *entity.Event) model.FindEventParticipantResp {
	return model.FindEventParticipantResp{
		EventID:   event.ID,
		EventName: event.Name,
		EventDate: event.DateTime,
		Available: event.IsAvailable,
		Cancelled: event.IsCancelled,
		Organizer: model.Participant{
			UserID: event.Organizer.ID,
			Role:   event.Organizer.Role,
			Name:   event.Organizer.Name,
			Email:  event.Organizer.Email,
		},
		Participants: lo.Map(event.Participants, func(item *entity.User, _ int) model.Participant {
			return model.Participant{
				UserID: item.ID,
				Role:   item.Role,
				Name:   item.Name,
				Email:  item.Email,
			}
		},
		),
	}
}
