package session

import (
	"analytics-api/models"
)

// UseCase ...
type UseCase interface {
	InsertSession(session models.Session, events models.Events) error
	GetTotalColumn(sessionID string) (int64, error)
	GetEventLimitBySessionID(sessionID string, limit, offset int, events *models.Events) error
	GetSessionByID(sessionID string, session *models.Session) error
	GetAllSession(listID []string, session models.Session) ([]models.Session, error)
	GetAllSessionID() ([]string, error)
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

// GetTotalColumn get total number column of session record
func (instance *useCase) GetTotalColumn(sessionID string) (int64, error) {
	totalColumn, err := instance.repo.GetTotalColumn(sessionID)
	if err != nil {
		return 0, err
	}
	return totalColumn, nil
}

// GetEventLimitBySessionID get limit event of session by session id
func (instance *useCase) GetEventLimitBySessionID(sessionID string, limit, offset int, events *models.Events) error {
	err := instance.repo.GetEventLimitBySessionID(sessionID, limit, offset, events)
	if err != nil {
		return err
	}
	return nil
}

// GetSessionByID get session by session id
func (instance *useCase) GetSessionByID(sessionID string, session *models.Session) error {
	err := instance.repo.GetSessionByID(sessionID, session)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session from list session id
func (instance *useCase) GetAllSession(listID []string, session models.Session) ([]models.Session, error) {
	sessions, err := instance.repo.GetAllSession(listID, session)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetAllSessionID get all session id from all record
func (instance *useCase) GetAllSessionID() ([]string, error) {
	listID, err := instance.repo.GetAllSessionID()
	if err != nil {
		return nil, err
	}
	return listID, nil
}

// InsertSession insert session
func (instance *useCase) InsertSession(session models.Session, events models.Events) error {
	for _, event := range events.Events {
		err := instance.repo.Insert(session, event)
		if err != nil {
			return err
		}
	}
	return nil
}
