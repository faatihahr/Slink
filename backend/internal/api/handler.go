package api

import (
	"fmt"
	"net/http"
	"strings"

	"slink-backend/internal/config"
	"slink-backend/internal/database"
	"slink-backend/internal/models"
	"slink-backend/internal/services"
	"slink-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	urlService  services.URLServiceInterface
	userService services.UserServiceInterface
	qrService   *services.QRService
	config      *config.Config
}

func NewHandler(supabaseClient *database.SupabaseClient, cfg *config.Config) *Handler {
	return &Handler{
		urlService:  services.NewURLServiceSupa(supabaseClient),
		userService: services.NewUserService(supabaseClient),
		qrService:   services.NewQRService(cfg.QRSize),
		config:      cfg,
	}
}

func (h *Handler) ShortenURL(c *gin.Context) {
	var req models.ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if !isValidURL(req.OriginalURL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL. Must start with http:// or https://",
		})
		return
	}

	userID := c.GetString("userID")
	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	url, err := h.urlService.CreateShortURL(req.OriginalURL, req.CustomAlias, userIDPtr)
	if err != nil {
		fmt.Printf("Error in CreateShortURL: %v\n", err)
		if err.Error() == "custom alias already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Custom alias already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create short URL",
		})
		return
	}

	fmt.Printf("URL created successfully: %+v\n", url)

	shortURL := fmt.Sprintf("%s/%s", h.config.BaseURL, url.ShortCode)
	qrCodeURL := fmt.Sprintf("%s/api/qr/%s", h.config.BaseURL, url.ShortCode)

	response := models.ShortenResponse{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		ShortCode:   url.ShortCode,
		ShortURL:    shortURL,
		QRCodeURL:   qrCodeURL,
		CustomAlias: url.CustomAlias,
	}

	fmt.Printf("Response built: %+v\n", response)
	fmt.Printf("Sending JSON response: %+v\n", response)
	c.JSON(http.StatusCreated, response)
}

func (h *Handler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := h.urlService.GetURLByCode(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
		})
		return
	}

	go func() {
		if err := h.urlService.IncrementHitCount(url.ShortCode); err != nil {
			fmt.Printf("Failed to increment hit count: %v\n", err)
		}
	}()

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

func (h *Handler) GenerateQR(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := h.urlService.GetURLByCode(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
		})
		return
	}

	fullURL := fmt.Sprintf("%s/%s", h.config.BaseURL, url.ShortCode)

	qrData, err := h.qrService.GenerateQRCode(fullURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate QR code",
		})
		return
	}

	c.Data(http.StatusOK, "image/png", qrData)
}

func isValidURL(url string) bool {
	if len(url) < 8 {
		return false
	}
	return url[:7] == "http://" || url[:8] == "https://"
}

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}

	response := models.AuthResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	user, err := h.userService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}

	response := models.AuthResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, userResponse)
}

func (h *Handler) GetLinksByUser(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	links, err := h.urlService.GetLinksByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user links",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"links": links,
	})
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
