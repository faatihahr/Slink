package database

import (
	"context"
	"fmt"
	"log"

	"slink-backend/internal/models"

	supa "github.com/lengzuo/supa"
)

type SupabaseClient struct {
	client *supa.Client
}

func ConnectSupabase(url, key, projectRef string) (*SupabaseClient, error) {
	config := supa.Config{
		ApiKey:     key,
		ProjectRef: projectRef,
		Debug:      true,
	}
	client, err := supa.New(config)
	if err != nil {
		return nil, err
	}

	log.Println("Supabase connected successfully")
	return &SupabaseClient{client: client}, nil
}

func (s *SupabaseClient) GetClient() *supa.Client {
	return s.client
}

// User operations
func (s *SupabaseClient) CreateUser(user *models.User) error {
	var result models.User
	err := s.client.DB.From("users").Insert(user).Execute(context.Background(), &result)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	user.ID = result.ID
	user.CreatedAt = result.CreatedAt
	user.UpdatedAt = result.UpdatedAt
	return nil
}

func (s *SupabaseClient) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := s.client.DB.From("users").Select("*").Eq("id", userID).Single().Execute(context.Background(), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return &user, nil
}

func (s *SupabaseClient) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.client.DB.From("users").Select("*").Eq("email", email).Single().Execute(context.Background(), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return &user, nil
}
