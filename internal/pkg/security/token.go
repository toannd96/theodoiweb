package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenDetails struct {
	UserID      string
	AccessToken string
	AccessUUID  string
	AtExpires   int64
	// RefreshToken string
	// RefreshUUID  string
	// RtExpires    int64
}

func CreateToken(userID string) (*TokenDetails, error) {
	td := &TokenDetails{
		AtExpires:  time.Now().Add(time.Hour * 24).Unix(),
		AccessUUID: uuid.New().String(),
		// RtExpires:   time.Now().Add(time.Hour * 24).Unix(),
		// RefreshUUID: uuid.NewV4().String(),
	}

	var err error

	// access token
	atClaims := jwt.MapClaims{
		"access_uuid": td.AccessUUID,
		"user_id":     userID,
		"exp":         td.AtExpires,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// refresh token
	// rtClaims := jwt.MapClaims{
	// 	"refresh_uuid": td.RefreshUUID,
	// 	"user_id":      userID,
	// 	"exp":          td.RtExpires,
	// }
	// rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	// td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	// if err != nil {
	// 	return nil, err
	// }

	return td, nil
}
