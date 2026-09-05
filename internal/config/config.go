package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	AppConfig *Config
)

type Config struct {
	DatabaseURL         string
	Port                string
	CORSOrigin          string
	JWTSecret           string
	JWTExpired          string
	RefreshTokenExpired string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, loading from system env: " + err.Error())
	}
	AppConfig = &Config{
		Port:                getEnv("PORT", "3000"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		CORSOrigin:          getEnv("CORS_ORIGIN", ""),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		JWTExpired:          getEnv("JWT_EXPIRED", "15m"),
		RefreshTokenExpired: getEnv("REFRESH_TOKEN_EXPIRED", "7d"),
	}
}

func getEnv(key string, fallback string) string {
	if value, exist := os.LookupEnv(key); exist && value != "" {
		return value
	}
	return fallback
}

func ConnectDB() {
	if AppConfig == nil || AppConfig.DatabaseURL == "" {
		log.Fatal("DATABASE_URL belum diatur di .env")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(AppConfig.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke Neon DB: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Gagal mengambil instance DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Neon DB tidak merespons: %v", err)
	}

	log.Println("Berhasil terhubung ke Neon PostgreSQL!")

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

}
