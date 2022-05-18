package event

import (
	"analytics-api/models"
)

// UseCase ...
type UseCase interface {
	GetCountEventBySessionID(sessionID string, events *models.Events) (int, error)
	GetEventBySessionID(sessionID string, events *models.Events, limit, skip int) error
	InsertEvent(sessionID string) error
	FindSessionID(id string) (int64, error)
	UpdateEvent(events models.Events) error
}

type useCase struct {
	repo Repository
}

// NewUseCase ...
func NewUseCase() UseCase {
	return &useCase{
		repo: NewRepository(),
	}
}

// GetCountEventBySessionID get count event of session by session id
func (instance *useCase) GetCountEventBySessionID(sessionID string, events *models.Events) (int, error) {
	count, err := instance.repo.GetCountEventBySessionID(sessionID, events)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetEventBySessionID get limit event of session by session id
func (instance *useCase) GetEventBySessionID(sessionID string, events *models.Events, limit, skip int) error {
	err := instance.repo.GetEventBySessionID(sessionID, events, limit, skip)
	if err != nil {
		return err
	}
	return nil
}

// InsertEvent insert session id by event
func (instance *useCase) InsertEvent(sessionID string) error {
	err := instance.repo.Insert(sessionID)
	if err != nil {
		return err
	}
	return nil
}

// FindSessionID find session id if not exists count is 0
func (instance *useCase) FindSessionID(id string) (int64, error) {
	count, err := instance.repo.FindSessionID(id)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateEvent update event of session by session id
func (instance *useCase) UpdateEvent(events models.Events) error {
	err := instance.repo.Update(events)
	if err != nil {
		return err
	}
	return nil
}
