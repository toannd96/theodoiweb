package auth

import (
	"analytics-api/configs"
	"analytics-api/internal/pkg/security"

	"time"

	"github.com/sirupsen/logrus"
)

// Repository ...
type Repository interface {
	InsertAuth(userID string, tokenDetails *security.TokenDetails) error
	GetAuth(accessUUID string) (string, error)
	DeleteAccessToken(accessUUID string) error
	DeleteRefreshToken(refresUUID string) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

func (instance *repository) InsertAuth(userID string, tokenDetails *security.TokenDetails) error {
	at := time.Unix(tokenDetails.AtExpires, 0)
	// rt := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	errAccess := configs.Redis.Client.Set(tokenDetails.AccessUUID, userID, at.Sub(now)).Err()
	if errAccess != nil {
		logrus.Error("Redis set at error ", errAccess)
		return errAccess
	}
	// errRefresh := configs.Redis.Client.Set(tokenDetails.RefreshUUID, userID, rt.Sub(now)).Err()
	// if errRefresh != nil {
	// 	return errRefresh
	// }
	return nil
}

func (instance *repository) GetAuth(accessUUID string) (string, error) {
	userID, err := configs.Redis.Client.Get(accessUUID).Result()
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (instance *repository) DeleteAccessToken(accessUUID string) error {
	deleteAt, err := configs.Redis.Client.Del(accessUUID).Result()
	if err != nil || deleteAt != 1 {
		return err
	}
	return nil
}

func (instance *repository) DeleteRefreshToken(refresUUID string) error {
	deleteAt, err := configs.Redis.Client.Del(refresUUID).Result()
	if err != nil || deleteAt != 1 {
		return err
	}
	return nil
}
