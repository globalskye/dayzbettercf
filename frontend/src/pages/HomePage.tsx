import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { usePlayerModal } from '../context/PlayerModalContext'
import { PlayerActionsModal } from '../components/PlayerActionsModal'
import {
  getProfileStates,
  syncBatch,
  fetchGroups,
  fetchTracked,
  formatPlaytime,
  formatDate,
  type Player,
  type ProfileStateItem,
  type Group,
} from '../api/client'
import { LoadingSpinner } from '../components/LoadingSpinner'
import './HomePage.css'

/** Строка списка: либо только state из CF, либо игрок из базы (после «В базу»). */
type SearchRow = ProfileStateItem & Partial<Pick<Player, 'playtime_sec' | 'bans_count' | 'updated_at' | 'last_server_identifier'>>

export function HomePage() {
  const [searchParams] = useSearchParams()
  const [searchResults, setSearchResults] = useState<SearchRow[]>([])
  const [cfQuery, setCfQuery] = useState('')
  const [loading, setLoading] = useState(false)
  const [syncing, setSyncing] = useState(false)
  const [popupPlayer, setPopupPlayer] = useState<SearchRow | null>(null)
  const openPlayerModal = usePlayerModal()
  const [groups, setGroups] = useState<Group[]>([])
  const [trackedIds, setTrackedIds] = useState<Set<string>>(new Set())
  const [error, setError] = useState<string | null>(null)
  const [lightSync, setLightSync] = useState(true)
  const [viewMode, setViewMode] = useState<'table' | 'cards'>('table')

  const qFromUrl = searchParams.get('q') ?? ''

  useEffect(() => {
    fetchGroups().then((r) => setGroups(Array.isArray(r?.groups) ? r.groups : [])).catch(() => {})
    fetchTracked()
      .then((r) => setTrackedIds(new Set((r?.players ?? []).map((p) => p.cftools_id))))
      .catch(() => {})
  }, [])

  useEffect(() => {
    if (!qFromUrl.trim()) return
    setCfQuery(qFromUrl)
    setLoading(true)
    setError(null)
    getProfileStates(qFromUrl.trim())
      .then((r) => setSearchResults(r?.states ?? []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [qFromUrl])

  const handleCftoolsSearch = () => {
    if (!cfQuery.trim()) return
    setLoading(true)
    setError(null)
    getProfileStates(cfQuery.trim())
      .then((r) => {
        setSearchResults(r?.states ?? [])
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  const handleAddToBase = () => {
    if (searchResults.length === 0) return
    const ids = searchResults.map((r) => r.cftools_id)
    setSyncing(true)
    setError(null)
    syncBatch(ids, lightSync)
      .then((r) => {
        const list = (r?.players ?? []).map((p) => ({
          cftools_id: p.cftools_id,
          display_name: p.display_name,
          avatar: p.avatar,
          online: p.online,
          server_name: p.last_server_identifier,
          last_server_identifier: p.last_server_identifier,
          playtime_sec: p.playtime_sec,
          bans_count: p.bans_count,
          updated_at: p.updated_at,
        }))
        setSearchResults(list)
      })
      .catch((e) => setError(e.message))
      .finally(() => setSyncing(false))
  }

  const handlePopupClose = () => setPopupPlayer(null)

  const total = searchResults.length
  const onlineCount = searchResults.filter((p) => p?.online).length
  const serverLabel = (p: SearchRow) => p.server_name ?? p.last_server_identifier ?? '—'

  return (
    <div className="home">
      <section className="search-section">
        <div className="search-row">
          <label className="search-label">Поиск по CF</label>
          <input
            type="text"
            placeholder="Ник или identifier..."
            value={cfQuery}
            onChange={(e) => setCfQuery(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleCftoolsSearch()}
          />
          <button onClick={handleCftoolsSearch} disabled={loading || !cfQuery.trim()} className={loading ? 'btn-loading' : ''}>
            {loading ? (
              <span className="btn-loading-wrap">
                <LoadingSpinner variant="dots" size="small" />
                <span>Поиск...</span>
              </span>
            ) : (
              'Искать'
            )}
          </button>
          <label className="checkbox-label">
            <input type="checkbox" checked={lightSync} onChange={(e) => setLightSync(e.target.checked)} />
            Лёгкий sync (онлайн). Без галки — с банами.
          </label>
        </div>
        <div className="search-hint">
          Поиск по CFtools (без записи в БД). Результаты — список с онлайн/сервером. Чтобы добавить в базу — нажмите «В базу». Просмотр базы — по кнопке «База» в меню.
        </div>
      </section>

      {popupPlayer && (
        <PlayerActionsModal
          player={popupPlayer as Player}
          groups={groups}
          isTracked={trackedIds.has(popupPlayer.cftools_id)}
          onClose={handlePopupClose}
          onSuccess={() => {
            fetchTracked().then((r) => setTrackedIds(new Set((r?.players ?? []).map((p) => p.cftools_id)))).catch(() => {})
            handlePopupClose()
          }}
        />
      )}

      {error && <div className="error-msg">{error}</div>}

      <section className="players-section">
        <div className="players-header">
          <h2>
            Игроки {total > 0 && <span className="count">({total})</span>}
            {onlineCount > 0 && <span className="online-count"> · онлайн: {onlineCount}</span>}
          </h2>
          {!loading && searchResults.length > 0 && (
            <div className="players-toolbar">
              <button
                type="button"
                className="btn-add-base"
                onClick={handleAddToBase}
                disabled={syncing}
              >
                {syncing ? (
                  <span className="btn-loading-wrap">
                    <LoadingSpinner variant="dots" size="small" />
                    <span>Добавляем…</span>
                  </span>
                ) : (
                  'В базу'
                )}
              </button>
              <div className="view-toggle">
                <button
                  type="button"
                  className={viewMode === 'table' ? 'active' : ''}
                  onClick={() => setViewMode('table')}
                  title="Таблица"
                >
                  ≡
                </button>
                <button
                  type="button"
                  className={viewMode === 'cards' ? 'active' : ''}
                  onClick={() => setViewMode('cards')}
                  title="Карточки"
                >
                  ⊞
                </button>
              </div>
            </div>
          )}
        </div>

        {loading ? (
          <div className="loading-block">
            <LoadingSpinner variant="spinner" size="medium" />
            <p>Запрос к CFtools...</p>
          </div>
        ) : searchResults.length === 0 ? (
          <div className="empty">
            Введите ник или identifier и нажмите «Искать» — появятся данные из CF (онлайн/сервер). Добавить в базу — кнопка «В базу».
          </div>
        ) : viewMode === 'table' ? (
          <div className="players-table-wrap">
            <table className="players-table">
              <thead>
                <tr>
                  <th />
                  <th>Ник</th>
                  <th>Онлайн</th>
                  <th>Сервер</th>
                  <th>Игра</th>
                  <th>Банов</th>
                  <th>Обновлено</th>
                </tr>
              </thead>
              <tbody>
                {searchResults.map((p) => (
                  <tr key={p.cftools_id} className={p.online ? 'row-online' : ''}>
                    <td>
                      {p.avatar ? (
                        <img src={p.avatar} alt="" className="player-avatar-sm" />
                      ) : (
                        <div className="player-avatar-placeholder-sm" />
                      )}
                    </td>
                    <td>
                      <button
                        type="button"
                        className="player-link-btn"
                        onClick={() => openPlayerModal?.(p.cftools_id)}
                      >
                        {p.display_name || p.cftools_id}
                      </button>
                      <button
                        type="button"
                        className="btn-add-popup"
                        onClick={(e) => {
                          e.preventDefault()
                          setPopupPlayer(p)
                        }}
                        title="В группу / отслеживание"
                      >
                        +
                      </button>
                    </td>
                    <td>
                      {p.online ? <span className="badge online">online</span> : <span className="badge offline">offline</span>}
                    </td>
                    <td className="server-cell" title={serverLabel(p)}>{serverLabel(p)}</td>
                    <td>{p.playtime_sec != null ? formatPlaytime(p.playtime_sec) : '—'}</td>
                    <td>{p.bans_count != null ? p.bans_count : '—'}</td>
                    <td className="date-cell">{p.updated_at ? formatDate(p.updated_at) : '—'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="player-grid">
            {searchResults.map((p) => (
              <div
                key={p.cftools_id}
                role="button"
                tabIndex={0}
                className={`player-card ${p.online ? 'card-online' : ''}`}
                onClick={() => openPlayerModal?.(p.cftools_id)}
                onKeyDown={(e) => e.key === 'Enter' && openPlayerModal?.(p.cftools_id)}
              >
                <div className="player-card-header">
                  {p.avatar ? (
                    <img src={p.avatar} alt="" className="player-avatar" />
                  ) : (
                    <div className="player-avatar-placeholder" />
                  )}
                  <div className="player-card-info">
                    <span className="player-name">{p.display_name}</span>
                    <button
                      type="button"
                      className="btn-add-popup"
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                        setPopupPlayer(p)
                      }}
                      title="В группу / отслеживание"
                    >
                      +
                    </button>
                    {p.online && <span className="badge online">online</span>}
                    {(p.server_name ?? p.last_server_identifier) && (
                      <span className="badge server">{p.server_name ?? p.last_server_identifier}</span>
                    )}
                    {(p.bans_count ?? 0) > 0 && (
                      <span className="badge bans">{p.bans_count} бан(ов)</span>
                    )}
                  </div>
                </div>
                <div className="player-card-stats">
                  {(p.server_name ?? p.last_server_identifier) && (
                    <>
                      <span className="server-badge" title="Сервер">{p.server_name ?? p.last_server_identifier}</span>
                      <span>•</span>
                    </>
                  )}
                  <span title="Время в игре">{p.playtime_sec != null ? formatPlaytime(p.playtime_sec) : '—'}</span>
                </div>
                {p.updated_at && (
                  <div className="player-card-date">
                    Обновлено: {formatDate(p.updated_at)}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  )
}
