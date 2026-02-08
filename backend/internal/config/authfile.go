package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const defaultAuthFile = "cftools_auth.json"

type AuthFile struct {
	CdnAuth      string `json:"cdn_auth"`
	CfClearance  string `json:"cf_clearance,omitempty"`
	Session      string `json:"session,omitempty"`
	UserInfo     string `json:"user_info,omitempty"`
	Acsrf        string `json:"acsrf,omitempty"`
}

// AuthFilePath возвращает путь к файлу авторизации (рядом с БД или в cwd)
func AuthFilePath(cfg *Config) string {
	if p := os.Getenv("CFTOOLS_AUTH_FILE"); p != "" {
		return p
	}
	// Если БД в file:path — кладём auth рядом
	if len(cfg.DatabaseURL) > 5 && cfg.DatabaseURL[:5] == "file:" {
		path := cfg.DatabaseURL[5:]
		for i, c := range path {
			if c == '?' || c == '&' {
				path = path[:i]
				break
			}
		}
		dir := filepath.Dir(path)
		if dir != "." {
			return filepath.Join(dir, defaultAuthFile)
		}
	}
	return defaultAuthFile
}

// LoadAuthFile загружает авторизацию из файла
func LoadAuthFile(path string) *AuthFile {
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("auth file read: %v", err)
		}
		return nil
	}
	var a AuthFile
	if err := json.Unmarshal(data, &a); err != nil {
		log.Printf("auth file parse: %v", err)
		return nil
	}
	return &a
}

// SaveAuthFile сохраняет авторизацию в файл
func SaveAuthFile(path string, a *AuthFile) error {
	if a == nil || a.CdnAuth == "" {
		return nil
	}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}
	return os.WriteFile(path, data, 0600)
}
