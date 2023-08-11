package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/joho/godotenv"
	"os"
	"time"
	balance_app "users-balance-monitoring"
	"users-balance-monitoring/pkg/repository"
)

const tokenTl = 12 * time.Hour

var (
	salt       = os.Getenv("SALT")
	signingKey = os.Getenv("SIGNING_KEY")
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"userId"`
}
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user balance_app.User) (int, error) {
	user.Password = passwordHash(user.Password)
	return s.repo.CreateUser(user)
}

func passwordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, passwordHash(password))

	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(acsessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(acsessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims error")
	}
	return claims.UserId, nil
}
