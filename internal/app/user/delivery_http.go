package user

import (
	"analytics-api/internal/app/auth"
	dur "analytics-api/internal/pkg/duration"
	"analytics-api/internal/pkg/security"
	"analytics-api/models"
	"net/http"
	"time"

	"analytics-api/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type httpDelivery struct {
	userUseCase UseCase
	authUsecase auth.UseCase
}

var validate = validator.New()

type RequestSignUp struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	FullName string `json:"full_name,omitempty" validate:"required,min=2,max=100"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
}

type RequestSignIn struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
}

type RequestUpdateUser struct {
	FullName string `json:"full_name,omitempty"`
	Password string `json:"password,omitempty"`
}

// InitRoutes ...
func (instance *httpDelivery) InitRoutes(r *gin.RouterGroup) {

	// Register routes user
	userRoutes := r.Group("user")
	{
		userRoutes.POST("/signup", instance.SignUp)
		userRoutes.POST("/signin", instance.SignIn)
		userRoutes.POST("/signout", middleware.JWTMiddleware(), instance.SignOut)
		userRoutes.GET("/profile", middleware.JWTMiddleware(), instance.GetUser)
		userRoutes.PUT("/profile/update", middleware.JWTMiddleware(), instance.UpdateUser)
	}
}

func (instance *httpDelivery) SignUp(c *gin.Context) {
	var request RequestSignUp
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	validationErr := validate.Struct(request)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": validationErr.Error()})
		return
	}

	count, err := instance.userUseCase.FindUser(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while check for the email"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"msg": "this email already exists"})
		return
	} else {
		hash, err := security.HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while hash password"})
			return
		}
		userID, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while create uuid for the user"})
			return
		}

		createdAt, err := dur.ParseTime(time.Now().Format("2006-01-02, 15:04:05"))
		if err != nil {
			logrus.Error(c, err)
			return
		}

		user := models.User{
			ID:        userID.String(),
			FullName:  request.FullName,
			Email:     request.Email,
			Password:  hash,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		}

		insertErr := instance.userUseCase.InsertUser(user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while insert user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (instance *httpDelivery) SignIn(c *gin.Context) {
	var request RequestSignIn
	var user models.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	validationErr := validate.Struct(request)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": validationErr.Error()})
		return
	}

	err := instance.userUseCase.GetUserByEmail(request.Email, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "email not exists"})
		return
	}

	// check password
	isTheSame := security.DoPasswordsMatch(user.Password, request.Password)
	if !isTheSame {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "passowrd is incorrect"})
		return
	}

	// create token
	token, err := security.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while create token"})
		return
	}

	InsertAuthErr := instance.authUsecase.InsertAuth(user.ID, token)
	if InsertAuthErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while insert auth token"})
		return
	}

	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken

	// create cookie for client
	c.SetCookie("access_token", token.AccessToken, 900, "/", "localhost", false, true)
	c.SetCookie("refresh_token", token.RefreshToken, 86400, "/", "localhost", false, true)

	c.JSON(http.StatusOK, user)
}

func (instance *httpDelivery) SignOut(c *gin.Context) {
	accessToken, err := security.ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "extract access token failed"})
		return
	}

	delAtErr := instance.authUsecase.DeleteAccessToken(accessToken.AccessUUID)
	if delAtErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "error occured while del access token"})
		return
	}

	refreshToken, err := security.ExtractRefreshTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "extract refresh token failed"})
		return
	}

	delRtErr := instance.authUsecase.DeleteRefreshToken(refreshToken.RefreshUUID)
	if delRtErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "error occured while del refresh token"})
		return
	}

	c.SetCookie("access_token", "", -1, "", "", false, true)
	c.SetCookie("refresh_token", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"msg": "sign out success"})
}

func (instance *httpDelivery) GetUser(c *gin.Context) {
	var user models.User
	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "extract token metadata failed"})
		return
	}

	userID, err := instance.authUsecase.GetAuth(tokenAuth.AccessUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "get token auth failed"})
		return
	}

	getUserErr := instance.userUseCase.GetUserByID(userID, &user)
	if getUserErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while get the user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (instance *httpDelivery) UpdateUser(c *gin.Context) {
	var request RequestUpdateUser
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(request)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "extract token metadata failed"})
		return
	}

	userID, err := instance.authUsecase.GetAuth(tokenAuth.AccessUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "get token auth failed"})
		return
	}

	hash, err := security.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while hash password"})
		return
	}

	updatedAt, err := dur.ParseTime(time.Now().Format("2006-01-02, 15:04:05"))
	if err != nil {
		logrus.Error(c, err)
		return
	}

	if request.FullName == "" {
		user := models.User{
			ID:        userID,
			Password:  hash,
			UpdatedAt: updatedAt,
		}
		err := instance.userUseCase.UpdatePassword(userID, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while update user password"})
			return
		}
	}

	if request.Password == "" {
		user := models.User{
			ID:        userID,
			FullName:  request.FullName,
			UpdatedAt: updatedAt,
		}
		err := instance.userUseCase.UpdateFullName(userID, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while update user full name"})
			return
		}
	}

	if request.FullName != "" && request.Password != "" {
		user := models.User{
			ID:        userID,
			FullName:  request.FullName,
			Password:  hash,
			UpdatedAt: updatedAt,
		}
		updateUserErr := instance.userUseCase.UpdateUser(userID, &user)
		if updateUserErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while update user"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"msg": "update user success"})
}
