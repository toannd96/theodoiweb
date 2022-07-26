package auth

import "analytics-api/internal/pkg/security"

// UseCase ...
type UseCase interface {
	InsertAuth(userID string, tokenDetails *security.TokenDetails) error
	GetAuth(accessUUID string) (string, error)
	DeleteAccessToken(accessUUID string) error
	DeleteRefreshToken(refresUUID string) error
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

func (instance *useCase) InsertAuth(userID string, tokenDetails *security.TokenDetails) error {
	err := instance.repo.InsertAuth(userID, tokenDetails)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) GetAuth(accessUUID string) (string, error) {
	userID, err := instance.repo.GetAuth(accessUUID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (instance *useCase) DeleteAccessToken(accessUUID string) error {
	err := instance.repo.DeleteAccessToken(accessUUID)
	if err != nil {
		return err
	}
	return nil
}

func (instance *useCase) DeleteRefreshToken(refresUUID string) error {
	err := instance.repo.DeleteRefreshToken(refresUUID)
	if err != nil {
		return err
	}
	return nil
}
