package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"dayzsmartcf/backend/internal/auth"
	"dayzsmartcf/backend/internal/config"
	"dayzsmartcf/backend/internal/cftools"
	"dayzsmartcf/backend/internal/db"
	"dayzsmartcf/backend/internal/player"
	"dayzsmartcf/backend/internal/server"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database: %v", err)
		os.Exit(1)
	}
	defer database.Close()
	if err := db.Migrate(database, "migrations"); err != nil {
		log.Fatalf("Migrate: %v", err)
		os.Exit(1)
	}
	log.Println("Database: ok")

	repo := player.NewRepository(database)
	if os.Getenv("SEED_SAMPLE") == "1" {
		if err := repo.SeedSample(); err != nil {
			log.Printf("SeedSample: %v", err)
		}
	}

	authRepo := auth.NewRepo(database)
	exists, _ := authRepo.Exists()
	if !exists && cfg.AdminUser != "" && cfg.AdminPass != "" {
		if err := authRepo.Create(cfg.AdminUser, cfg.AdminPass, auth.RoleAdmin); err != nil {
			log.Printf("Create admin user: %v", err)
		} else {
			log.Printf("Created admin user: %s", cfg.AdminUser)
		}
	}

	cf := cftools.New(cfg)
	log.Println("Logging in to CFtools...")
	if err := cf.Login(); err != nil {
		log.Printf("CFtools login failed (server will start anyway): %v", err)
		log.Println("Update CFtools auth in Settings after logging in.")
	} else {
		log.Println("CFtools: logged in successfully")
	}

	syncSvc := player.NewSyncService(cf, repo)
	tracker := player.NewTracker(cf, repo)
	tracker.Start()

	srv := server.New(cfg, cf, repo, syncSvc, authRepo)
	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, srv.Router()); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
		os.Exit(1)
	}
}
