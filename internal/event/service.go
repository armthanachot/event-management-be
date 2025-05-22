package event

import (
	"event-management-system/dto/entity"
	"event-management-system/dto/model"
)

type (
    Service interface {
        Create(payload *entity.Event) error
        FindAll(criteria model.FindAllEventCriteria) ([]*entity.Event, int64, error)
        FindByID(id uint) (*entity.Event, error)
        Update(id uint, payload *entity.Event) error
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

func (s *service) FindAll(criteria model.FindAllEventCriteria) ([]*entity.Event, int64, error) {
    result, count, err := s.repository.FindAll(criteria)
    if err != nil {
        return nil, 0, err
    }

    return result, count, nil
}

func (s *service) FindByID(id uint) (*entity.Event, error) {
    result, err := s.repository.FindByID(id)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (s *service) Update(id uint, payload *entity.Event) error {
    err := s.repository.Update(id, payload)
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

