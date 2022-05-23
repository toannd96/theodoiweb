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
	UpdateSession(id string, session models.Session) error

	GetEvents(id string, session *models.Session) (models.Events, error)
	GetEventByLimitSkip(id string, session *models.Session, limit, skip int) error

	GetSessionTimestamp(id string) (int64, error)
	InsertSessionTimestamp(id string, timeStart int64) error
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

// GetEvents get event of session by session id
func (instance *useCase) GetEvents(id string, session *models.Session) (models.Events, error) {
	events, err := instance.repo.GetEvents(id, session)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetEventByLimitSkip get limit event of session by session id
func (instance *useCase) GetEventByLimitSkip(id string, session *models.Session, limit, skip int) error {
	err := instance.repo.GetEventByLimitSkip(id, session, limit, skip)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSession update duration session and event by session id
func (instance *useCase) UpdateSession(id string, session models.Session) error {
	err := instance.repo.UpdateSession(id, session)
	if err != nil {
		return err
	}
	return nil
}

// GetSessionTimestamp get first timestamp of session by id
func (instance *useCase) GetSessionTimestamp(id string) (int64, error) {
	timeStart, err := instance.repo.GetSessionTimestamp(id)
	if err != nil {
		return 0, err
	}
	return timeStart, nil
}

// InsertSessionTimestamp insert first timestamp by session id
func (instance *useCase) InsertSessionTimestamp(id string, timeStart int64) error {
	err := instance.repo.InsertSessionTimestamp(id, timeStart)
	if err != nil {
		return err
	}
	return nil
}
