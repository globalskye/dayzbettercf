package handlers

import (
	"encoding/json"
	"net/http"

	"dayzsmartcf/backend/internal/player"
)

// DBWipe очищает все данные базы (игроки, группы, история). Только для админа. Users не удаляются.
func DBWipe(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		if err := repo.WipeAllData(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "База очищена"})
	}
}
