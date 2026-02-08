import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import {
  fetchAuthSettings,
  updateAuthSettings,
  checkAuth,
  wipeDatabase,
  fetchAdminUsers,
  adminCreateUser,
  adminUpdateUser,
  adminDeleteUser,
  fetchUserRequestLogs,
  type AuthSettingsBody,
  type AdminUser,
  type RequestLogEntry,
} from '../api/client'

const ROLES = [
  { value: 'admin', label: 'admin' },
  { value: 'editor', label: 'editor' },
  { value: 'viewer', label: 'viewer' },
]
import './SettingsPage.css'

export function SettingsPage() {
  const [cdnAuth, setCdnAuth] = useState('')
  const [cfClearance, setCfClearance] = useState('')
  const [session, setSession] = useState('')
  const [userInfo, setUserInfo] = useState('')
  const [acsrf, setAcsrf] = useState('')
  const [configured, setConfigured] = useState(false)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [checking, setChecking] = useState(false)
  const [checkResult, setCheckResult] = useState<'ok' | 'fail' | null>(null)
  const [wiping, setWiping] = useState(false)
  const [wipeMessage, setWipeMessage] = useState<string | null>(null)

  const [users, setUsers] = useState<AdminUser[]>([])
  const [usersLoading, setUsersLoading] = useState(false)
  const [newUsername, setNewUsername] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [newRole, setNewRole] = useState('viewer')
  const [addingUser, setAddingUser] = useState(false)
  const [editUser, setEditUser] = useState<AdminUser | null>(null)
  const [editRole, setEditRole] = useState('')
  const [editPassword, setEditPassword] = useState('')
  const [savingUser, setSavingUser] = useState(false)
  const [deletingId, setDeletingId] = useState<number | null>(null)
  const [userError, setUserError] = useState<string | null>(null)

  const [logsUserId, setLogsUserId] = useState<number | null>(null)
  const [logs, setLogs] = useState<RequestLogEntry[]>([])
  const [logsLoading, setLogsLoading] = useState(false)

  useEffect(() => {
    fetchAuthSettings()
      .then((r) => setConfigured(Boolean(r?.configured)))
      .catch(() => setConfigured(false))
      .finally(() => setLoading(false))
  }, [])

  const loadUsers = () => {
    setUsersLoading(true)
    fetchAdminUsers()
      .then((r) => setUsers(r?.users ?? []))
      .catch(() => setUsers([]))
      .finally(() => setUsersLoading(false))
  }
  useEffect(() => {
    loadUsers()
  }, [])

  useEffect(() => {
    if (logsUserId == null) {
      setLogs([])
      return
    }
    setLogsLoading(true)
    fetchUserRequestLogs(logsUserId, 500)
      .then((r) => setLogs(r?.logs ?? []))
      .catch(() => setLogs([]))
      .finally(() => setLogsLoading(false))
  }, [logsUserId])

  const handleAddUser = () => {
    if (!newUsername.trim() || !newPassword.trim()) return
    setAddingUser(true)
    setUserError(null)
    adminCreateUser(newUsername.trim(), newPassword, newRole)
      .then((r) => {
        setUsers(r?.users ?? [])
        setNewUsername('')
        setNewPassword('')
        setNewRole('viewer')
      })
      .catch((e) => setUserError(e.message))
      .finally(() => setAddingUser(false))
  }

  const handleSaveEditUser = () => {
    if (!editUser) return
    const payload: { role?: string; password?: string } = {}
    if (editRole !== editUser.role) payload.role = editRole
    if (editPassword.trim()) payload.password = editPassword.trim()
    if (Object.keys(payload).length === 0) {
      setEditUser(null)
      return
    }
    setSavingUser(true)
    setUserError(null)
    adminUpdateUser(editUser.id, payload)
      .then((r) => {
        setUsers(r?.users ?? [])
        setEditUser(null)
        setEditRole('')
        setEditPassword('')
      })
      .catch((e) => setUserError(e.message))
      .finally(() => setSavingUser(false))
  }

  const handleDeleteUser = (u: AdminUser) => {
    if (!window.confirm(`Удалить пользователя ${u.username}?`)) return
    setDeletingId(u.id)
    setUserError(null)
    adminDeleteUser(u.id)
      .then(loadUsers)
      .catch((e) => setUserError(e.message))
      .finally(() => setDeletingId(null))
  }

  const openEdit = (u: AdminUser) => {
    setEditUser(u)
    setEditRole(u.role)
    setEditPassword('')
  }

  const handleSave = () => {
    if (!cdnAuth.trim()) {
      setError('cdn-auth обязателен')
      return
    }
    setSaving(true)
    setError(null)
    setMessage(null)
    const body: AuthSettingsBody = {
      cdn_auth: cdnAuth.trim(),
      cf_clearance: cfClearance.trim() || undefined,
      session: session.trim() || undefined,
      user_info: userInfo.trim() || undefined,
      acsrf: acsrf.trim() || undefined,
    }
    updateAuthSettings(body)
      .then(() => {
        setMessage('Авторизация обновлена')
        setConfigured(true)
      })
      .catch((e) => setError(e.message))
      .finally(() => setSaving(false))
  }

  const handleWipeDb = () => {
    if (!window.confirm('Очистить всю базу? Будут удалены все игроки, группы, история онлайна и отслеживание. Учётные записи (логин) сохранятся. Отменить действие нельзя.')) return
    setWiping(true)
    setWipeMessage(null)
    setError(null)
    wipeDatabase()
      .then((r) => setWipeMessage(r?.message ?? 'База очищена'))
      .catch((e) => setError(e.message))
      .finally(() => setWiping(false))
  }

  const handleCheck = () => {
    setChecking(true)
    setCheckResult(null)
    setError(null)
    checkAuth()
      .then((r) => {
        setCheckResult(r.ok ? 'ok' : 'fail')
        if (!r.ok && r.error) setError(r.error)
      })
      .catch((e) => {
        setCheckResult('fail')
        setError(e.message)
      })
      .finally(() => setChecking(false))
  }

  if (loading) {
    return (
      <div className="settings-page">
        <div className="loading">Загрузка...</div>
      </div>
    )
  }

  return (
    <div className="settings-page">
      <Link to="/" className="back-link">← Назад</Link>

      <section className="settings-section">
        <h2>CFtools авторизация</h2>
        <p className="settings-hint">
          Если авторизация слетела — зайди на <a href="https://auth.cftools.cloud" target="_blank" rel="noreferrer">auth.cftools.cloud</a>,
          открой DevTools → Application → Cookies и скопируй значения.
        </p>
        {configured && (
          <p className="settings-status configured">Авторизация настроена</p>
        )}
        <button
          type="button"
          className="check-auth-btn"
          onClick={handleCheck}
          disabled={checking || !configured}
          title="Проверить, работают ли cookies (реальный запрос к CF API)"
        >
          {checking ? 'Проверка...' : 'Проверить авторизацию'}
        </button>
        {checkResult === 'ok' && <p className="settings-status success">✓ Авторизация работает</p>}
        {checkResult === 'fail' && (
          <p className="settings-status fail">
            ✗ Не работает — обнови cookies в DevTools и сохрани заново
          </p>
        )}

        <div className="auth-form">
          <label>
            <span>cdn-auth (обязательно)</span>
            <textarea
              rows={2}
              placeholder="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
              value={cdnAuth}
              onChange={(e) => setCdnAuth(e.target.value)}
            />
          </label>
          <label>
            <span>cf_clearance</span>
            <input
              type="text"
              placeholder="LhcQYxxMAg9HVdqefqhEBJFm7VKPtlA8Nwe1cKQ_gTk-..."
              value={cfClearance}
              onChange={(e) => setCfClearance(e.target.value)}
            />
          </label>
          <label>
            <span>session</span>
            <input
              type="text"
              placeholder=".eJxtzD0KAjEQQOG7pHaX..."
              value={session}
              onChange={(e) => setSession(e.target.value)}
            />
          </label>
          <label>
            <span>user_info (JSON)</span>
            <textarea
              rows={3}
              placeholder='{"profile": {"display_name": "..."}, "user": {"cftools_id": "..."}}'
              value={userInfo}
              onChange={(e) => setUserInfo(e.target.value)}
            />
          </label>
          <label>
            <span>acsrf</span>
            <input
              type="text"
              placeholder="XdsU0Gw4A2LFiwECHefU8Kwor51gffqgaoTbfkvg1Bo="
              value={acsrf}
              onChange={(e) => setAcsrf(e.target.value)}
            />
          </label>

          {error && <div className="error-msg">{error}</div>}
          {message && <div className="success-msg">{message}</div>}

          <button
            type="button"
            className="save-btn"
            onClick={handleSave}
            disabled={saving || !cdnAuth.trim()}
          >
            {saving ? 'Сохранение...' : 'Сохранить'}
          </button>
        </div>
      </section>

      <section className="settings-section settings-section-db">
        <h2>Управление базой</h2>
        <p className="settings-hint">
          Очистка удалит всех игроков, группы, историю онлайна и список отслеживания. Учётные записи (логин) не затрагиваются.
        </p>
        {wipeMessage && <p className="settings-status success">{wipeMessage}</p>}
        <button
          type="button"
          className="wipe-db-btn"
          onClick={handleWipeDb}
          disabled={wiping}
          title="Удалить все данные из базы"
        >
          {wiping ? 'Очистка...' : 'Очистить всю базу'}
        </button>
      </section>

      <section className="settings-section settings-section-users admin-panel">
        <div className="admin-panel-header">
          <h2>Пользователи</h2>
          <p className="admin-panel-desc">
            Роли: <span className="role-pill role-admin">admin</span> — настройки и доступ ко всему,{' '}
            <span className="role-pill role-editor">editor</span> — игроки и группы,{' '}
            <span className="role-pill role-viewer">viewer</span> — только просмотр.
          </p>
        </div>
        {userError && <div className="error-msg admin-err">{userError}</div>}
        {usersLoading ? (
          <div className="admin-loading">Загрузка списка…</div>
        ) : (
          <>
            <div className="admin-users-card">
              <div className="admin-users-add-block">
                <h3 className="admin-users-add-title">Добавить пользователя</h3>
                <div className="admin-users-add-fields">
                  <label className="admin-users-field">
                    <span className="admin-users-label">Логин</span>
                    <input
                      type="text"
                      placeholder="Имя для входа"
                      value={newUsername}
                      onChange={(e) => setNewUsername(e.target.value)}
                      className="admin-users-input"
                      autoComplete="off"
                    />
                  </label>
                  <label className="admin-users-field">
                    <span className="admin-users-label">Пароль</span>
                    <input
                      type="password"
                      placeholder="••••••••"
                      value={newPassword}
                      onChange={(e) => setNewPassword(e.target.value)}
                      className="admin-users-input"
                      autoComplete="new-password"
                    />
                  </label>
                  <label className="admin-users-field">
                    <span className="admin-users-label">Роль</span>
                    <select value={newRole} onChange={(e) => setNewRole(e.target.value)} className="admin-users-select">
                      {ROLES.map((r) => (
                        <option key={r.value} value={r.value}>{r.label}</option>
                      ))}
                    </select>
                  </label>
                  <button
                    type="button"
                    className="admin-users-btn-add"
                    onClick={handleAddUser}
                    disabled={addingUser || !newUsername.trim() || !newPassword}
                    title="Создать учётную запись"
                  >
                    {addingUser ? (
                      'Добавление…'
                    ) : (
                      <>+ Добавить</>
                    )}
                  </button>
                </div>
              </div>

              <h3 className="admin-users-list-title">Список пользователей</h3>
              {users.length === 0 ? (
                <div className="admin-users-empty">Пока ни одного пользователя. Добавьте первого выше.</div>
              ) : (
                <div className="admin-users-table-wrap">
                  <table className="admin-users-table">
                    <thead>
                      <tr>
                        <th>Логин</th>
                        <th>Роль</th>
                        <th>Создан</th>
                        <th className="th-actions">Действия</th>
                      </tr>
                    </thead>
                    <tbody>
                      {users.map((u) => (
                        <tr key={u.id}>
                          <td className="td-username">{u.username}</td>
                          <td>
                            <span className={`role-pill role-${u.role}`}>{u.role}</span>
                          </td>
                          <td className="admin-users-created">
                            {u.created_at ? new Date(u.created_at).toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' }) : '—'}
                          </td>
                          <td className="td-actions">
                            <div className="admin-users-actions">
                              <button
                                type="button"
                                className="admin-users-btn-edit"
                                onClick={() => openEdit(u)}
                                title="Изменить роль или пароль"
                              >
                                Изменить
                              </button>
                              <button
                                type="button"
                                className="admin-users-btn-delete"
                                onClick={() => handleDeleteUser(u)}
                                disabled={deletingId === u.id}
                                title="Удалить пользователя"
                              >
                                {deletingId === u.id ? '…' : 'Удалить'}
                              </button>
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>

            {editUser && (
              <div
                className="admin-users-edit-modal"
                onClick={() => { setEditUser(null); setEditRole(''); setEditPassword(''); }}
                role="presentation"
              >
                <div className="admin-users-edit-inner" onClick={(e) => e.stopPropagation()} role="dialog" aria-modal="true" aria-labelledby="edit-user-title">
                  <div className="admin-users-edit-head">
                    <h4 id="edit-user-title">Редактировать: {editUser.username}</h4>
                    <button
                      type="button"
                      className="admin-users-edit-close"
                      onClick={() => { setEditUser(null); setEditRole(''); setEditPassword(''); }}
                      title="Закрыть"
                      aria-label="Закрыть"
                    >
                      ×
                    </button>
                  </div>
                  <label className="admin-users-field">
                    <span className="admin-users-label">Роль</span>
                    <select value={editRole} onChange={(e) => setEditRole(e.target.value)} className="admin-users-select">
                      {ROLES.map((r) => (
                        <option key={r.value} value={r.value}>{r.label}</option>
                      ))}
                    </select>
                  </label>
                  <label className="admin-users-field">
                    <span className="admin-users-label">Новый пароль (не заполняйте, чтобы оставить текущий)</span>
                    <input
                      type="password"
                      placeholder="••••••••"
                      value={editPassword}
                      onChange={(e) => setEditPassword(e.target.value)}
                      className="admin-users-input"
                      autoComplete="new-password"
                    />
                  </label>
                  <div className="admin-users-edit-actions">
                    <button type="button" className="admin-users-btn-save" onClick={handleSaveEditUser} disabled={savingUser}>
                      {savingUser ? 'Сохранение...' : 'Сохранить'}
                    </button>
                    <button type="button" className="admin-users-btn-cancel" onClick={() => { setEditUser(null); setEditRole(''); setEditPassword(''); }}>
                      Отмена
                    </button>
                  </div>
                </div>
              </div>
            )}
          </>
        )}
      </section>

      <section className="settings-section settings-section-logs">
        <h2>Логи запросов</h2>
        <p className="settings-hint">
          Запросы к API игроков, групп и отслеживания по каждому пользователю.
        </p>
        <label className="admin-logs-select-label">
          Пользователь:
          <select
            value={logsUserId ?? ''}
            onChange={(e) => setLogsUserId(e.target.value ? Number(e.target.value) : null)}
            className="admin-logs-select"
          >
            <option value="">— выберите —</option>
            {users.map((u) => (
              <option key={u.id} value={u.id}>{u.username} ({u.role})</option>
            ))}
          </select>
        </label>
        {logsLoading ? (
          <p className="loading-small">Загрузка логов...</p>
        ) : logs.length > 0 ? (
          <div className="admin-logs-table-wrap">
            <table className="admin-logs-table">
              <thead>
                <tr>
                  <th>Время</th>
                  <th>Метод</th>
                  <th>Путь</th>
                </tr>
              </thead>
              <tbody>
                {logs.map((log) => (
                  <tr key={log.id}>
                    <td className="admin-logs-time">{log.created_at ? new Date(log.created_at).toLocaleString('ru-RU') : '—'}</td>
                    <td>{log.method}</td>
                    <td className="admin-logs-path" title={log.path}>{log.path}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : logsUserId != null ? (
          <p className="settings-hint">Нет записей за выбранный период.</p>
        ) : null}
      </section>
    </div>
  )
}
