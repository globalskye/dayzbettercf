import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  fetchPlayer,
  syncPlayer,
  fetchGroups,
  addToGroup,
  fetchTracked,
  addTracked,
  removeTracked,
  fetchPlayerHistory,
  formatPlaytime,
  formatDate,
  type Player,
  type Group,
  type HistoryRecord,
} from '../api/client'
import { PlayerHistorySection } from '../components/PlayerHistorySection'
import './PlayerPage.css'

export function PlayerPage() {
  const { id } = useParams<{ id: string }>()
  const [player, setPlayer] = useState<Player | null>(null)
  const [groups, setGroups] = useState<Group[]>([])
  const [loading, setLoading] = useState(true)
  const [syncing, setSyncing] = useState(false)
  const [addGroupId, setAddGroupId] = useState<number | null>(null)
  const [addGroupAlias, setAddGroupAlias] = useState('')
  const [isTracked, setIsTracked] = useState(false)
  const [trackedCount, setTrackedCount] = useState(0)
  const [history, setHistory] = useState<HistoryRecord[]>([])
  const [loadingHistory, setLoadingHistory] = useState(false)
  const [historyError, setHistoryError] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  const load = () => {
    if (!id) return
    setLoading(true)
    setError(null)
    syncPlayer(id, false)
      .then((p) => {
        setPlayer(p)
        setLoading(false)
      })
      .catch(() => {
        fetchPlayer(id)
          .then((p) => {
            setPlayer(p)
            setLoading(false)
          })
          .catch((e) => {
            setError(e.message)
            setLoading(false)
          })
      })
  }

  useEffect(() => {
    load()
  }, [id])

  useEffect(() => {
    fetchGroups()
      .then((r) => setGroups(Array.isArray(r?.groups) ? r.groups : []))
      .catch(() => setGroups([]))
  }, [])

  useEffect(() => {
    if (!id) return
    fetchTracked()
      .then((r) => {
        const list = Array.isArray(r?.players) ? r.players : []
        setTrackedCount(list.length)
        setIsTracked(list.some((p) => p.cftools_id === id))
      })
      .catch(() => {})
  }, [id])

  // История — только после загрузки игрока (уже в базе после sync/fetch)
  useEffect(() => {
    if (!id || !player) return
    setLoadingHistory(true)
    setHistoryError(null)
    fetchPlayerHistory(id, 200)
      .then((r) => setHistory(Array.isArray(r?.history) ? r.history : []))
      .catch((e) => {
        setHistory([])
        setHistoryError(e?.message?.includes('not found') ? 'Игрок не в базе' : 'Не удалось загрузить историю')
      })
      .finally(() => setLoadingHistory(false))
  }, [id, player])

  const handleSync = () => {
    if (!id) return
    setSyncing(true)
    setError(null)
    syncPlayer(id)
      .then(setPlayer)
      .catch((e) => setError(e.message))
      .finally(() => setSyncing(false))
  }

  const handleAddToGroup = () => {
    if (!id || !addGroupId) return
    setError(null)
    addToGroup(addGroupId, id, addGroupAlias.trim() || undefined)
      .then(() => {
        setAddGroupId(null)
        setAddGroupAlias('')
      })
      .catch((e) => setError(e.message))
  }

  const handleToggleTracked = () => {
    if (!id) return
    setError(null)
    if (isTracked) {
      removeTracked(id).then(() => setIsTracked(false)).catch((e) => setError(e.message))
    } else {
      addTracked(id)
        .then(() => setIsTracked(true))
        .catch((e) => setError(e.message))
    }
  }

  if (loading) return <div className="player-page loading">Загрузка...</div>
  if (error || !player)
    return (
      <div className="player-page">
        <div className="error-msg">{error || 'Игрок не найден'}</div>
        <Link to="/" className="back-link">← Назад</Link>
      </div>
    )

  return (
    <div className="player-page">
      <Link to="/" className="back-link">← Назад к списку</Link>

      <header className="player-header">
        <div className="player-header-main">
          {player.avatar ? (
            <img src={player.avatar} alt="" className="player-avatar-lg" />
          ) : (
            <div className="player-avatar-placeholder-lg" />
          )}
          <div className="player-header-info">
            <h1>{player.display_name}</h1>
            <div className="player-header-badges">
              {player.online && <span className="badge online">online</span>}
              {player.is_bot && <span className="badge">bot</span>}
              {player.bans_count > 0 && (
                <span className="badge bans">{player.bans_count} бан(ов)</span>
              )}
            </div>
            <p className="player-id">CFtools ID: {player.cftools_id}</p>
            <div className="player-actions">
              <button
                className="sync-btn"
                onClick={handleSync}
                disabled={syncing}
              >
                {syncing ? 'Обновление...' : 'Обновить из CFtools'}
              </button>
              {trackedCount < 10 && (
                <button
                  type="button"
                  className={`tracked-btn ${isTracked ? 'tracked' : ''}`}
                  onClick={handleToggleTracked}
                  title={isTracked ? 'Убрать из отслеживания' : 'Добавить в отслеживание (обновление каждые 5 мин)'}
                >
                  {isTracked ? 'В отслеживании' : 'В отслеживание'}
                </button>
              )}
              {Array.isArray(groups) && groups.length > 0 && (
                <span className="add-to-group">
                  <select
                    value={addGroupId ?? ''}
                    onChange={(e) => {
                      setAddGroupId(e.target.value ? Number(e.target.value) : null)
                      setAddGroupAlias('')
                    }}
                  >
                    <option value="">Добавить в группу...</option>
                    {(groups ?? []).map((g) => (
                      <option key={g.id} value={g.id}>{g.name}</option>
                    ))}
                  </select>
                  {addGroupId && (
                    <>
                      <input
                        type="text"
                        className="add-alias-input"
                        placeholder="Алиас (как запомнил)"
                        value={addGroupAlias}
                        onChange={(e) => setAddGroupAlias(e.target.value)}
                      />
                      <button type="button" className="add-btn" onClick={handleAddToGroup}>
                        Добавить
                      </button>
                    </>
                  )}
                </span>
              )}
            </div>
          </div>
        </div>
      </header>

      {error && <div className="error-msg">{error}</div>}

      <div className="player-grid-info">
        {player.steam64 && (
          <section className="info-card info-card-steam">
            <h3>Steam</h3>
            <div className="steam-block">
              {player.steam_avatar && (
                <img src={player.steam_avatar} alt="" className="steam-avatar" />
              )}
              <div>
                <a href={`https://steamcommunity.com/profiles/${player.steam64}`} target="_blank" rel="noreferrer" className="steam-link">
                  {player.steam_persona || player.steam64}
                </a>
                <p className="steam-id">Steam64: {player.steam64}</p>
                <p className="steam-bans">
                  VAC: {player.steam_vac_bans ?? 0} · EAC/Game: {player.steam_game_bans ?? 0}
                </p>
              </div>
            </div>
            <a href={`https://app.cftools.cloud/profile/${player.cftools_id}`} target="_blank" rel="noreferrer" className="cftools-link">
              Открыть в CFtools →
            </a>
          </section>
        )}

        <section className="info-card">
          <h3>Статистика</h3>
          <dl>
            <dt>Время в игре</dt>
            <dd>{formatPlaytime(player.playtime_sec)}</dd>
            <dt>Сессий</dt>
            <dd>{player.sessions_count}</dd>
            <dt>Связанных аккаунтов</dt>
            <dd>{player.linked_accounts_count}</dd>
            <dt>Банов</dt>
            <dd>{player.bans_count}</dd>
          </dl>
        </section>

        <PlayerHistorySection
          history={history}
          loading={loadingHistory}
          error={historyError}
        />

        <section className="info-card">
          <h3>Даты</h3>
          <dl>
            <dt>Последняя активность</dt>
            <dd>{formatDate(player.last_activity_at)}</dd>
            <dt>Последнее обновление</dt>
            <dd>{formatDate(player.last_seen_at)}</dd>
            <dt>Добавлен в базу</dt>
            <dd>{formatDate(player.created_at)}</dd>
          </dl>
        </section>

        {player.nicknames && player.nicknames.length > 0 && (
          <section className="info-card">
            <h3>Ники / алиасы</h3>
            <div className="tag-list">
              {player.nicknames.map((n) => (
                <span key={n} className="tag">{n}</span>
              ))}
            </div>
          </section>
        )}

        {(() => {
          const links = player.linked_accounts?.length
            ? player.linked_accounts
            : (player.linked_cftools_ids ?? []).map((cftools_id) => ({ cftools_id, confirmed: false, trusted: false }))
          if (links.length === 0) return null
          return (
            <section className="info-card">
              <h3>Связанные аккаунты ({links.length})</h3>
              <p className="link-list-hint">
                <strong>confirmed</strong> — CFTools подтвердил связь (например с одного устройства); <strong>trusted</strong> — отмечен доверенным. Клик — переход в профиль.
              </p>
              <div className="link-list">
                {links.slice(0, 20).map((link) => (
                  <Link key={link.cftools_id} to={`/players/${link.cftools_id}`} className="link-item">
                    <span className="link-item-id">{link.cftools_id}</span>
                    {(link.confirmed || link.trusted) && (
                      <span className="link-item-badges">
                        {link.confirmed && <span className="badge badge-confirmed">confirmed</span>}
                        {link.trusted && <span className="badge badge-trusted">trusted</span>}
                      </span>
                    )}
                  </Link>
                ))}
              </div>
              {links.length > 20 && <span className="more">+{links.length - 20} ещё</span>}
            </section>
          )
        })()}

        {player.raw_bans && (() => {
          try {
            const data = JSON.parse(player.raw_bans)
            const bans = data.bans as Array<{ reason?: string; created_at?: string }>
            const banlists = (data.banlists || {}) as Record<string, { identifier?: string }>
            if (bans && bans.length > 0) {
              return (
                <section className="info-card info-card-bans">
                  <h3>CFtools баны</h3>
                  <ul className="bans-list">
                    {bans.map((b: { reason?: string; created_at?: string; banlist_id?: string }, i: number) => (
                      <li key={i}>
                        <span className="ban-reason">{b.reason || '—'}</span>
                        <span className="ban-date">{formatDate(b.created_at)}</span>
                        {b.banlist_id && banlists[b.banlist_id] && (
                          <span className="ban-list">{banlists[b.banlist_id].identifier}</span>
                        )}
                      </li>
                    ))}
                  </ul>
                </section>
              )
            }
          } catch {
            /* ignore */
          }
          return null
        })()}

        {player.raw_battleye && (() => {
          try {
            const data = JSON.parse(player.raw_battleye)
            const records = (data.records || []) as Array<{ id?: string; date?: string }>
            if (records && records.length > 0) {
              return (
                <section className="info-card info-card-battleye">
                  <h3>BattlEye</h3>
                  <ul className="bans-list">
                    {records.map((r: { id?: string; date?: string }, i: number) => (
                      <li key={i}>
                        <span className="ban-reason">{r.id || '—'}</span>
                        <span className="ban-date">{formatDate(r.date)}</span>
                      </li>
                    ))}
                  </ul>
                </section>
              )
            }
          } catch {
            /* ignore */
          }
          return null
        })()}
      </div>
    </div>
  )
}
