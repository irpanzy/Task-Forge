package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/config"
)

func GenerateToken(userID int64, role, email string, publicID uuid.UUID) (string, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return "", errors.New("JWT_SECRET belum diatur di .env")
	}

	duration := ParseDuration(config.AppConfig.JWTExpired, 15*time.Minute)

	claims := jwt.MapClaims{
		"user_id":   userID,
		"public_id": publicID.String(),
		"role":      role,
		"email":     email,
		"exp":       time.Now().Add(duration).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID int64, publicID uuid.UUID) (string, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return "", errors.New("JWT_SECRET belum diatur di .env")
	}

	duration := ParseDuration(config.AppConfig.RefreshTokenExpired, 7*24*time.Hour)

	claims := jwt.MapClaims{
		"user_id":   userID,
		"public_id": publicID.String(),
		"exp":       time.Now().Add(duration).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return nil, errors.New("JWT_SECRET belum diatur di .env")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak valid: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid atau kadaluarsa")
	}

	return claims, nil
}

func ParseDuration(str string, fallback time.Duration) time.Duration {
	if str == "" {
		return fallback
	}

	if strings.HasSuffix(str, "d") {
		daysStr := strings.TrimSuffix(str, "d")
		if days, err := strconv.Atoi(daysStr); err == nil {
			return time.Duration(days) * 24 * time.Hour
		}
	}

	d, err := time.ParseDuration(str)
	if err != nil {
		return fallback
	}
	return d
}
