package website

// UseCase ...
type UseCase interface {
	FindWebsite(userID, hostName string) (int64, error)
	FindWebsiteByID(userID, websiteID string) (int64, error)
	InsertWebsite(userID string, aWebsite website) error
	GetWebsite(userID, websiteID string, aWebsite *website) error
	GetAllWebsite(userID string) (*websites, error)
	DeleteWebsite(userID, websiteID string) error
	DeleteSession(userID, websiteID string) error
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

func (instance *useCase) FindWebsite(userID, url string) (int64, error) {
	count, err := instance.repo.FindWebsite(userID, url)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *useCase) FindWebsiteByID(userID, websiteID string) (int64, error) {
	count, err := instance.repo.FindWebsiteByID(userID, websiteID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *useCase) InsertWebsite(userID string, aWebsite website) error {
	err := instance.repo.InsertWebsite(userID, aWebsite)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetWebsite(userID, websiteID string, aWebsite *website) error {
	err := instance.repo.GetWebsite(userID, websiteID, aWebsite)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetAllWebsite(userID string) (*websites, error) {
	websites, err := instance.repo.GetAllWebsite(userID)
	if err != nil {
		return nil, err
	}
	return websites, nil
}

func (instance *useCase) DeleteWebsite(userID, websiteID string) error {
	err := instance.repo.DeleteWebsite(userID, websiteID)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) DeleteSession(userID, websiteID string) error {
	err := instance.repo.DeleteSession(userID, websiteID)
	if err != nil {
		return err
	}
	return nil
}
