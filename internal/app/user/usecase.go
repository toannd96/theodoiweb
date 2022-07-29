package user

import (
	"analytics-api/models"
)

// UseCase ...
type UseCase interface {
	FindUser(email string) (int64, error)
	InsertUser(user models.User) error
	GetUserByEmail(email string, user *models.User) error
	GetUserByID(userID string, user *models.User) error
	UpdateUser(userID string, user *models.User) error
	UpdateFullName(userID string, user *models.User) error
	UpdatePassword(userID string, user *models.User) error
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

func (instance *useCase) FindUser(email string) (int64, error) {
	count, err := instance.repo.FindUser(email)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *useCase) InsertUser(user models.User) error {
	err := instance.repo.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetUserByEmail(email string, user *models.User) error {
	err := instance.repo.GetUserByEmail(email, user)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetUserByID(userID string, user *models.User) error {
	err := instance.repo.GetUserByID(userID, user)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdateFullName(userID string, user *models.User) error {
	err := instance.repo.UpdateFullName(userID, user)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdatePassword(userID string, user *models.User) error {
	err := instance.repo.UpdatePassword(userID, user)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) UpdateUser(userID string, user *models.User) error {
	err := instance.repo.UpdateUser(userID, user)
	if err != nil {
		return err
	}
	return nil
}
