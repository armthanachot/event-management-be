// you must be migrate entity first (put entity in the pkg/db/db.go in function Migrate under db.AutoMigrate and g.ApplyInterface)
// after that run script file script/migrate.sh
package event

import (
	"event-management-system/dto/entity"
	"event-management-system/dto/model"
	"event-management-system/query"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(payload *entity.Event) error
		CreateParticipant(payload []*entity.EventParticipant) error
		FindAll(criteria model.FindAllEventCriteria) ([]*entity.Event, int64, error)
		FindAllParticipant(criteria model.FindAllEventCriteria) ([]*entity.Event, int64, error)
		FindByID(id uint) (*entity.Event, error)
		Update(id uint, payload *entity.Event) error
		UpdateEventParticipant(event_id uint, payload []*entity.EventParticipant) error
		Delete(id uint) error
	}
	repository struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) Repository {
	query.SetDefault(db)
	return &repository{db}
}

func (r *repository) Create(payload *entity.Event) error {
	instance := query.Event
	err := instance.Create(payload)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateParticipant(payload []*entity.EventParticipant) error {
	instance := query.EventParticipant
	err := instance.CreateInBatches(payload, len(payload))
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(criteria model.FindAllEventCriteria) ([]*entity.Event, int64, error) {
	instance := query.Event
	limit := criteria.Limit
	offset := criteria.Offset
	// search := "%" + criteria.Search + "%"
	result, err := instance.Limit(limit).Offset(offset).Find()
	if err != nil {
		return nil, 0, err
	}
	count, err := instance.Count()
	if err != nil {
		return nil, 0, err
	}
	return result, count, nil
}

func (r *repository) FindAllParticipant(criteria model.FindAllEventCriteria) ([]*entity.Event, int64, error) {
	instance := query.Event
	limit := criteria.Limit
	offset := criteria.Offset
	result, err := instance.Preload(field.Associations).Limit(limit).Offset(offset).Find()
	if err != nil {
		return nil, 0, err
	}
	count, err := instance.Count()
	if err != nil {
		return nil, 0, err
	}
	return result, count, nil
}

func (r *repository) FindByID(id uint) (*entity.Event, error) {
	instance := query.Event
	result, err := instance.Preload(field.Associations).Where(instance.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) Update(id uint, payload *entity.Event) error {
	instance := query.Event
	_, err := instance.Where(instance.ID.Eq(id)).Updates(payload)
	if err != nil {
		return err
	}
	return nil
}

//avoid this solution, using temp table then update
func (r *repository) UpdateEventParticipant(event_id uint, payload []*entity.EventParticipant) error {
	tx := r.db.Begin()
	instance := query.EventParticipant

	_, err := instance.Where(instance.EventID.Eq(event_id)).Delete() //soft delete
	if err != nil {
		tx.Rollback()
		return err
	}

	err = instance.CreateInBatches(payload, len(payload))
	if err != nil {
		tx.Rollback()
	}

	tx.Commit()
	return nil
}

func (r *repository) Delete(id uint) error {
	instance := query.Event
	// remove Unscoped if you want to use soft delete
	_, err := instance.Unscoped().Where(instance.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}
	return nil
}
