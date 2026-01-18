package services

import (
	"fmt"
	"slink-backend/internal/database"
	"slink-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Register(req models.RegisterRequest) (*models.User, error)
	Login(req models.LoginRequest) (*models.User, error)
	GetUserByID(userID string) (*models.User, error)
}

type UserService struct {
	db *database.SupabaseClient
}

func NewUserService(db *database.SupabaseClient) UserServiceInterface {
	return &UserService{db: db}
}

func (s *UserService) Register(req models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
	}

	// Insert into database
	err = s.db.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *UserService) Login(req models.LoginRequest) (*models.User, error) {
	// Get user by email
	user, err := s.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	return s.db.GetUserByID(userID)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.db.GetUserByEmail(email)
}
