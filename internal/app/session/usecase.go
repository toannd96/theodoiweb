package session

import (
	"analytics-api/models"
)

// UseCase ...
type UseCase interface {
	GetAllSession(sessions models.Sessions) (models.Sessions, error)
	GetSession(id string, session *models.Session) error
	InsertSession(session models.Session) error
	FindSession(id string) (int64, error)

	GetCountEvent(id string, session *models.Session) (int, error)
	GetEvent(id string, session *models.Session, limit, skip int) error
	UpdateEvent(session models.Session) error
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

// GetSession get session by session id
func (instance *useCase) GetSession(id string, session *models.Session) error {
	err := instance.repo.GetSession(id, session)
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
	err := instance.repo.InsertSession(session)
	if err != nil {
		return err
	}
	return nil
}

// FindSession find session id to check exists
func (instance *useCase) FindSession(id string) (int64, error) {
	count, err := instance.repo.FindSession(id)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCountEvent get count event of session by session id
func (instance *useCase) GetCountEvent(id string, session *models.Session) (int, error) {
	count, err := instance.repo.GetCountEvent(id, session)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetEvent get limit event of session by session id
func (instance *useCase) GetEvent(id string, session *models.Session, limit, skip int) error {
	err := instance.repo.GetEvent(id, session, limit, skip)
	if err != nil {
		return err
	}
	return nil
}

// UpdateEvent update event of session by session id
func (instance *useCase) UpdateEvent(session models.Session) error {
	err := instance.repo.UpdateEvent(session)
	if err != nil {
		return err
	}
	return nil
}
