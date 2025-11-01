package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	RefreshToken(ctx context.Context, tokenString string) (string, error)
	ValidateToken(tokenString string, expectedType string) (uuid.UUID, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo repository.UserRepository, secret string, accessTTL time.Duration, refreshTTL time.Duration) AuthService {
	return &authService{
		userRepo:   repo,
		jwtSecret:  secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	id, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {

		return nil, err
	}
	user.ID = id

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := s.generateToken(user.ID, s.accessTTL, "access")
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateToken(user.ID, s.refreshTTL, "refresh")
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) generateToken(userID uuid.UUID, ttl time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"sub":        userID.String(),
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(ttl).Unix(),
		"token_type": tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ValidateToken(tokenString string, expectedType string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType != expectedType {
			return uuid.Nil, errors.New("invalid token type")
		}

		userIDStr, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, errors.New("invalid token claims (sub)")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, errors.New("invalid user ID in token")
		}

		return userID, nil
	}

	return uuid.Nil, errors.New("invalid token")
}

func (s *authService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	userID, err := s.ValidateToken(tokenString, "refresh")
	if err != nil {
		return "", err
	}

	newAccessToken, err := s.generateToken(userID, s.accessTTL, "access")
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
