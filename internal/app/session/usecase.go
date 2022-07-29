package session

import (
	"analytics-api/models"
)

// UseCase ...
type UseCase interface {
	GetAllSession(userID, websiteID string, listSessionID []string, session models.Session) ([]models.Session, error)
	GetAllSessionID(userID, websiteID string, session models.Session) ([]string, error)
	GetSession(userID, sessionID string, session *models.Session) error
	GetCountSession(userID, sessionID string) (int64, error)
	InsertSession(session models.Session, events []models.Event) error

	GetEventByLimitSkip(userID, sessionID string, limit, skip int) ([]*models.Event, error)

	GetSessionTimestamp(sessionID string) (int64, error)
	InsertSessionTimestamp(sessionID string, timeStart int64) error
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
func (instance *useCase) GetSession(userID, sessionID string, session *models.Session) error {
	err := instance.repo.GetSession(userID, sessionID, session)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session
func (instance *useCase) GetAllSession(userID, websiteID string, listSessionID []string, session models.Session) ([]models.Session, error) {
	listSession, err := instance.repo.GetAllSession(userID, websiteID, listSessionID, session)
	if err != nil {
		return nil, err
	}
	return listSession, nil
}

// GetAllSessionID get all id of session
func (instance *useCase) GetAllSessionID(userID, websiteID string, session models.Session) ([]string, error) {
	listSessionID, err := instance.repo.GetAllSessionID(userID, websiteID, session)
	if err != nil {
		return nil, err
	}
	return listSessionID, nil
}

// InsertSession insert session
func (instance *useCase) InsertSession(session models.Session, events []models.Event) error {
	for _, event := range events {
		err := instance.repo.InsertSession(session, event)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCountSession find session id to check exists
func (instance *useCase) GetCountSession(userID, sessionID string) (int64, error) {
	count, err := instance.repo.GetCountSession(userID, sessionID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetEventByLimitSkip get limit event of session by session id
func (instance *useCase) GetEventByLimitSkip(userID, sessionID string, limit, skip int) ([]*models.Event, error) {
	events, err := instance.repo.GetEventByLimitSkip(userID, sessionID, limit, skip)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// GetSessionTimestamp get first timestamp of session by id
func (instance *useCase) GetSessionTimestamp(sessionID string) (int64, error) {
	timeStart, err := instance.repo.GetSessionTimestamp(sessionID)
	if err != nil {
		return 0, err
	}
	return timeStart, nil
}

// InsertSessionTimestamp insert first timestamp by session id
func (instance *useCase) InsertSessionTimestamp(sessionID string, timeStart int64) error {
	err := instance.repo.InsertSessionTimestamp(sessionID, timeStart)
	if err != nil {
		return err
	}
	return nil
}
