import { useEffect, useState } from 'react'
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
import { PlayerHistorySection } from './PlayerHistorySection'
import './PlayerModal.css'

export function PlayerModal({
  playerId,
  onClose,
  onOpenPlayer,
}: {
  playerId: string
  onClose: () => void
  onOpenPlayer?: (id: string) => void
}) {
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
    if (!playerId) return
    setLoading(true)
    setError(null)
    // Сначала полный запрос к CF (sync), чтобы профиль был актуальным; при ошибке — из базы
    syncPlayer(playerId, false)
      .then((p) => {
        setPlayer(p)
        setLoading(false)
      })
      .catch(() => {
        fetchPlayer(playerId)
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
  }, [playerId])

  useEffect(() => {
    fetchGroups().then((r) => setGroups(Array.isArray(r?.groups) ? r.groups : [])).catch(() => {})
  }, [])

  useEffect(() => {
    if (!playerId) return
    fetchTracked()
      .then((r) => {
        const list = Array.isArray(r?.players) ? r.players : []
        setTrackedCount(list.length)
        setIsTracked(list.some((p) => p.cftools_id === playerId))
      })
      .catch(() => {})
  }, [playerId])

  // История — только после загрузки игрока (он уже в базе после sync/fetch), иначе 404
  useEffect(() => {
    if (!playerId || !player) return
    setLoadingHistory(true)
    setHistoryError(null)
    fetchPlayerHistory(playerId, 200)
      .then((r) => {
        setHistory(Array.isArray(r?.history) ? r.history : [])
      })
      .catch((e) => {
        setHistory([])
        setHistoryError(e?.message?.includes('not found') || e?.message?.includes('404') ? 'Игрок не в базе — история недоступна' : 'Не удалось загрузить историю')
      })
      .finally(() => setLoadingHistory(false))
  }, [playerId, player])

  const handleSync = () => {
    if (!playerId) return
    setSyncing(true)
    setError(null)
    syncPlayer(playerId)
      .then(setPlayer)
      .catch((e) => setError(e.message))
      .finally(() => setSyncing(false))
  }

  const handleAddToGroup = () => {
    if (!playerId || !addGroupId) return
    setError(null)
    addToGroup(addGroupId, playerId, addGroupAlias.trim() || undefined)
      .then(() => {
        setAddGroupId(null)
        setAddGroupAlias('')
      })
      .catch((e) => setError(e.message))
  }

  const handleToggleTracked = () => {
    if (!playerId) return
    setError(null)
    if (isTracked) {
      removeTracked(playerId).then(() => setIsTracked(false)).catch((e) => setError(e.message))
    } else {
      addTracked(playerId).then(() => setIsTracked(true)).catch((e) => setError(e.message))
    }
  }

  return (
    <div className="player-modal-overlay" onClick={onClose}>
      <div className="player-modal" onClick={(e) => e.stopPropagation()}>
        <button type="button" className="player-modal-close" onClick={onClose} title="Закрыть">
          ×
        </button>
        {loading ? (
          <div className="player-modal-loading">Загрузка...</div>
        ) : error || !player ? (
          <div className="player-modal-error">{error || 'Игрок не найден'}</div>
        ) : (
          <>
            <header className="player-modal-header">
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
                    {player.last_server_identifier && (
                      <span className="badge server">{player.last_server_identifier}</span>
                    )}
                    {player.is_bot && <span className="badge">bot</span>}
                    {player.bans_count > 0 && (
                      <span className="badge bans">{player.bans_count} бан(ов)</span>
                    )}
                  </div>
                  <div className="player-actions">
                    <button className="sync-btn" onClick={handleSync} disabled={syncing}>
                      {syncing ? 'Обновление...' : 'Обновить'}
                    </button>
                    {trackedCount < 10 && (
                      <button
                        type="button"
                        className={`tracked-btn ${isTracked ? 'tracked' : ''}`}
                        onClick={handleToggleTracked}
                      >
                        {isTracked ? 'В отслеживании' : 'В отслеживание'}
                      </button>
                    )}
                    {groups.length > 0 && (
                      <span className="add-to-group">
                        <select
                          value={addGroupId ?? ''}
                          onChange={(e) =>
                            setAddGroupId(e.target.value ? Number(e.target.value) : null)
                          }
                        >
                          <option value="">В группу...</option>
                          {groups.map((g) => (
                            <option key={g.id} value={g.id}>
                              {g.name}
                            </option>
                          ))}
                        </select>
                        {addGroupId && (
                          <>
                            <input
                              type="text"
                              className="add-alias-input"
                              placeholder="Алиас"
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

            <div className="player-modal-body">
              <div className="player-grid-info">
                {player.steam64 && (
                  <section className="info-card info-card-steam">
                    <h3>Steam</h3>
                    <div className="steam-block">
                      {player.steam_avatar && (
                        <img src={player.steam_avatar} alt="" className="steam-avatar" />
                      )}
                      <div>
                        <a
                          href={`https://steamcommunity.com/profiles/${player.steam64}`}
                          target="_blank"
                          rel="noreferrer"
                          className="steam-link"
                        >
                          {player.steam_persona || player.steam64}
                        </a>
                        <p className="steam-bans">
                          VAC: {player.steam_vac_bans ?? 0} · EAC: {player.steam_game_bans ?? 0}
                        </p>
                      </div>
                    </div>
                    <a
                      href={`https://app.cftools.cloud/profile/${player.cftools_id}`}
                      target="_blank"
                      rel="noreferrer"
                      className="cftools-link"
                    >
                      CFtools →
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
                    <dt>Банов</dt>
                    <dd>{player.bans_count}</dd>
                  </dl>
                </section>

                <PlayerHistorySection
                  history={history}
                  loading={loadingHistory}
                  error={historyError}
                  compact
                  showHeatmap={false}
                />

                <section className="info-card">
                  <h3>Даты</h3>
                  <dl>
                    <dt>Последняя активность</dt>
                    <dd>{formatDate(player.last_activity_at)}</dd>
                    <dt>Обновлено</dt>
                    <dd>{formatDate(player.updated_at)}</dd>
                  </dl>
                </section>

                {player.nicknames && player.nicknames.length > 0 && (
                  <section className="info-card">
                    <h3>Ники</h3>
                    <div className="tag-list">
                      {player.nicknames.map((n) => (
                        <span key={n} className="tag">
                          {n}
                        </span>
                      ))}
                    </div>
                  </section>
                )}

                {player.linked_cftools_ids && player.linked_cftools_ids.length > 0 && (
                  <section className="info-card">
                    <h3>Связанные ({player.linked_cftools_ids.length})</h3>
                    <div className="link-list">
                      {player.linked_cftools_ids.slice(0, 10).map((lid) => (
                        <button
                          key={lid}
                          type="button"
                          className="link-item"
                          onClick={() => onOpenPlayer?.(lid) ?? onClose()}
                        >
                          {lid}
                        </button>
                      ))}
                    </div>
                  </section>
                )}

                {player.raw_bans &&
                  (() => {
                    try {
                      const data = JSON.parse(player.raw_bans)
                      const bans = data.bans as Array<{ reason?: string; created_at?: string }>
                      if (bans && bans.length > 0) {
                        return (
                          <section className="info-card info-card-bans">
                            <h3>CFtools баны</h3>
                            <ul className="bans-list">
                              {bans.slice(0, 5).map((b: { reason?: string; created_at?: string }, i: number) => (
                                <li key={i}>
                                  <span className="ban-reason">{b.reason || '—'}</span>
                                  <span className="ban-date">{formatDate(b.created_at)}</span>
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
          </>
        )}
      </div>
    </div>
  )
}
