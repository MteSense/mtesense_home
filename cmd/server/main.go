package main

import (
	"log"
	"net/http"

	"mtesense_home/internal/auth"
	"mtesense_home/internal/config"
	"mtesense_home/internal/db"
	httpapi "mtesense_home/internal/http"
	"mtesense_home/internal/storage"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("database open failed: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := auth.NewService(database, cfg.JWTSecret).EnsureAdmin(cfg.AdminUsername, cfg.AdminPassword); err != nil {
		log.Fatalf("admin bootstrap failed: %v", err)
	}
	if err := storage.NewService(cfg.UploadDir).Ensure(); err != nil {
		log.Fatalf("upload directory bootstrap failed: %v", err)
	}

	app := httpapi.NewRouter(cfg, database)

	log.Printf("MteSense Home listening on http://localhost%s", cfg.Address())
	if err := http.ListenAndServe(cfg.Address(), app); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
