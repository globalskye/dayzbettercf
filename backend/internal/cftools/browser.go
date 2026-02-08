package cftools

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

// fetchCloudflareCookies opens auth.cftools.cloud in headless browser to pass
// Cloudflare protection and returns cookies (including cf_clearance) for API requests.
func fetchCloudflareCookies(headless bool) ([]*http.Cookie, error) {
	u := launcher.New().
		Leakless(false). // отключаем leakless.exe — Windows Defender может блокировать
		Headless(headless).
		Set("no-sandbox", "").
		Set("disable-dev-shm-usage", "").
		MustLaunch()

	browser := rod.New().
		ControlURL(u).
		Timeout(60 * time.Second).
		MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)
	page = page.Timeout(45 * time.Second)

	if err := page.Navigate(originURL); err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}

	// Wait for Cloudflare challenge to complete ("Checking your browser" ~5–15 sec)
	_ = page.WaitLoad()
	time.Sleep(10 * time.Second)

	cookies, err := proto.NetworkGetAllCookies{}.Call(browser)
	if err != nil {
		return nil, fmt.Errorf("get cookies: %w", err)
	}

	var result []*http.Cookie
	for _, c := range cookies.Cookies {
		// Include cftools.cloud domain cookies (cf_clearance, session, etc.)
		domain := c.Domain
		if domain == "" {
			continue
		}
		if !strings.Contains(domain, "cftools.cloud") {
			continue
		}
		result = append(result, &http.Cookie{
			Name:   c.Name,
			Value:  c.Value,
			Domain: domain,
			Path:   c.Path,
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no cookies received, Cloudflare may still be challenging")
	}

	return result, nil
}
