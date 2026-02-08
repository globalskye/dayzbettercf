package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"dayzsmartcf/backend/internal/config"
	"dayzsmartcf/backend/internal/cftools"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello from DayZ Smart CF API!",
	})
}

func CFtoolsStatus(cftoolsClient *cftools.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"logged_in": cftoolsClient.IsLoggedIn(),
		})
	}
}

func CFtoolsLogin(cftoolsClient *cftools.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := cftoolsClient.Login(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "login_failed",
				"message": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"message": "Logged in to CFtools",
		})
	}
}

// AuthSettingsGet возвращает статус: настроена ли авторизация
func AuthSettingsGet(cftoolsClient *cftools.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"configured": cftoolsClient.IsLoggedIn(),
		})
	}
}

// AuthSettingsCheck проверяет авторизацию — делает реальный запрос к CF API
func AuthSettingsCheck(cftoolsClient *cftools.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := cftoolsClient.VerifyAuth(); err != nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ok":    false,
				"error": err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok": true,
		})
	}
}

type authSettingsBody struct {
	CdnAuth      string `json:"cdn_auth"`
	CfClearance  string `json:"cf_clearance"`
	Session      string `json:"session"`
	UserInfo     string `json:"user_info"`
	Acsrf        string `json:"acsrf"`
}

// AuthSettingsUpdate обновляет cookies CFtools и сохраняет в файл (persist между перезапусками)
func AuthSettingsUpdate(cftoolsClient *cftools.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body authSettingsBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if body.CdnAuth == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "cdn_auth required"})
			return
		}
		cftoolsClient.UpdateAuth(body.CdnAuth, body.CfClearance, body.Session, body.UserInfo, body.Acsrf)

		// Сохраняем в auth.json, чтобы не терялось при перезапуске
		path := config.AuthFilePath(cfg)
		if err := config.SaveAuthFile(path, &config.AuthFile{
			CdnAuth:     body.CdnAuth,
			CfClearance: body.CfClearance,
			Session:     body.Session,
			UserInfo:    body.UserInfo,
			Acsrf:       body.Acsrf,
		}); err != nil {
			// Не фатально — auth уже в памяти
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "warning": "auth saved in memory but file write failed: " + err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

// CFToolsStates — как в old: GET /cftools/states?q=identifier → GlobalQuery + playState по каждому, без записи в БД
func CFToolsStates(cf *cftools.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "missing q"})
			return
		}
		resp, err := cf.GlobalQuery(q)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		maxResults := 30
		if len(resp.Results) > maxResults {
			resp.Results = resp.Results[:maxResults]
		}
		type stateItem struct {
			CftoolsID   string `json:"cftools_id"`
			DisplayName string `json:"display_name"`
			Avatar      string `json:"avatar,omitempty"`
			Online      bool   `json:"online"`
			ServerName  string `json:"server_name,omitempty"`
		}
		states := make([]stateItem, 0, len(resp.Results))
		for _, x := range resp.Results {
			item := stateItem{
				CftoolsID:   x.User.CftoolsID,
				DisplayName: x.User.DisplayName,
				Avatar:      x.User.Avatar,
			}
			data, _ := cf.ProfilePlayState(x.User.CftoolsID)
			if len(data) > 0 {
				var ps struct {
					PlayState struct {
						Online bool `json:"online"`
						Server *struct {
							Name string `json:"name"`
							ID   string `json:"id"`
						} `json:"server"`
					} `json:"playState"`
				}
				if json.Unmarshal(data, &ps) == nil {
					item.Online = ps.PlayState.Online
					if ps.PlayState.Server != nil {
						if ps.PlayState.Server.Name != "" {
							item.ServerName = ps.PlayState.Server.Name
						} else {
							item.ServerName = ps.PlayState.Server.ID
						}
					}
				}
			}
			states = append(states, item)
		}
		// Сначала кто онлайн на сервере
		sort.Slice(states, func(i, j int) bool {
			if states[i].Online != states[j].Online {
				return states[i].Online
			}
			return false
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"states": states,
			"count":  len(states),
		})
	}
}
