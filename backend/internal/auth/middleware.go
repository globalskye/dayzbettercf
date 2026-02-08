package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// Path prefixes we log (players, groups, tracked)
var logRequestPrefixes = []string{"/api/v1/players", "/api/v1/groups", "/api/v1/tracked"}

type contextKey string

const UserContextKey contextKey = "user"

func RequireAuth(secret string, userRepo *Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				writeAuthError(w, http.StatusUnauthorized, "missing token")
				return
			}
			claims, err := ParseToken(token, secret)
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			user, err := userRepo.GetByID(claims.UserID)
			if err != nil || user == nil {
				writeAuthError(w, http.StatusUnauthorized, "user not found")
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user == nil {
				writeAuthError(w, http.StatusUnauthorized, "not authenticated")
				return
			}
			if !user.HasRole(roles...) {
				writeAuthError(w, http.StatusForbidden, "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserFromContext(ctx context.Context) *User {
	u, _ := ctx.Value(UserContextKey).(*User)
	return u
}

// LogRequests middleware logs requests to players/groups/tracked API per user (for admin audit).
func LogRequests(repo *Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user != nil {
				path := r.URL.Path
				for _, p := range logRequestPrefixes {
					if strings.HasPrefix(path, p) {
						_ = repo.LogRequest(user.ID, r.Method, path)
						break
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func extractToken(r *http.Request) string {
	if h := r.Header.Get("Authorization"); h != "" {
		if parts := strings.SplitN(h, " ", 2); len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return parts[1]
		}
	}
	return r.URL.Query().Get("token")
}

func writeAuthError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
