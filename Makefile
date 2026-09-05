# ----------------------------------------------------
# TaskForge Automation Shortcuts
# ----------------------------------------------------

# Menjalankan server web Fiber
run:
	@go run cmd/api/main.go

# Menjalankan seeder database
seed:
	@go run cmd/seed/main.go

# Compile binary executable
build:
	@go build -o bin/taskforge.exe ./cmd/api

# Menjalankan seluruh test
test:
	@go test -v ./...

# Sinkronisasi dependensi Go
tidy:
	@go mod tidy

# ----------------------------------------------------
# Database Migrations
# ----------------------------------------------------

# Buat file migrasi baru (contoh: make migrate-create name=create_boards_table)
migrate-create:
	@migrate create -ext sql -dir database/migrations -seq $(name)

# Eksekusi migrasi UP ke Neon DB
migrate-up:
	@go run cmd/migrate/main.go up

# Rollback 1 langkah migrasi DOWN
migrate-down:
	@go run cmd/migrate/main.go down 1

# Cek versi migrasi saat ini
migrate-version:
	@go run cmd/migrate/main.go version

# Force versi migrasi jika dirty state (contoh: make migrate-force version=1)
migrate-force:
	@go run cmd/migrate/main.go force $(version)

.PHONY: run seed build test tidy migrate-create migrate-up migrate-down migrate-version migrate-force
