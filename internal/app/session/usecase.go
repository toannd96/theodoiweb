package session

import (
	"analytics-api/models"
)

// UseCase ...
type UseCase interface {
	GetAllSession(sessions models.Sessions) (models.Sessions, error)
	GetSessionByID(sessionID string, session *models.Session) error
	InsertSession(session models.Session) error
	FindSessionID(id string) (int64, error)
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

// GetSessionByID get session by session id
func (instance *useCase) GetSessionByID(sessionID string, session *models.Session) error {
	err := instance.repo.GetSessionByID(sessionID, session)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session
func (instance *useCase) GetAllSession(sessions models.Sessions) (models.Sessions, error) {
	sessions, err := instance.repo.GetAllSession(sessions)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// InsertSession insert session
func (instance *useCase) InsertSession(session models.Session) error {
	err := instance.repo.Insert(session)
	if err != nil {
		return err
	}
	return nil
}

// FindSessionID find session id to check exists
func (instance *useCase) FindSessionID(id string) (int64, error) {
	count, err := instance.repo.FindSessionID(id)
	if err != nil {
		return 0, err
	}
	return count, nil
}
