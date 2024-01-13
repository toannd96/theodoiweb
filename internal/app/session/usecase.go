package session

// UseCase ...
type UseCase interface {
	GetAllSession(userID, websiteID string, listSessionID []string, session session) ([]session, error)
	GetAllSessionID(userID, websiteID string, session session) ([]string, error)

	GetSessionIDToday(userID, websiteID string, session session) ([]string, error)
	GetSession(userID, sessionID string, session *session) error
	GetCountSession(userID, sessionID string) (int64, error)
	InsertSession(session session, events []event) error

	GetEventByLimitSkip(userID, sessionID string, limit, skip int) ([]*event, error)

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
func (instance *useCase) GetSession(userID, sessionID string, aSession *session) error {
	err := instance.repo.GetSession(userID, sessionID, aSession)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session
func (instance *useCase) GetAllSession(userID, websiteID string, listSessionID []string, aSession session) ([]session, error) {
	listSession, err := instance.repo.GetAllSession(userID, websiteID, listSessionID, aSession)
	if err != nil {
		return nil, err
	}
	return listSession, nil
}

// GetAllSessionID get all id of session all time
func (instance *useCase) GetAllSessionID(userID, websiteID string, aSession session) ([]string, error) {
	listSessionID, err := instance.repo.GetAllSessionID(userID, websiteID, aSession)
	if err != nil {
		return nil, err
	}
	return listSessionID, nil
}

// GetAllSessionID get all id of session today
func (instance *useCase) GetSessionIDToday(userID, websiteID string, aSession session) ([]string, error) {
	listSessionID, err := instance.repo.GetSessionIDToday(userID, websiteID, aSession)
	if err != nil {
		return nil, err
	}
	return listSessionID, nil
}

// InsertSession insert session
func (instance *useCase) InsertSession(aSession session, events []event) error {
	for _, event := range events {
		err := instance.repo.InsertSession(aSession, event)
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
func (instance *useCase) GetEventByLimitSkip(userID, sessionID string, limit, skip int) ([]*event, error) {
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
