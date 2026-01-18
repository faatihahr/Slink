package models

import (
	"time"
)

type URL struct {
	ID          string    `json:"id" db:"id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	CustomAlias *string   `json:"custom_alias,omitempty" db:"custom_alias"`
	HitCount    int       `json:"hit_count" db:"hit_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ShortenRequest struct {
	OriginalURL string  `json:"original_url" binding:"required"`
	CustomAlias *string `json:"custom_alias,omitempty"`
}

type ShortenResponse struct {
	ID          string  `json:"id"`
	OriginalURL string  `json:"original_url"`
	ShortCode   string  `json:"short_code"`
	ShortURL    string  `json:"short_url"`
	QRCodeURL   string  `json:"qr_code_url"`
	CustomAlias *string `json:"custom_alias,omitempty"`
}
