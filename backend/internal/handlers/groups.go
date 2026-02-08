package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"

	"dayzsmartcf/backend/internal/player"
)

func GroupsList(repo *player.Repository, syncSvc *player.SyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortParam := r.URL.Query().Get("sort")
		if sortParam == "" {
			sortParam = "online"
		}
		enrich := r.URL.Query().Get("enrich") != "0"
		list, err := repo.ListGroups(sortParam)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if enrich {
			for _, g := range list {
				enrichMembersFromCF(syncSvc, &g.Members, sortParam)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"groups": list})
	}
}

func GroupsCreate(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := chi.URLParam(r, "name")
		if name == "" {
			http.Error(w, `{"error":"missing name"}`, http.StatusBadRequest)
			return
		}
		g, err := repo.CreateGroup(name)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)
	}
}

func GroupsGet(repo *player.Repository, syncSvc *player.SyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}
		sortParam := r.URL.Query().Get("sort")
		if sortParam == "" {
			sortParam = "online"
		}
		g, err := repo.GetGroup(id, sortParam)
		if err != nil || g == nil {
			http.Error(w, `{"error":"group not found"}`, http.StatusNotFound)
			return
		}
		// Не подтягиваем данные из CF при открытии группы — никнеймы и прочее обновляются только по кнопке «Обновить данные».
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)
	}
}

func GroupsDelete(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}
		if err := repo.DeleteGroup(id); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func GroupsAddMember(repo *player.Repository, syncSvc *player.SyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		groupID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		cftoolsID := chi.URLParam(r, "cftoolsId")
		alias := r.URL.Query().Get("alias")
		if cftoolsID == "" {
			http.Error(w, `{"error":"missing cftoolsId"}`, http.StatusBadRequest)
			return
		}
		p, _ := repo.GetByCftoolsID(cftoolsID)
		if p == nil {
			http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
			return
		}
		if err := repo.AddGroupMember(groupID, p.ID, alias); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		g, _ := repo.GetGroup(groupID, "online")
		if g != nil {
			enrichMembersFromCF(syncSvc, &g.Members, "online")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)
	}
}

func GroupsUpdateMemberAlias(repo *player.Repository, syncSvc *player.SyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch && r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		groupID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		cftoolsID := chi.URLParam(r, "cftoolsId")
		if cftoolsID == "" {
			http.Error(w, `{"error":"missing cftoolsId"}`, http.StatusBadRequest)
			return
		}
		var body struct {
			Alias string `json:"alias"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
			return
		}
		p, _ := repo.GetByCftoolsID(cftoolsID)
		if p == nil {
			http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
			return
		}
		if err := repo.UpdateGroupMemberAlias(groupID, p.ID, body.Alias); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		g, _ := repo.GetGroup(groupID, "online")
		if g != nil {
			enrichMembersFromCF(syncSvc, &g.Members, "online")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)
	}
}

func GroupsRemoveMember(repo *player.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
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
		if err := repo.RemoveGroupMember(groupID, p.ID); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// enrichMembersFromCF подтягивает данные игроков из CF и сортирует участников.
func enrichMembersFromCF(syncSvc *player.SyncService, members *[]player.Member, sortParam string) {
	if members == nil {
		return
	}
	for i := range *members {
		p, _ := syncSvc.FetchPlayerFromCF((*members)[i].CftoolsID)
		(*members)[i].Player = p
	}
	sortMembers(members, sortParam)
}

func sortMembers(members *[]player.Member, by string) {
	if members == nil {
		return
	}
	m := *members
	switch by {
	case "playtime":
		sort.Slice(m, func(i, j int) bool {
			pi, pj := m[i].Player, m[j].Player
			if pi == nil || pj == nil {
				return pi != nil
			}
			return pi.PlaytimeSec > pj.PlaytimeSec
		})
	case "bans":
		sort.Slice(m, func(i, j int) bool {
			pi, pj := m[i].Player, m[j].Player
			if pi == nil || pj == nil {
				return pi != nil
			}
			if pi.BansCount != pj.BansCount {
				return pi.BansCount > pj.BansCount
			}
			return pi.PlaytimeSec > pj.PlaytimeSec
		})
	case "name":
		sort.Slice(m, func(i, j int) bool {
			pi, pj := m[i].Player, m[j].Player
			ni, nj := m[i].Alias, m[j].Alias
			if ni == "" && pi != nil {
				ni = pi.DisplayName
			}
			if nj == "" && pj != nil {
				nj = pj.DisplayName
			}
			return ni < nj
		})
	default: // "online"
		sort.Slice(m, func(i, j int) bool {
			pi, pj := m[i].Player, m[j].Player
			if pi == nil || pj == nil {
				return pi != nil
			}
			if pi.Online != pj.Online {
				return pi.Online
			}
			return pi.DisplayName < pj.DisplayName
		})
	}
}
