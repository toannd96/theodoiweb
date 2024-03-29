package website

import (
	"analytics-api/configs"
	"analytics-api/internal/app/auth"
	"analytics-api/internal/pkg/middleware"
	"analytics-api/internal/pkg/security"
	str "analytics-api/internal/pkg/string"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type httpDelivery struct {
	websiteUseCase UseCase
	authUsecase    auth.UseCase
}

// var validate = validator.New()

// type RequestWebsite struct {
// 	Name string `json:"name" validate:"required,min=2,max=100"`
// 	URL  string `json:"url" validate:"required,min=3"`
// }

// InitRoutes ...
func (instance *httpDelivery) InitRoutes(r *gin.RouterGroup) {
	websiteRoutes := r.Group("/website")
	{
		websiteRoutes.GET("/dashboard", middleware.JWTMiddleware(), instance.Dashboard)
		websiteRoutes.GET("/:website_id", middleware.JWTMiddleware(), instance.GetWebsite)
		websiteRoutes.GET("/list", middleware.JWTMiddleware(), instance.GetAllWebsite)
		websiteRoutes.GET("/tracking/:website_id", middleware.JWTMiddleware(), instance.Tracking)

		websiteRoutes.GET("/add", middleware.JWTMiddleware(), instance.ShowAddWebsite)
		websiteRoutes.POST("/add", middleware.JWTMiddleware(), instance.AddWebsite)

		websiteRoutes.GET("/delete/:website_id", middleware.JWTMiddleware(), instance.DeleteWebsite)
	}
}

func (instance *httpDelivery) Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{})
}

func (instance *httpDelivery) ShowAddWebsite(c *gin.Context) {
	c.HTML(http.StatusOK, "website.html", gin.H{})
}

// Tracking guide tracking code to website
func (instance *httpDelivery) Tracking(c *gin.Context) {
	websiteID := c.Param("website_id")

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

	c.HTML(http.StatusOK, "tracking.html", gin.H{
		"URL":       configs.AppURL,
		"UserID":    userID,
		"WebsiteID": websiteID,
	})
}

func (instance *httpDelivery) GetWebsite(c *gin.Context) {
	websiteID := c.Param("website_id")
	var aWebsite website
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

	getWebsiteErr := instance.websiteUseCase.GetWebsite(userID, websiteID, &aWebsite)
	if getWebsiteErr != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	c.JSON(http.StatusOK, aWebsite)
}

func (instance *httpDelivery) GetAllWebsite(c *gin.Context) {
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

	websites, err := instance.websiteUseCase.GetAllWebsite(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	if len(*websites) == 0 {
		c.HTML(http.StatusOK, "website.html", gin.H{})
	}

	c.HTML(http.StatusOK, "websites.html", gin.H{
		"Websites": websites,
	})
}

func (instance *httpDelivery) AddWebsite(c *gin.Context) {
	url := c.PostForm("url")
	category := c.PostForm("category")

	hostName, err := str.ParseURL(url)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
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

	count, err := instance.websiteUseCase.FindWebsite(userID, hostName)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"msg": "this website already exists"})
		return
	} else {

		websiteID := str.GetMD5Hash(hostName)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			return
		}

		createdAt := time.Now().Format("2006-01-02, 15:04:05")

		aWebsite := website{
			ID:        websiteID,
			UserID:    userID,
			Category:  category,
			HostName:  hostName,
			URL:       url,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		}

		insertErr := instance.websiteUseCase.InsertWebsite(userID, aWebsite)
		if insertErr != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/website/list")
	}
}

func (instance *httpDelivery) DeleteWebsite(c *gin.Context) {
	websiteID := c.Param("website_id")
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

	deleteWebsiteErr := instance.websiteUseCase.DeleteWebsite(userID, websiteID)
	if deleteWebsiteErr != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	deleteSessionErr := instance.websiteUseCase.DeleteSession(userID, websiteID)
	if deleteSessionErr != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/website/list")
}
