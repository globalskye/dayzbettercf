package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	DatabaseURL string // sqlite:file.db или postgres://...

	JWTSecret   string
	AdminUser   string
	AdminPass   string

	CFtoolsIdentifier   string
	CFtoolsPasswordHash string
	CFtoolsHeadless     bool

	// Токен-режим: логинишься сам, копируешь cookies из DevTools после логина
	CFtoolsCdnAuth    string
	CFtoolsSession    string
	CFtoolsUserInfo   string
	CFtoolsCfClearance string
	CFtoolsAcsrf      string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using env vars")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	headless := os.Getenv("CFTOOLS_HEADLESS") != "false"
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "file:dayzsmartcf.db?_pragma=foreign_keys(1)"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dayzsmartcf-secret-change-in-production"
	}
	cfg := &Config{
		Port:                 port,
		DatabaseURL:          dbURL,
		Environment:          env,
		JWTSecret:            jwtSecret,
		AdminUser:            os.Getenv("ADMIN_USER"),
		AdminPass:            os.Getenv("ADMIN_PASS"),
		CFtoolsIdentifier:    os.Getenv("CFTOOLS_IDENTIFIER"),
		CFtoolsPasswordHash:  os.Getenv("CFTOOLS_PASSWORD_HASH"),
		CFtoolsHeadless:      headless,
		CFtoolsCdnAuth:       os.Getenv("CFTOOLS_CDN_AUTH"),
		CFtoolsSession:       os.Getenv("CFTOOLS_SESSION"),
		CFtoolsUserInfo:      os.Getenv("CFTOOLS_USER_INFO"),
		CFtoolsCfClearance:   os.Getenv("CFTOOLS_CF_CLEARANCE"),
		CFtoolsAcsrf:         os.Getenv("CFTOOLS_ACSRF"),
	}

	// Файл auth.json переопределяет .env — авторизация сохраняется между перезапусками
	if af := LoadAuthFile(AuthFilePath(cfg)); af != nil && af.CdnAuth != "" {
		cfg.CFtoolsCdnAuth = af.CdnAuth
		if af.CfClearance != "" {
			cfg.CFtoolsCfClearance = af.CfClearance
		}
		if af.Session != "" {
			cfg.CFtoolsSession = af.Session
		}
		if af.UserInfo != "" {
			cfg.CFtoolsUserInfo = af.UserInfo
		}
		if af.Acsrf != "" {
			cfg.CFtoolsAcsrf = af.Acsrf
		}
		log.Println("Auth loaded from", AuthFilePath(cfg))
	}

	return cfg
}
