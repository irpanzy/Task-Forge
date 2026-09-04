package seed

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/irpanzy/Task-Forge/internal/config"
	"github.com/irpanzy/Task-Forge/internal/model"
	"github.com/irpanzy/Task-Forge/pkg/utils"
)

func SeedAdmin() {
	ctx := context.Background()

	adminEmail := os.Getenv("SEED_ADMIN_EMAIL")
	adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")
	adminRole := os.Getenv("SEED_ADMIN_ROLE")

	if adminEmail == "" || adminPassword == "" {
		log.Println("Peringatan: Lewati seeding admin karena SEED_ADMIN_EMAIL atau SEED_ADMIN_PASSWORD belum diatur di .env")
		return
	}

	if adminRole == "" {
		adminRole = "admin"
	}

	var user model.User
	result := config.DB.Where("email = ?", adminEmail).First(&user)

	if result.RowsAffected > 0 {
		fmt.Printf("Admin (%s) sudah ada di database.\n", adminEmail)
		return
	}

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		log.Fatal("Gagal hash password admin: ", err)
	}

	admin := model.User{
		Email:    adminEmail,
		Password: hashedPassword,
		Role:     model.Role(adminRole),
		Name:     "Administrator",
	}

	if err := config.DB.WithContext(ctx).Create(&admin).Error; err != nil {
		log.Fatalf("Gagal membuat user admin: %v", err)
	}

	fmt.Printf("User admin berhasil dibuat dari .env: %s\n", adminEmail)
}
