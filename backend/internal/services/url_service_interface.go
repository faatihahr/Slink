package services

import "slink-backend/internal/models"

type URLServiceInterface interface {
	CreateShortURL(originalURL string, customAlias *string, userID *string) (*models.URL, error)
	GetURLByCode(code string) (*models.URL, error)
	IncrementHitCount(shortCode string) error
	GetLinksByUser(userID string) ([]models.URL, error)
}
