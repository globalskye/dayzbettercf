package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dayzsmartcf/backend/internal/player"
)

func TrackedList(repo *player.Repository, _ *player.SyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortParam := r.URL.Query().Get("sort")
		if sortParam == "" {
			sortParam = "online"
		}
		list, err := repo.ListTracked("")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		sortTracked(list, sortParam)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"players": list})
	}
}

func sortTracked(list []*player.Player, by string) {
	switch by {
	case "playtime":
		sort.Slice(list, func(i, j int) bool { return list[i].PlaytimeSec > list[j].PlaytimeSec })
	case "bans":
		sort.Slice(list, func(i, j int) bool {
			if list[i].BansCount != list[j].BansCount {
				return list[i].BansCount > list[j].BansCount
			}
			return list[i].PlaytimeSec > list[j].PlaytimeSec
		})
	default: // "online"
		sort.Slice(list, func(i, j int) bool {
			if list[i].Online != list[j].Online {
				return list[i].Online
			}
			return list[i].DisplayName < list[j].DisplayName
		})
	}
}

func TrackedAdd(repo *player.Repository, syncSvc *player.SyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cftoolsID := chi.URLParam(r, "cftoolsId")
		if cftoolsID == "" {
			http.Error(w, `{"error":"missing cftoolsId"}`, http.StatusBadRequest)
			return
		}
		p, err := repo.GetByCftoolsID(cftoolsID)
		if err != nil || p == nil {
			// Игрока нет в базе — подтягиваем из CF и сохраняем, затем добавляем в отслеживание
			if syncSvc != nil {
				p, err = syncSvc.SyncPlayer(cftoolsID, true)
			}
			if err != nil || p == nil {
				http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
				return
			}
		}
		if err := repo.AddTracked(p.ID); err != nil {
			w.Header().Set("Content-Type", "application/json")
			if err == player.ErrTrackedLimit {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "tracked players limit reached (max 10)"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		list, _ := repo.ListTracked("online")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"players": list})
	}
}

func TrackedRemove(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cftoolsID := chi.URLParam(r, "cftoolsId")
		if cftoolsID == "" {
			http.Error(w, `{"error":"missing cftoolsId"}`, http.StatusBadRequest)
			return
		}
		p, _ := repo.GetByCftoolsID(cftoolsID)
		if p == nil {
			http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
			return
		}
		if err := repo.RemoveTracked(p.ID); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func TrackedHistory(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cftoolsID := chi.URLParam(r, "cftoolsId")
		if cftoolsID == "" {
			http.Error(w, `{"error":"missing cftoolsId"}`, http.StatusBadRequest)
			return
		}
		p, _ := repo.GetByCftoolsID(cftoolsID)
		if p == nil {
			http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
			return
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 {
			limit = 5000
		}
		if limit > 10000 {
			limit = 10000
		}
		history, err := repo.GetPlayerHistory(p.ID, limit)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"player":  p,
			"history": history,
		})
	}
}

func PlayerHistory(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cftoolsID := chi.URLParam(r, "id")
		if cftoolsID == "" {
			http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
			return
		}
		p, _ := repo.GetByCftoolsID(cftoolsID)
		if p == nil {
			http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
			return
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 {
			limit = 200
		}
		history, err := repo.GetPlayerHistory(p.ID, limit)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"history": history,
		})
	}
}
