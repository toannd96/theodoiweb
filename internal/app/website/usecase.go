package website

import (
	"analytics-api/models"
)

// UseCase ...
type UseCase interface {
	FindWebsite(userID, hostName string) (int64, error)
	FindWebsiteByID(userID, websiteID string) (int64, error)
	InsertWebsite(userID string, website models.Website) error
	GetWebsite(userID, websiteID string, website *models.Website) error
	GetAllWebsite(userID string) (*models.Websites, error)
	UpdateNameWebsite(userID, websiteID string, website *models.Website) error
	UpdateURLWebsite(userID, websiteID string, website *models.Website) error
	UpdateTrackedWebsite(userID, websiteID string, website *models.Website) error
	UpdateWebsite(userID, websiteID string, website *models.Website) error
	DeleteWebsite(userID, websiteID string) error
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

func (instance *useCase) InsertWebsite(userID string, website models.Website) error {
	err := instance.repo.InsertWebsite(userID, website)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetWebsite(userID, websiteID string, website *models.Website) error {
	err := instance.repo.GetWebsite(userID, websiteID, website)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetAllWebsite(userID string) (*models.Websites, error) {
	websites, err := instance.repo.GetAllWebsite(userID)
	if err != nil {
		return nil, err
	}
	return websites, nil
}

func (instance *useCase) UpdateNameWebsite(userID, websiteID string, website *models.Website) error {
	err := instance.repo.UpdateNameWebsite(userID, websiteID, website)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdateURLWebsite(userID, websiteID string, website *models.Website) error {
	err := instance.repo.UpdateURLWebsite(userID, websiteID, website)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdateTrackedWebsite(userID, websiteID string, website *models.Website) error {
	err := instance.repo.UpdateTrackedWebsite(userID, websiteID, website)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdateWebsite(userID, websiteID string, website *models.Website) error {
	err := instance.repo.UpdateWebsite(userID, websiteID, website)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) DeleteWebsite(userID, websiteID string) error {
	err := instance.repo.DeleteWebsite(userID, websiteID)
	if err != nil {
		return err
	}
	return nil
}
