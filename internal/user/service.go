package user

import (
	"event-management-system/dto/entity"
	"event-management-system/dto/model"
)

type (
    Service interface {
        Create(payload *entity.User) error
        FindAll(criteria model.FindAllUserCriteria) ([]*entity.User, int64, error)
        FindByID(id uint) (*entity.User, error)
        Update(id uint, payload *entity.User) error
    }


    service struct {
        repository Repository
    }
)

func NewService(repository Repository) Service {
    return &service{repository}
}

func (s *service) Create(payload *entity.User) error {
    err := s.repository.Create(payload)
    if err != nil {
        return err
    }
    return nil
}

func (s *service) FindAll(criteria model.FindAllUserCriteria) ([]*entity.User, int64, error) {
    result, count, err := s.repository.FindAll(criteria)
    if err != nil {
        return nil, 0, err
    }

    resp := []*entity.User{}
    for _, v := range result {
        resp = append(resp, v)
    }

    return resp, count, nil
}

func (s *service) FindByID(id uint) (*entity.User, error) {
    user, err := s.repository.FindByID(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (s *service) Update(id uint, payload *entity.User) error {
    err := s.repository.Update(id, payload)
    if err != nil {
        return err
    }
    return nil
}