package player

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"dayzsmartcf/backend/internal/cftools"
)

// isCftoolsIDLike возвращает true, если s похож на CFTools ID (24 hex-символа, опционально с суффиксом "+").
// CFTools иногда отдаёт ID в omega.aliases — такие «ники» не показываем.
func isCftoolsIDLike(s string) bool {
	s = strings.TrimSuffix(strings.TrimSpace(s), "+")
	if len(s) != 24 {
		return false
	}
	for _, c := range s {
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
			continue
		}
		return false
	}
	return true
}

type SyncService struct {
	cf   *cftools.Client
	repo *Repository
}

func NewSyncService(cf *cftools.Client, repo *Repository) *SyncService {
	return &SyncService{cf: cf, repo: repo}
}

const maxSearchResults = 30

func (s *SyncService) SearchAndSync(identifier string, light bool) ([]*Player, error) {
	resp, err := s.cf.GlobalQuery(identifier)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var saved []*Player

	for _, r := range resp.Results {
		if len(saved) >= maxSearchResults {
			break
		}
		cftoolsID := r.User.CftoolsID
		if cftoolsID == "" || seen[cftoolsID] {
			continue
		}
		seen[cftoolsID] = true

		p, err := s.fetchAndSavePlayer(cftoolsID, r.User.DisplayName, r.User.Avatar, r.Identifier, light)
		if err != nil {
			log.Printf("sync player %s: %v", cftoolsID, err)
			continue
		}
		saved = append(saved, p)
	}

	return saved, nil
}

func (s *SyncService) SyncPlayer(cftoolsID string, light bool) (*Player, error) {
	return s.fetchAndSavePlayer(cftoolsID, "", "", "", light)
}

// SyncBatch syncs multiple players to DB by cftools_ids (from CF search results).
func (s *SyncService) SyncBatch(cftoolsIDs []string, light bool) ([]*Player, error) {
	var saved []*Player
	for _, id := range cftoolsIDs {
		if id == "" {
			continue
		}
		p, err := s.fetchAndSavePlayer(id, "", "", "", light)
		if err != nil {
			log.Printf("sync batch %s: %v", id, err)
			continue
		}
		saved = append(saved, p)
	}
	return saved, nil
}

// FetchPlayerFromCF запрашивает актуальные данные игрока из CFtools API без записи в БД.
// Используется для групп и отслеживания — всегда свежие данные из CF.
func (s *SyncService) FetchPlayerFromCF(cftoolsID string) (*Player, error) {
	statusData, _ := s.cf.ProfileStatus(cftoolsID)
	playStateData, _ := s.cf.ProfilePlayState(cftoolsID)
	overviewData, _ := s.cf.ProfileOverview(cftoolsID)
	structureData, _ := s.cf.ProfileStructure(cftoolsID)
	p := buildPlayerFromCFData(cftoolsID, statusData, playStateData, overviewData, structureData)
	p.UpdatedAt = time.Now().UTC()
	return p, nil
}

