package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dayzsmartcf/backend/internal/cftools"
	"dayzsmartcf/backend/internal/player"
)

func PlayersSyncBatch(sync *player.SyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			CftoolsIDs []string `json:"cftools_ids"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if len(body.CftoolsIDs) == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"players": []interface{}{}, "count": 0})
			return
		}
		light := r.URL.Query().Get("light") != "0"
		players, err := sync.SyncBatch(body.CftoolsIDs, light)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"players": players,
			"count":   len(players),
		})
	}
}

func PlayersSearch(sync *player.SyncService, repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			http.Error(w, `{"error":"missing q"}`, http.StatusBadRequest)
			return
		}
		light := r.URL.Query().Get("light") == "1" // только status+playState+overview — меньше запросов к CF

		players, err := sync.SearchAndSync(q, light)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"players": players,
			"count":   len(players),
		})
	}
}

func PlayersList(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := player.ListOptions{
			Limit:      parseInt(r.URL.Query().Get("limit"), 50, 200),
			Offset:     parseInt(r.URL.Query().Get("offset"), 0, 10000),
			OnlyOnline: r.URL.Query().Get("online") == "1",
			OnlyBanned: r.URL.Query().Get("banned") == "1",
			Sort:       r.URL.Query().Get("sort"),
		}
		if opts.Sort == "" {
			opts.Sort = "online"
		}

		players, err := repo.ListAll(opts)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		total, _ := repo.Count(&opts)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"players": players,
			"total":   total,
		})
	}
}

func PlayersGet(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cftoolsID := chi.URLParam(r, "id")
		if cftoolsID == "" {
			http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
			return
		}

		p, err := repo.GetByCftoolsID(cftoolsID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if p == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "player not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

func PlayersSyncOne(sync *player.SyncService, repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cftoolsID := chi.URLParam(r, "id")
		if cftoolsID == "" {
			http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
			return
		}

		light := r.URL.Query().Get("light") == "1"
		p, err := sync.SyncPlayer(cftoolsID, light)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

func PlayersSearchCFtools(cf *cftools.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			http.Error(w, `{"error":"missing q"}`, http.StatusBadRequest)
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
		type item struct {
			CftoolsID   string `json:"cftools_id"`
			DisplayName string `json:"display_name"`
			Avatar      string `json:"avatar,omitempty"`
			Identifier  string `json:"identifier,omitempty"`
		}
		results := make([]item, 0, len(resp.Results))
		for _, x := range resp.Results {
			results = append(results, item{
				CftoolsID:   x.User.CftoolsID,
				DisplayName: x.User.DisplayName,
				Avatar:      x.User.Avatar,
				Identifier:  x.Identifier,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": results,
			"count":   len(results),
		})
	}
}

func parseInt(s string, def, max int) int {
	n, _ := strconv.Atoi(s)
	if n <= 0 {
		return def
	}
	if n > max {
		return max
	}
	return n
}

func PlayersSearchLocal(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			http.Error(w, `{"error":"missing q"}`, http.StatusBadRequest)
			return
		}

		opts := &player.ListOptions{
			Limit:      parseInt(r.URL.Query().Get("limit"), 5000, 10000),
			OnlyOnline: r.URL.Query().Get("online") == "1",
			OnlyBanned: r.URL.Query().Get("banned") == "1",
			Sort:       r.URL.Query().Get("sort"),
		}
		if opts.Sort == "" {
			opts.Sort = "online"
		}

		players, err := repo.SearchByNickname(q, opts.Limit, opts)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"players": players,
			"count":   len(players),
		})
	}
}
