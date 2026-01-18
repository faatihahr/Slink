package services

import (
	"context"
	cryptoRand "crypto/rand"
	"fmt"
	"log"
	"time"

	"slink-backend/internal/database"
	"slink-backend/internal/models"
	"slink-backend/internal/utils"

	"github.com/google/uuid"
)

type URLRecord struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CustomAlias *string   `json:"custom_alias,omitempty"`
	UserID      *string   `json:"user_id,omitempty"`
	HitCount    int       `json:"hit_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type URLServiceSupa struct {
	supabase *database.SupabaseClient
}

func NewURLServiceSupa(supabaseClient *database.SupabaseClient) *URLServiceSupa {
	return &URLServiceSupa{supabase: supabaseClient}
}

func (s *URLServiceSupa) CreateShortURL(originalURL string, customAlias *string, userID *string) (*models.URL, error) {
	if !utils.IsValidURL(originalURL) {
		return nil, fmt.Errorf("invalid URL")
	}

	if customAlias != nil {
		client := s.supabase.GetClient()
		var results []URLRecord

		err := client.DB.From("urls").Select("id").Eq("custom_alias", *customAlias).Execute(context.Background(), &results)
		if err != nil {
			return nil, err
		}
		if len(results) > 0 {
			return nil, fmt.Errorf("custom alias already exists")
		}

		results = []URLRecord{}
		err = client.DB.From("urls").Select("id").Eq("short_code", *customAlias).Execute(context.Background(), &results)
		if err != nil {
			return nil, err
		}
		if len(results) > 0 {
			return nil, fmt.Errorf("custom alias already exists")
		}
	}

	var shortCode string
	if customAlias != nil {
		if !utils.IsValidCustomAlias(*customAlias) {
			return nil, fmt.Errorf("invalid custom alias")
		}
		shortCode = *customAlias
	} else {
		var err error
		shortCode, err = s.generateUniqueShortCode()
		if err != nil {
			return nil, err
		}
	}

	urlRecord := URLRecord{
		ID:          uuid.New().String(),
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
		HitCount:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if customAlias != nil {
		urlRecord.CustomAlias = customAlias
	}

	client := s.supabase.GetClient()
	var result URLRecord
	err := client.DB.From("urls").Insert(urlRecord).Execute(context.Background(), &result)
	if err != nil {
		return nil, err
	}

	if result.ID == "" {
		return nil, fmt.Errorf("failed to create URL")
	}

	return &models.URL{
		ID:          result.ID,
		OriginalURL: result.OriginalURL,
		ShortCode:   result.ShortCode,
		CustomAlias: result.CustomAlias,
		HitCount:    result.HitCount,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (s *URLServiceSupa) GetURLByCode(code string) (*models.URL, error) {
	client := s.supabase.GetClient()
	var results []URLRecord

	// Query by short_code first
	err := client.DB.From("urls").Select("*").Eq("short_code", code).Execute(context.Background(), &results)
	if err != nil {
		return nil, err
	}

	// If not found by short_code, try custom_alias
	if len(results) == 0 {
		err = client.DB.From("urls").Select("*").Eq("custom_alias", code).Execute(context.Background(), &results)
		if err != nil {
			return nil, err
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	// Convert to models.URL
	return &models.URL{
		ID:          results[0].ID,
		OriginalURL: results[0].OriginalURL,
		ShortCode:   results[0].ShortCode,
		CustomAlias: results[0].CustomAlias,
		HitCount:    results[0].HitCount,
		CreatedAt:   results[0].CreatedAt,
		UpdatedAt:   results[0].UpdatedAt,
	}, nil
}

func (s *URLServiceSupa) IncrementHitCount(shortCode string) error {
	client := s.supabase.GetClient()

	// Get current hit count
	var results []struct {
		HitCount int `json:"hit_count"`
	}
	err := client.DB.From("urls").Select("hit_count").Eq("short_code", shortCode).Execute(context.Background(), &results)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return fmt.Errorf("URL not found")
	}

	// Update hit count
	newHitCount := results[0].HitCount + 1
	updateData := struct {
		HitCount  int       `json:"hit_count"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		HitCount:  newHitCount,
		UpdatedAt: time.Now(),
	}

	err = client.DB.From("urls").Update(updateData).Eq("short_code", shortCode).Execute(context.Background(), nil)

	return err
}

func (s *URLServiceSupa) GetLinksByUser(userID string) ([]models.URL, error) {
	client := s.supabase.GetClient()
	var results []URLRecord

	err := client.DB.From("urls").Select("*").Eq("user_id", userID).Execute(context.Background(), &results)
	if err != nil {
		return nil, err
	}

	// Convert to models.URL slice
	urls := make([]models.URL, len(results))
	for i, result := range results {
		urls[i] = models.URL{
			ID:          result.ID,
			OriginalURL: result.OriginalURL,
			ShortCode:   result.ShortCode,
			CustomAlias: result.CustomAlias,
			HitCount:    result.HitCount,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
	}

	return urls, nil
}

func (s *URLServiceSupa) generateUniqueShortCode() (string, error) {
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		code := generateRandomCodeSupa()

		// Check if code exists
		client := s.supabase.GetClient()
		var results []URLRecord
		err := client.DB.From("urls").Select("id").Eq("short_code", code).Execute(context.Background(), &results)
		if err != nil {
			log.Printf("Error checking short code existence: %v", err)
			continue
		}

		if len(results) == 0 {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique short code after %d attempts", maxAttempts)
}

func generateRandomCodeSupa() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	b := make([]byte, length)
	cryptoRand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b)
}