func (s *SyncService) fetchAndSavePlayer(cftoolsID, displayName, avatar, searchIdentifier string, light bool) (*Player, error) {
	statusData, _ := s.cf.ProfileStatus(cftoolsID)
	playStateData, _ := s.cf.ProfilePlayState(cftoolsID)
	overviewData, _ := s.cf.ProfileOverview(cftoolsID)
	structureData, _ := s.cf.ProfileStructure(cftoolsID)
	var steamData, bansData, battleyeData []byte
	if !light {
		steamData, _ = s.cf.ProfileSteam(cftoolsID)
		bansData, _ = s.cf.ProfileBans(cftoolsID)
		battleyeData, _ = s.cf.ProfileBattlEyeBanStatus(cftoolsID)
		_, _ = s.cf.ProfileActivities(cftoolsID)
	}

	p := buildPlayerFromCFData(cftoolsID, statusData, playStateData, overviewData, structureData)
	if p == nil {
		p = &Player{CftoolsID: cftoolsID, DisplayName: displayName, Avatar: avatar}
	}
	p.RawStatus = string(statusData)
	p.RawPlayState = string(playStateData)
	p.RawStructure = string(structureData)
	p.RawOverview = string(overviewData)
	p.RawBans = string(bansData)
	p.RawBattlEye = string(battleyeData)
	now := time.Now().UTC()
	p.LastSeenAt = &now

	// Parse Steam (only in full sync)
	if len(steamData) > 0 {
		var steam struct {
			SteamID string `json:"steam64"`
			Profile struct {
				Avatar      string `json:"avatar"`
				Avatarfull  string `json:"avatarfull"`
				PersonaName string `json:"personaname"`
			} `json:"profile"`
			Bans struct {
				NumberOfGameBans int `json:"NumberOfGameBans"`
				NumberOfVACBans  int `json:"NumberOfVACBans"`
			} `json:"bans"`
		}
		if json.Unmarshal(steamData, &steam) == nil {
			p.Steam64 = steam.SteamID
			if steam.Profile.Avatarfull != "" {
				p.SteamAvatar = steam.Profile.Avatarfull
			} else {
				p.SteamAvatar = steam.Profile.Avatar
			}
			p.SteamPersona = steam.Profile.PersonaName
			p.SteamVacBans = steam.Bans.NumberOfVACBans
			p.SteamGameBans = steam.Bans.NumberOfGameBans
		}
	}

	// Upsert
	playerID, err := s.repo.UpsertPlayer(p)
	if err != nil {
		return nil, err
	}

	// Лог обновления в БД
	_ = s.repo.LogSync(playerID, p.CftoolsID, p.DisplayName)

	// Save nicknames (не сохраняем CFTools ID как ник — API иногда отдаёт их в aliases)
	nicknames := make(map[string]string)
	if p.DisplayName != "" && !isCftoolsIDLike(p.DisplayName) {
		nicknames[p.DisplayName] = "display_name"
	}
	if searchIdentifier != "" && !isCftoolsIDLike(searchIdentifier) && searchIdentifier != p.CftoolsID {
		nicknames[searchIdentifier] = "search"
	}
	if len(overviewData) > 0 {
		var ov struct {
			Omega struct {
				Aliases []string `json:"aliases"`
			} `json:"omega"`
		}
		if json.Unmarshal(overviewData, &ov) == nil {
			for _, a := range ov.Omega.Aliases {
				if a != "" && !isCftoolsIDLike(a) && a != p.CftoolsID {
					nicknames[a] = "alias"
				}
			}
		}
	}
	for nick, src := range nicknames {
		_ = s.repo.UpsertNickname(playerID, nick, src)
	}

	// Save links with confirmed/trusted from overview
	_ = s.repo.DeletePlayerLinks(playerID)
	if len(overviewData) > 0 {
		var ov struct {
			AlternateAccounts struct {
				Links []struct {
					CftoolsID string `json:"cftools_id"`
					Confirmed bool   `json:"confirmed"`
					Trusted   bool   `json:"trusted"`
				} `json:"links"`
			} `json:"alternate_accounts"`
		}
		if json.Unmarshal(overviewData, &ov) == nil {
			for _, link := range ov.AlternateAccounts.Links {
				_ = s.repo.UpsertPlayerLink(playerID, link.CftoolsID, link.Confirmed, link.Trusted)
			}
		}
	}

	// Save servers
	_ = s.repo.DeletePlayerServers(playerID)
	if len(structureData) > 0 {
		var st struct {
			Servers []struct {
				ID         string `json:"id"`
				Identifier string `json:"identifier"`
				Game       int    `json:"game"`
			} `json:"servers"`
		}
		if json.Unmarshal(structureData, &st) == nil {
			for _, sv := range st.Servers {
				_ = s.repo.UpsertPlayerServer(playerID, sv.ID, sv.Identifier, sv.Game)
			}
		}
	}

	return s.repo.GetByCftoolsID(cftoolsID)
}

