package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	acsrfURL = "https://api.cftools.cloud/olymp/v1/@me/acsrf-token"
	authURL  = "https://api.cftools.cloud/olymp/v1/@me/native-login"
)

func GetAuthCookie(CFUsername, CFPasswordHash string) ([]*http.Cookie, error) {
	client := &http.Client{}

	// 1. GET acsrf token
	req, err := http.NewRequest("GET", acsrfURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create acsrf request: %w", err)
	}
	setAuthHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("acsrf request: %w", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	acsrfCookies := resp.Cookies()
	var acsrf string
	for _, c := range acsrfCookies {
		if c.Name == "acsrf" {
			acsrf = c.Value
			break
		}
	}
	if acsrf == "" {
		return nil, fmt.Errorf("acsrf cookie not found in response")
	}

	// 2. POST login
	loginData := map[string]string{
		"acsrf_token": acsrf,
		"identifier":  CFUsername,
		"password":    CFPasswordHash,
	}
	body, err := json.Marshal(loginData)
	if err != nil {
		return nil, fmt.Errorf("marshal login: %w", err)
	}

	loginReq, err := http.NewRequest("POST", authURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create login request: %w", err)
	}
	setAuthHeaders(loginReq)
	loginReq.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	for _, c := range acsrfCookies {
		loginReq.AddCookie(c)
	}

	loginResp, err := client.Do(loginReq)
	if err != nil {
		return nil, fmt.Errorf("login request: %w", err)
	}
	defer loginResp.Body.Close()

	loginCookies := loginResp.Cookies()
	// объединяем: сначала cookies от логина, потом acsrf
	out := make([]*http.Cookie, 0, len(loginCookies)+len(acsrfCookies))
	out = append(out, loginCookies...)
	out = append(out, acsrfCookies...)

	for _, c := range out {
		log.Println("cookie:", c.Name)
	}
	return out, nil
}

func setAuthHeaders(req *http.Request) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://auth.cftools.cloud")
	req.Header.Set("Referer", "https://auth.cftools.cloud/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
}
