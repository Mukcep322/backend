package service

import (
	"context"
	"fmt"
	"time"

	"trainers-backend/internal/models"
	"trainers-backend/internal/repository"
	"trainers-backend/internal/telegram"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	userRepo  *repository.UserRepo
	redis     *redis.Client
	jwtSecret string
	botToken  string
}

func NewAuthService(userRepo *repository.UserRepo, redis *redis.Client, jwtSecret, botToken string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		redis:     redis,
		jwtSecret: jwtSecret,
		botToken:  botToken,
	}
}

type Claims struct {
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
	TelegramID int64  `json:"telegram_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) AuthenticateTelegram(ctx context.Context, initData string) (*models.User, string, error) {
	// Валидация initData
	userData, err := telegram.ValidateInitData(initData, s.botToken)
	if err != nil {
		return nil, "", fmt.Errorf("invalid init data: %w", err)
	}

	// Парсим telegram_id
	var telegramID int64
	fmt.Sscanf(userData["id"], "%d", &telegramID)

	// Ищем или создаем пользователя
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		// Создаём нового пользователя
		user = &models.User{
			TelegramID: telegramID,
			Username:   userData["username"],
			FirstName:  userData["first_name"],
			LastName:   userData["last_name"],
			Role:       "client",
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, "", err
		}
	} else {
		// Обновляем данные существующего пользователя
		user.Username = userData["username"]
		user.FirstName = userData["first_name"]
		user.LastName = userData["last_name"]
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, "", err
		}
	}

	// Генерируем JWT
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	// Сохраняем в Redis (опционально)
	s.redis.Set(ctx, "session:"+user.ID, token, 24*time.Hour)

	return user, token, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:     user.ID,
		Role:       user.Role,
		TelegramID: user.TelegramID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
