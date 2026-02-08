package cftools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"dayzsmartcf/backend/internal/config"
)

const (
	baseURL   = "https://api.cftools.cloud"
	originURL = "https://auth.cftools.cloud"
)

type Client struct {
	cfg      *config.Config
	client   *http.Client
	acsrf    string
	cookies  []*http.Cookie
}

func New(cfg *config.Config) *Client {
	return &Client{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (c *Client) baseRequest(method, path string, body io.Reader) (*http.Request, error) {
	u, _ := url.JoinPath(baseURL, path)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", originURL)
	req.Header.Set("Referer", originURL+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}

func (c *Client) mergeCookies(newCookies []*http.Cookie) {
	for _, nc := range newCookies {
		found := false
		for i, oc := range c.cookies {
			if strings.EqualFold(oc.Name, nc.Name) {
				c.cookies[i] = nc
				found = true
				break
			}
		}
		if !found {
			c.cookies = append(c.cookies, nc)
		}
	}
}

func (c *Client) fetchStatus() error {
	req, err := c.baseRequest("GET", "/olymp/v1/@me/status", nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	c.mergeCookies(resp.Cookies())
	return nil
}

func (c *Client) fetchPersona() error {
	req, err := c.baseRequest("GET", "/app/v1/@me/persona", nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	c.mergeCookies(resp.Cookies())
	return nil
}

func (c *Client) GetACSRFToken() (string, error) {
	req, err := c.baseRequest("GET", "/olymp/v1/@me/acsrf-token", nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("acsrf token failed: %d %s", resp.StatusCode, string(body))
	}

	c.mergeCookies(resp.Cookies())
	for _, cookie := range c.cookies {
		if strings.EqualFold(cookie.Name, "acsrf") {
			c.acsrf = cookie.Value
			return cookie.Value, nil
		}
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &result); err == nil && result.Token != "" {
		c.acsrf = result.Token
		return result.Token, nil
	}

	return "", fmt.Errorf("acsrf token not found in response")
}

func (c *Client) cookiesFromToken() []*http.Cookie {
	domain := ".cftools.cloud"
	var cookies []*http.Cookie
	if c.cfg.CFtoolsCdnAuth != "" {
		cookies = append(cookies, &http.Cookie{Name: "cdn-auth", Value: c.cfg.CFtoolsCdnAuth, Domain: domain, Path: "/"})
	}
	if c.cfg.CFtoolsSession != "" {
		cookies = append(cookies, &http.Cookie{Name: "session", Value: c.cfg.CFtoolsSession, Domain: domain, Path: "/"})
	}
	if c.cfg.CFtoolsUserInfo != "" {
		cookies = append(cookies, &http.Cookie{Name: "user_info", Value: c.cfg.CFtoolsUserInfo, Domain: domain, Path: "/"})
	}
	if c.cfg.CFtoolsCfClearance != "" {
		cookies = append(cookies, &http.Cookie{Name: "cf_clearance", Value: c.cfg.CFtoolsCfClearance, Domain: domain, Path: "/"})
	}
	if c.cfg.CFtoolsAcsrf != "" {
		c.acsrf = c.cfg.CFtoolsAcsrf
		cookies = append(cookies, &http.Cookie{Name: "acsrf", Value: c.cfg.CFtoolsAcsrf, Domain: domain, Path: "/"})
	}
	return cookies
}

func (c *Client) Login() error {
	// Режим токена: только cookies из .env, без login/acsrf эндпоинтов
	if c.cfg.CFtoolsCdnAuth != "" {
		c.cookies = c.cookiesFromToken()
		if len(c.cookies) == 0 {
			return fmt.Errorf("CFTOOLS_CDN_AUTH set but no valid cookies")
		}
		if c.cfg.CFtoolsAcsrf != "" {
			c.acsrf = c.cfg.CFtoolsAcsrf
		}
		return nil
	}

	if c.cfg.CFtoolsIdentifier == "" || c.cfg.CFtoolsPasswordHash == "" {
		return fmt.Errorf("CFTOOLS_IDENTIFIER and CFTOOLS_PASSWORD_HASH must be set in .env, or use CFTOOLS_CDN_AUTH")
	}

	// Get Cloudflare cookies first (cf_clearance required for API access)
	if len(c.cookies) == 0 {
		cookies, err := fetchCloudflareCookies(c.cfg.CFtoolsHeadless)
		if err != nil {
			return fmt.Errorf("cloudflare cookies: %w", err)
		}
		c.cookies = cookies
	}

	if c.acsrf == "" {
		if _, err := c.GetACSRFToken(); err != nil {
			return fmt.Errorf("get acsrf: %w", err)
		}
	}

	// Status до логина — инициализирует сессию (user_info, session), как на странице auth
	_ = c.fetchStatus()

	payload := map[string]interface{}{
		"acsrf_token": c.acsrf,
		"password":    c.cfg.CFtoolsPasswordHash,
		"identifier":  c.cfg.CFtoolsIdentifier,
		"_v":          2,
		"_i":          c.cfg.CFtoolsIdentifier,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal body: %w", err)
	}

	req, err := c.baseRequest("POST", "/olymp/v1/@me/native-login", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: %d %s", resp.StatusCode, string(body))
	}

	c.mergeCookies(resp.Cookies())

	// После логина — user_info, persona
	_ = c.fetchStatus()
	_ = c.fetchPersona()

	return nil
}

func (c *Client) IsLoggedIn() bool {
	return len(c.cookies) > 0
}

// VerifyAuth проверяет, работают ли текущие cookies — делает реальный запрос к CF API
func (c *Client) VerifyAuth() error {
	_, err := c.GetACSRFToken()
	return err
}

// UpdateAuth устанавливает cookies из значений, обновляемых с фронта (cdn-auth, cf_clearance, session, user_info, acsrf)
func (c *Client) UpdateAuth(cdnAuth, cfClearance, session, userInfo, acsrf string) {
	domain := ".cftools.cloud"
	var cookies []*http.Cookie
	if cdnAuth != "" {
		cookies = append(cookies, &http.Cookie{Name: "cdn-auth", Value: cdnAuth, Domain: domain, Path: "/"})
	}
	if cfClearance != "" {
		cookies = append(cookies, &http.Cookie{Name: "cf_clearance", Value: cfClearance, Domain: domain, Path: "/"})
	}
	if session != "" {
		cookies = append(cookies, &http.Cookie{Name: "session", Value: session, Domain: domain, Path: "/"})
	}
	if userInfo != "" {
		cookies = append(cookies, &http.Cookie{Name: "user_info", Value: userInfo, Domain: domain, Path: "/"})
	}
	if acsrf != "" {
		c.acsrf = acsrf
		cookies = append(cookies, &http.Cookie{Name: "acsrf", Value: acsrf, Domain: domain, Path: "/"})
	}
	if len(cookies) > 0 {
		c.cookies = cookies
	}
}
