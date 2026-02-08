package cftools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const appOriginURL = "https://app.cftools.cloud"

func (c *Client) appRequest(method, path string, body io.Reader) (*http.Request, error) {
	u, _ := url.JoinPath(baseURL, path)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", appOriginURL)
	req.Header.Set("Referer", appOriginURL+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}

type GlobalQueryResult struct {
	Identifier string `json:"identifier"`
	User       struct {
		CftoolsID   string `json:"cftools_id"`
		DisplayName string `json:"display_name"`
		Avatar      string `json:"avatar,omitempty"`
	} `json:"user"`
}

type GlobalQueryResponse struct {
	Results []GlobalQueryResult `json:"results"`
	Status  bool               `json:"status"`
}

func (c *Client) GlobalQuery(identifier string) (*GlobalQueryResponse, error) {
	// В режиме токена (cdn-auth) acsrf не используется — передаём пустой
	payload := map[string]string{
		"acsrf_token": c.acsrf,
		"identifier":  identifier,
	}
	bodyBytes, _ := json.Marshal(payload)

	req, err := c.appRequest("POST", "/app/v1/global-query", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	start := time.Now()
	resp, err := c.client.Do(req)
	dur := time.Since(start)
	log.Printf("[CF] POST /app/v1/global-query -> %d (%v)", statusOrErr(resp, err), dur)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.mergeCookies(resp.Cookies())

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("global-query: %d %s", resp.StatusCode, string(data))
	}

	var result GlobalQueryResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ProfileStatus(cftoolsID string) ([]byte, error) {
	return c.profileGet(cftoolsID, "status")
}

func (c *Client) ProfilePlayState(cftoolsID string) ([]byte, error) {
	return c.profileGet(cftoolsID, "playState")
}

func (c *Client) ProfileStructure(cftoolsID string) ([]byte, error) {
	return c.profileGet(cftoolsID, "structure")
}

func (c *Client) ProfileOverview(cftoolsID string) ([]byte, error) {
	return c.profileGet(cftoolsID, "overview")
}

func (c *Client) ProfileActivities(cftoolsID string) ([]byte, error) {
	return c.profileGet(cftoolsID, "activities")
}

func (c *Client) ProfileSteam(cftoolsID string) ([]byte, error) {
	return c.profileGet(cftoolsID, "steam")
}

func (c *Client) ProfileBans(cftoolsID string) ([]byte, error) {
	return c.profileGet(cftoolsID, "bans")
}

func (c *Client) ProfileBattlEyeBanStatus(cftoolsID string) ([]byte, error) {
	return c.profileGet(cftoolsID, "publisher-services/battleye/ban-status")
}

func (c *Client) profileGet(cftoolsID, suffix string) ([]byte, error) {
	path := "/app/v1/profile/" + cftoolsID + "/" + suffix
	req, err := c.appRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	dur := time.Since(start)
	log.Printf("[CF] GET %s -> %d (%v)", path, statusOrErr(resp, err), dur)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.mergeCookies(resp.Cookies())

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("profile %s: %d %s", suffix, resp.StatusCode, string(data))
	}

	return data, nil
}

func statusOrErr(resp *http.Response, err error) int {
	if err != nil {
		return -1
	}
	if resp == nil {
		return -1
	}
	return resp.StatusCode
}
