package auth

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	passwordHash string
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetByUsername(username string) (*User, error) {
	var u User
	var createdAt string
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, role, created_at FROM users WHERE username = ?`,
		username,
	).Scan(&u.ID, &u.Username, &u.passwordHash, &u.Role, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &u, nil
}

func (r *Repo) GetByID(id int64) (*User, error) {
	var u User
	var createdAt string
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, role, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Username, &u.passwordHash, &u.Role, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &u, nil
}

func (r *Repo) Create(username, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)`,
		username, string(hash), role,
	)
	return err
}

func (r *Repo) List() ([]*User, error) {
	rows, err := r.db.Query(
		`SELECT id, username, role, created_at FROM users ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*User
	for rows.Next() {
		var u User
		var createdAt string
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &createdAt); err != nil {
			return nil, err
		}
		u.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		list = append(list, &u)
	}
	return list, rows.Err()
}

func (r *Repo) UpdateRole(id int64, role string) error {
	_, err := r.db.Exec(`UPDATE users SET role = ? WHERE id = ?`, role, id)
	return err
}

func (r *Repo) UpdatePassword(id int64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, string(hash), id)
	return err
}

func (r *Repo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func (r *Repo) CountAdmins() (int, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE role = ?`, RoleAdmin).Scan(&n)
	return n, err
}

func (r *Repo) Exists() (bool, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n)
	return n > 0, err
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password)) == nil
}

func (u *User) HasRole(roles ...string) bool {
	for _, r := range roles {
		if u.Role == r {
			return true
		}
	}
	return false
}

// RequestLogEntry — запись лога запроса (игроки/группы/отслеживание)
type RequestLogEntry struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"` // RFC3339 or datetime for JSON
}

func (r *Repo) LogRequest(userID int64, method, path string) error {
	_, err := r.db.Exec(
		`INSERT INTO request_logs (user_id, method, path) VALUES (?, ?, ?)`,
		userID, method, path,
	)
	return err
}

func (r *Repo) GetRequestLogs(userID int64, limit int) ([]RequestLogEntry, error) {
	if limit <= 0 {
		limit = 500
	}
	if limit > 5000 {
		limit = 5000
	}
	rows, err := r.db.Query(
		`SELECT id, user_id, method, path, created_at FROM request_logs WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []RequestLogEntry
	for rows.Next() {
		var e RequestLogEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Method, &e.Path, &e.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}
