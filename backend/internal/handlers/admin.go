package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dayzsmartcf/backend/internal/auth"
)

func AdminListUsers(repo *auth.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := repo.List()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"users": list})
	}
}

func AdminCreateUser(repo *auth.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Username == "" || body.Password == "" {
			http.Error(w, `{"error":"username and password required"}`, http.StatusBadRequest)
			return
		}
		if body.Role == "" {
			body.Role = auth.RoleViewer
		}
		if body.Role != auth.RoleAdmin && body.Role != auth.RoleEditor && body.Role != auth.RoleViewer {
			http.Error(w, `{"error":"role must be admin, editor or viewer"}`, http.StatusBadRequest)
			return
		}
		if err := repo.Create(body.Username, body.Password, body.Role); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		list, _ := repo.List()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"users": list})
	}
}

func AdminUpdateUser(repo *auth.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch && r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}
		var body struct {
			Role     *string `json:"role"`
			Password *string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		if body.Role != nil {
			role := *body.Role
			if role != auth.RoleAdmin && role != auth.RoleEditor && role != auth.RoleViewer {
				http.Error(w, `{"error":"role must be admin, editor or viewer"}`, http.StatusBadRequest)
				return
			}
			u, _ := repo.GetByID(id)
			if u != nil && u.Role == auth.RoleAdmin {
				admins, _ := repo.CountAdmins()
				if admins <= 1 {
					http.Error(w, `{"error":"cannot change role of the last admin"}`, http.StatusBadRequest)
					return
				}
			}
			if err := repo.UpdateRole(id, role); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
		}
		if body.Password != nil && *body.Password != "" {
			if err := repo.UpdatePassword(id, *body.Password); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
		}
		list, _ := repo.List()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"users": list})
	}
}

func AdminDeleteUser(repo *auth.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}
		u, _ := repo.GetByID(id)
		if u == nil {
			http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
			return
		}
		if u.Role == auth.RoleAdmin {
			admins, _ := repo.CountAdmins()
			if admins <= 1 {
				http.Error(w, `{"error":"cannot delete the last admin"}`, http.StatusBadRequest)
				return
			}
		}
		if err := repo.Delete(id); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func AdminGetRequestLogs(repo *auth.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		logs, err := repo.GetRequestLogs(userID, limit)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"logs": logs})
	}
}