// buildPlayerFromCFData собирает Player из ответов CF API (status, playState, overview, structure) без БД.
func buildPlayerFromCFData(cftoolsID string, statusData, playStateData, overviewData, structureData []byte) *Player {
	p := &Player{CftoolsID: cftoolsID}
	if len(statusData) > 0 {
		var st struct {
			Account struct {
				IsBot  bool `json:"is_bot"`
				Status int  `json:"status"`
			} `json:"account"`
			Profile struct {
				DisplayName string `json:"display_name"`
				Avatar      string `json:"avatar"`
			} `json:"profile"`
		}
		if json.Unmarshal(statusData, &st) == nil {
			p.IsBot = st.Account.IsBot
			p.AccountStatus = st.Account.Status
			p.DisplayName = st.Profile.DisplayName
			p.Avatar = st.Profile.Avatar
		}
	}
	if len(playStateData) > 0 {
		var ps struct {
			PlayState struct {
				Online bool `json:"online"`
				Server *struct {
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"server"`
			} `json:"playState"`
		}
		if json.Unmarshal(playStateData, &ps) == nil {
			p.Online = ps.PlayState.Online
			if ps.PlayState.Server != nil && ps.PlayState.Server.Name != "" {
				p.LastServerIdentifier = ps.PlayState.Server.Name
			} else if ps.PlayState.Server != nil && ps.PlayState.Server.ID != "" {
				p.LastServerIdentifier = ps.PlayState.Server.ID
			}
		}
	}
	if len(structureData) > 0 {
		var st struct {
			Bans struct {
				Count int `json:"count"`
			} `json:"bans"`
			Servers []struct {
				ID string `json:"id"`
			} `json:"servers"`
		}
		if json.Unmarshal(structureData, &st) == nil {
			p.BansCount = st.Bans.Count
			for _, sv := range st.Servers {
				p.ServerIDs = append(p.ServerIDs, sv.ID)
			}
		}
	}
	if len(overviewData) > 0 {
		var ov struct {
			AlternateAccounts struct {
				TotalCount int `json:"total_count"`
				Links      []struct {
					CftoolsID string `json:"cftools_id"`
				} `json:"links"`
			} `json:"alternate_accounts"`
			Omega struct {
				Playtime  int64    `json:"playtime"`
				Sessions  int      `json:"sessions"`
				UpdatedAt string   `json:"updated_at"`
				Aliases   []string `json:"aliases"`
			} `json:"omega"`
		}
		if json.Unmarshal(overviewData, &ov) == nil {
			p.LinkedAccountsCount = ov.AlternateAccounts.TotalCount
			p.PlaytimeSec = ov.Omega.Playtime
			p.SessionsCount = ov.Omega.Sessions
			p.LastActivityAt = parseTime(ov.Omega.UpdatedAt)
			for _, link := range ov.AlternateAccounts.Links {
				p.LinkedCftoolsIDs = append(p.LinkedCftoolsIDs, link.CftoolsID)
			}
			// Ники из CF (исключаем CFTools ID — API иногда отдаёт их в aliases)
			if len(ov.Omega.Aliases) > 0 {
				p.Nicknames = make([]string, 0, len(ov.Omega.Aliases)+1)
				for _, a := range ov.Omega.Aliases {
					if a != "" && !isCftoolsIDLike(a) && a != cftoolsID {
						p.Nicknames = append(p.Nicknames, a)
					}
				}
			}
			if p.DisplayName != "" {
				hasDisplay := false
				for _, a := range p.Nicknames {
					if a == p.DisplayName {
						hasDisplay = true
						break
					}
				}
				if !hasDisplay {
					p.Nicknames = append(p.Nicknames, p.DisplayName)
				}
			}
		}
	}
	return p
}
