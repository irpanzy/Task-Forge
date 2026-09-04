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
	JWTSecret           string
	JWTExpired          string
	RefreshTokenExpired string
	JWTExpiryMinutes    string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, loading from system env: " + err.Error())
	}
	AppConfig = &Config{
		Port:                getEnv("PORT", "3000"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		JWTExpired:          getEnv("JWT_EXPIRED", "15m"),
		RefreshTokenExpired: getEnv("REFRESH_TOKEN_EXPIRED", "7d"),
		JWTExpiryMinutes:    getEnv("JWT_EXPIRY_MINUTES", "15"),
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

	// 1. Batasi koneksi maksimal yang dibuka aplikasi secara bersamaan
	sqlDB.SetMaxOpenConns(20) 
	// 2. Jumlah koneksi standby/menganggur di pool
	sqlDB.SetMaxIdleConns(5)
	// 3. Waktu maksimal koneksi menganggur sebelum ditutup otomatis (PENTING untuk Neon)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	// 4. Umur maksimal satu koneksi sebelum dibuat ulang dari nol
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

}
