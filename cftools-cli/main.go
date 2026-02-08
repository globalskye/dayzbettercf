package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	baseURL      = "https://api.cftools.cloud"
	appOriginURL = "https://app.cftools.cloud"
)

var (
	query  = flag.String("query", "", "Поиск: ник или identifier для global-query")
	cftoolsID = flag.String("profile", "", "CFtools ID для запроса профиля (status, overview, playState, structure, activities)")
	action = flag.String("action", "global-query", "Действие: global-query | profile-status | profile-overview | profile-playState | profile-structure | profile-activities")
	token  = flag.String("token", "", "CFTOOLS_CDN_AUTH (или из .env)")
)

func main() {
	flag.Parse()

	_ = godotenv.Load()
	cdnAuth := *token
	if cdnAuth == "" {
		cdnAuth = os.Getenv("CFTOOLS_CDN_AUTH")
	}
	if cdnAuth == "" {
		fmt.Fprintln(os.Stderr, "Укажи CFTOOLS_CDN_AUTH в .env или флаг -token")
		os.Exit(1)
	}

	client := &http.Client{}
	cookies := []*http.Cookie{
		{Name: "cdn-auth", Value: cdnAuth, Domain: ".cftools.cloud", Path: "/"},
	}

	switch *action {
	case "global-query":
		if *query == "" {
			fmt.Fprintln(os.Stderr, "Укажи -query <ник>")
			os.Exit(1)
		}
		res, err := doGlobalQuery(client, cookies, *query)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printJSON(res)
	case "profile-status", "profile-overview", "profile-playState", "profile-structure", "profile-activities":
		id := *cftoolsID
		if id == "" {
			id = *query
		}
		if id == "" {
			fmt.Fprintln(os.Stderr, "Укажи -profile <cftools_id> или -query <id>")
			os.Exit(1)
		}
		suffix := strings.TrimPrefix(*action, "profile-") // profile-status -> status
		res, err := doProfile(client, cookies, id, suffix)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(res))
	default:
		fmt.Fprintf(os.Stderr, "Неизвестный -action: %s\n", *action)
		flag.Usage()
		os.Exit(1)
	}
}

func appRequest(client *http.Client, method, path string, body io.Reader, cookies []*http.Cookie) (*http.Response, error) {
	u, _ := url.JoinPath(baseURL, path)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", appOriginURL)
	req.Header.Set("Referer", appOriginURL+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	return client.Do(req)
}

func doGlobalQuery(client *http.Client, cookies []*http.Cookie, identifier string) (interface{}, error) {
	payload := map[string]string{
		"acsrf_token": "",
		"identifier":  identifier,
	}
	body, _ := json.Marshal(payload)
	resp, err := appRequest(client, "POST", "/app/v1/global-query", bytes.NewReader(body), cookies)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, string(data))
	}
	var out interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func doProfile(client *http.Client, cookies []*http.Cookie, cftoolsID, suffix string) ([]byte, error) {
	path := "/app/v1/profile/" + cftoolsID + "/" + suffix
	resp, err := appRequest(client, "GET", path, nil, cookies)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
