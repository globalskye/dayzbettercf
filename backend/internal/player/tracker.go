package player

import (
	"encoding/json"
	"log"
	"time"

	"dayzsmartcf/backend/internal/cftools"
)

const (
	playStateInterval = 10 * time.Second
	profileInterval   = 5 * time.Minute
)

type Tracker struct {
	cf     *cftools.Client
	repo   *Repository
	stopCh chan struct{}
}

func NewTracker(cf *cftools.Client, repo *Repository) *Tracker {
	return &Tracker{
		cf:     cf,
		repo:   repo,
		stopCh: make(chan struct{}),
	}
}

func (t *Tracker) Start() {
	go t.loopPlayState()
	go t.loopProfile()
	log.Printf("Tracker started: playState every %v, profile/nick every %v", playStateInterval, profileInterval)
}

func (t *Tracker) Stop() {
	close(t.stopCh)
}

func (t *Tracker) loopPlayState() {
	tick := time.NewTicker(playStateInterval)
	defer tick.Stop()
	time.Sleep(5 * time.Second)
	t.pollPlayState()
	for {
		select {
		case <-t.stopCh:
			return
		case <-tick.C:
			t.pollPlayState()
		}
	}
}

func (t *Tracker) loopProfile() {
	tick := time.NewTicker(profileInterval)
	defer tick.Stop()
	time.Sleep(30 * time.Second)
	t.pollProfile()
	for {
		select {
		case <-t.stopCh:
			return
		case <-tick.C:
			t.pollProfile()
		}
	}
}

func (t *Tracker) pollPlayState() {
	list, err := t.repo.ListTracked("")
	if err != nil {
		log.Printf("tracker playState: list: %v", err)
		return
	}
	for _, p := range list {
		t.updatePlayState(p.ID, p.CftoolsID, p.DisplayName)
		time.Sleep(300 * time.Millisecond)
	}
}

func (t *Tracker) pollProfile() {
	list, err := t.repo.ListTracked("")
	if err != nil {
		log.Printf("tracker profile: list: %v", err)
		return
	}
	for _, p := range list {
		t.updateProfile(p.ID, p.CftoolsID)
		time.Sleep(500 * time.Millisecond)
	}
}

func (t *Tracker) updatePlayState(playerID int64, cftoolsID, displayName string) {
	data, err := t.cf.ProfilePlayState(cftoolsID)
	if err != nil {
		log.Printf("tracker playState %s: %v", cftoolsID, err)
		return
	}
	var online bool
	var serverName string
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
			online = ps.PlayState.Online
			if ps.PlayState.Server != nil && ps.PlayState.Server.Name != "" {
				serverName = ps.PlayState.Server.Name
			} else if ps.PlayState.Server != nil && ps.PlayState.Server.ID != "" {
				serverName = ps.PlayState.Server.ID
			}
		}
	}
	_ = t.repo.UpdatePlayerOnlineStatus(playerID, online, serverName)
	last, _ := t.repo.GetLastPlayerHistory(playerID)
	// Записываем только при смене состояния (сравнение с последней записью), не каждые N секунд
	stateChanged := last == nil || last.Online != online
	if !stateChanged {
		return
	}
	now := time.Now().UTC()
	var sessionDurationSec, offlineDurationSec int64
	if last != nil {
		lastTs, parseErr := time.Parse(time.RFC3339, last.Ts)
		if parseErr != nil {
			lastTs, parseErr = time.Parse(time.RFC3339Nano, last.Ts)
		}
		if parseErr == nil {
			sec := int64(now.Sub(lastTs).Seconds())
			if online {
				offlineDurationSec = sec
			} else {
				sessionDurationSec = sec
			}
		}
	}
	_ = t.repo.AppendPlayerHistory(playerID, online, serverName, 0, 0, displayName, sessionDurationSec, offlineDurationSec)
}

func (t *Tracker) updateProfile(playerID int64, cftoolsID string) {
	statusData, _ := t.cf.ProfileStatus(cftoolsID)
	overviewData, _ := t.cf.ProfileOverview(cftoolsID)

	displayName := ""
	if len(statusData) > 0 {
		var st struct {
			Profile struct {
				DisplayName string `json:"display_name"`
			} `json:"profile"`
		}
		if json.Unmarshal(statusData, &st) == nil && st.Profile.DisplayName != "" {
			displayName = st.Profile.DisplayName
		}
	}
	if displayName != "" {
		_ = t.repo.UpdatePlayerDisplayName(playerID, displayName)
	}

	nicknames := []string{}
	if displayName != "" && !isCftoolsIDLike(displayName) && displayName != cftoolsID {
		nicknames = append(nicknames, displayName)
	}
	if len(overviewData) > 0 {
		var ov struct {
			Omega struct {
				Aliases []string `json:"aliases"`
			} `json:"omega"`
		}
		if json.Unmarshal(overviewData, &ov) == nil {
			for _, a := range ov.Omega.Aliases {
				if a != "" && !isCftoolsIDLike(a) && a != cftoolsID {
					nicknames = append(nicknames, a)
				}
			}
		}
	}
	for _, nick := range nicknames {
		if nick != "" {
			_ = t.repo.UpsertNickname(playerID, nick, "tracker")
		}
	}
}
