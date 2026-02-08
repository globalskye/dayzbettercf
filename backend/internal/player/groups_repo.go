package player

import (
	"database/sql"
	"time"
)

func (r *Repository) CreateGroup(name string) (*Group, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := r.db.Exec(`INSERT INTO groups (name, updated_at) VALUES (?, ?)`, name, now)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return r.GetGroup(id, "online")
}

func (r *Repository) GetGroup(id int64, membersSort string) (*Group, error) {
	if membersSort == "" {
		membersSort = "online"
	}
	var g Group
	var createdAt, updatedAt string
	err := r.db.QueryRow(`SELECT id, name, created_at, updated_at FROM groups WHERE id = ?`, id).Scan(
		&g.ID, &g.Name, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	g.CreatedAt = parseTimeValue(createdAt)
	g.UpdatedAt = parseTimeValue(updatedAt)
	g.Members, _ = r.GetGroupMembers(id, membersSort)
	return &g, nil
}

func (r *Repository) ListGroups(membersSort string) ([]*Group, error) {
	if membersSort == "" {
		membersSort = "online"
	}
	rows, err := r.db.Query(`SELECT id, name, created_at, updated_at FROM groups ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Group
	for rows.Next() {
		var g Group
		var createdAt, updatedAt string
		_ = rows.Scan(&g.ID, &g.Name, &createdAt, &updatedAt)
		g.CreatedAt = parseTimeValue(createdAt)
		g.UpdatedAt = parseTimeValue(updatedAt)
		g.Members, _ = r.GetGroupMembers(g.ID, membersSort)
		list = append(list, &g)
	}
	return list, nil
}

func (r *Repository) DeleteGroup(id int64) error {
	_, err := r.db.Exec(`DELETE FROM groups WHERE id = ?`, id)
	return err
}

func (r *Repository) GetGroupMembers(groupID int64, sort string) ([]Member, error) {
	order := "ORDER BY p.playtime_sec DESC, p.display_name"
	switch sort {
	case "online":
		order = "ORDER BY p.online DESC, COALESCE(p.last_seen_at,'') DESC, p.display_name"
	case "bans":
		order = "ORDER BY p.bans_count DESC, p.playtime_sec DESC, p.display_name"
	case "name":
		order = "ORDER BY p.display_name"
	default:
		order = "ORDER BY p.playtime_sec DESC, p.display_name"
	}
	rows, err := r.db.Query(`
		SELECT gm.group_id, gm.player_id, p.cftools_id, gm.alias, gm.created_at
		FROM group_members gm
		JOIN players p ON p.id = gm.player_id
		WHERE gm.group_id = ?
		`+order,
		groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Member
	for rows.Next() {
		var m Member
		var createdAt string
		_ = rows.Scan(&m.GroupID, &m.PlayerID, &m.CftoolsID, &m.Alias, &createdAt)
		m.CreatedAt = parseTimeValue(createdAt)
		list = append(list, m)
	}
	// Заполняем Player из БД (без CF), чтобы фронт мог отображать участников
	for i := range list {
		list[i].Player, _ = r.GetByCftoolsID(list[i].CftoolsID)
	}
	return list, nil
}

func (r *Repository) AddGroupMember(groupID int64, playerID int64, alias string) error {
	_, err := r.db.Exec(`INSERT INTO group_members (group_id, player_id, alias) VALUES (?, ?, ?)
		ON CONFLICT(group_id, player_id) DO UPDATE SET alias = excluded.alias`, groupID, playerID, alias)
	return err
}

func (r *Repository) UpdateGroupMemberAlias(groupID int64, playerID int64, alias string) error {
	_, err := r.db.Exec(`UPDATE group_members SET alias = ? WHERE group_id = ? AND player_id = ?`, alias, groupID, playerID)
	return err
}

func (r *Repository) RemoveGroupMember(groupID int64, playerID int64) error {
	_, err := r.db.Exec(`DELETE FROM group_members WHERE group_id = ? AND player_id = ?`, groupID, playerID)
	return err
}

func (r *Repository) GetByID(id int64) (*Player, error) {
	var cftoolsID string
	err := r.db.QueryRow(`SELECT cftools_id FROM players WHERE id = ?`, id).Scan(&cftoolsID)
	if err != nil {
		return nil, nil
	}
	return r.GetByCftoolsID(cftoolsID)
}
