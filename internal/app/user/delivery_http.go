package user

import (
	"analytics-api/internal/app/auth"
	"analytics-api/internal/pkg/security"
	str "analytics-api/internal/pkg/string"
	"analytics-api/models"
	"net/http"
	"time"

	"analytics-api/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type httpDelivery struct {
	userUseCase UseCase
	authUsecase auth.UseCase
}

// var validate = validator.New()

type RequestSignUp struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

type RequestSignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RequestUpdateUser struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

// InitRoutes ...
func (instance *httpDelivery) InitRoutes(r *gin.RouterGroup) {

	r.GET("/signup", instance.ShowSignupPage)
	r.POST("/signup", instance.SignUp)

	r.GET("/login", instance.ShowLoginPage)
	r.POST("/login", instance.Login)

	r.GET("/logout", middleware.JWTMiddleware(), instance.Logout)

	profileRoutes := r.Group("profile")
	{
		profileRoutes.GET("/details", middleware.JWTMiddleware(), instance.GetUser)

		// profileRoutes.GET("/update", middleware.JWTMiddleware(), instance.ShowDetailsUserPage)
		profileRoutes.POST("/update", middleware.JWTMiddleware(), instance.UpdateUser)
	}
}

func (instance *httpDelivery) ShowDetailsUserPage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", gin.H{})
}

func (instance *httpDelivery) ShowSignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func (instance *httpDelivery) SignUp(c *gin.Context) {
	email := c.PostForm("email")
	fullname := c.PostForm("fullname")
	password := c.PostForm("password")

	logrus.Info("email", email)
	logrus.Info("fullname", fullname)
	logrus.Info("password", password)

	count, err := instance.userUseCase.FindUser(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while check for the email"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"msg": "this email already exists"})
		return
	} else {
		hash, err := security.HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while hash password"})
			return
		}

		userID := str.GetMD5Hash(email)

		createdAt := time.Now().Format("2006-01-02, 15:04:05")

		user := models.User{
			ID:        userID,
			FullName:  fullname,
			Email:     email,
			Password:  hash,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		}

		insertErr := instance.userUseCase.InsertUser(user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error occured while insert user"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/login")
	}
}

func (instance *httpDelivery) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (instance *httpDelivery) Login(c *gin.Context) {
	var user models.User

	email := c.PostForm("email")
	password := c.PostForm("password")

	logrus.Info("email", email)
	logrus.Info("password", password)

	err := instance.userUseCase.GetUserByEmail(email, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "email not exists"})
		return
	}

	// check password
	isTheSame := security.DoPasswordsMatch(user.Password, password)
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

	// c.JSON(http.StatusOK, user)
	c.Redirect(http.StatusMovedPermanently, "/website/dashboard")
}

func (instance *httpDelivery) Logout(c *gin.Context) {
	accessToken, err := security.ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "extract access token failed"})
		return
	}

	logrus.Info("delete access token")
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

	logrus.Info("delete refresh token")
	delRtErr := instance.authUsecase.DeleteRefreshToken(refreshToken.RefreshUUID)
	if delRtErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "error occured while del refresh token"})
		return
	}

	logrus.Info("delete access from cookie token")
	c.SetCookie("access_token", "", -1, "", "", false, true)

	logrus.Info("delete refresh from cookie token")
	c.SetCookie("refresh_token", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
	c.Redirect(http.StatusMovedPermanently, "/login")
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

	logrus.Info(userID)

	getUserErr := instance.userUseCase.GetUserByID(userID, &user)
	if getUserErr != nil {
		logrus.Error(getUserErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while get the user"})
		return
	}

	logrus.Info(user)

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"User": user,
	})
}

func (instance *httpDelivery) UpdateUser(c *gin.Context) {
	fullName := c.PostForm("fullname")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirmPassword")

	logrus.Info("fullname", fullName)
	logrus.Info("password", password)

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

	hash, err := security.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while hash password"})
		return
	}

	updatedAt := time.Now().Format("2006-01-02, 15:04:05")

	if fullName == "" {
		if password == confirmPassword {
			user := models.User{
				Password:  hash,
				UpdatedAt: updatedAt,
			}
			err := instance.userUseCase.UpdatePassword(userID, &user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while update user password"})
				return
			}
		} else {
			c.Redirect(http.StatusMovedPermanently, "/profile/details")
		}

	}

	if password == "" {
		user := models.User{
			FullName:  fullName,
			UpdatedAt: updatedAt,
		}
		err := instance.userUseCase.UpdateFullName(userID, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while update user full name"})
			return
		}
	}

	if fullName != "" && password != "" {
		if password == confirmPassword {
			user := models.User{
				FullName:  fullName,
				Password:  hash,
				UpdatedAt: updatedAt,
			}
			updateUserErr := instance.userUseCase.UpdateUser(userID, &user)
			if updateUserErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while update user"})
				return
			}
		} else {
			c.Redirect(http.StatusMovedPermanently, "/profile/details")
		}
	}

	// c.JSON(http.StatusOK, gin.H{"msg": "update user success"})
	c.Redirect(http.StatusMovedPermanently, "/profile/details")
}
