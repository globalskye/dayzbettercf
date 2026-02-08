import { useEffect, useState } from 'react'
import { usePlayerModal } from '../context/PlayerModalContext'
import { PlayerActionsModal } from '../components/PlayerActionsModal'
import {
  fetchPlayers,
  searchPlayersLocal,
  fetchGroups,
  fetchTracked,
  formatPlaytime,
  formatDate,
  type Player,
  type Group,
} from '../api/client'
import { LoadingSpinner, LoadingOverlay } from '../components/LoadingSpinner'
import { HistoryPopover } from '../components/HistoryPopover'
import './BasePage.css'
import './HomePage.css'

export function BasePage() {
  const [players, setPlayers] = useState<Player[]>([])
  const [total, setTotal] = useState(0)
  const [localQuery, setLocalQuery] = useState('')
  const [loading, setLoading] = useState(true)
  const [popupPlayer, setPopupPlayer] = useState<Player | null>(null)
  const openPlayerModal = usePlayerModal()
  const [groups, setGroups] = useState<Group[]>([])
  const [trackedIds, setTrackedIds] = useState<Set<string>>(new Set())
  const [error, setError] = useState<string | null>(null)
  const [onlyOnline, setOnlyOnline] = useState(false)
  const [onlyBanned, setOnlyBanned] = useState(false)
  const [sort, setSort] = useState<'online' | 'playtime' | 'bans' | 'updated'>('online')
  const [viewMode, setViewMode] = useState<'table' | 'cards'>('table')

  const loadPlayers = () => {
    setLoading(true)
    setError(null)
    fetchPlayers({ online: onlyOnline, banned: onlyBanned, sort, limit: 200 })
      .then((r) => {
        const list = Array.isArray(r?.players) ? r.players : []
        setPlayers(list)
        setTotal(r?.total ?? list.length)
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    if (localQuery.trim()) {
      setLoading(true)
      setError(null)
      searchPlayersLocal(localQuery.trim(), { online: onlyOnline, banned: onlyBanned, sort, limit: 5000 })
        .then((r) => {
          const list = Array.isArray(r?.players) ? r.players : []
          setPlayers(list)
          setTotal(r?.count ?? list.length)
        })
        .catch((e) => setError(e.message))
        .finally(() => setLoading(false))
    } else {
      loadPlayers()
    }
  }, [onlyOnline, onlyBanned, sort])

  useEffect(() => {
    fetchGroups().then((r) => setGroups(Array.isArray(r?.groups) ? r.groups : [])).catch(() => {})
    fetchTracked()
      .then((r) => setTrackedIds(new Set((r?.players ?? []).map((p) => p.cftools_id))))
      .catch(() => {})
  }, [])

  const handleLocalSearch = () => {
    if (!localQuery.trim()) {
      loadPlayers()
      return
    }
    setLoading(true)
    setError(null)
    searchPlayersLocal(localQuery.trim(), { online: onlyOnline, banned: onlyBanned, sort, limit: 5000 })
      .then((r) => {
        const list = Array.isArray(r?.players) ? r.players : []
        setPlayers(list)
        setTotal(r?.count ?? list.length)
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  const handlePopupClose = () => setPopupPlayer(null)

  return (
    <div className="base-page home">
      <section className="search-section">
        <h2 className="base-title">База игроков</h2>
        <p className="search-hint">
          Поиск среди сохранённых в базе. Актуальные данные — на главной (поиск по CF).
        </p>
        <div className="search-row">
          <label className="search-label">Поиск в базе</label>
          <input
            type="text"
            placeholder="Ник или часть..."
            value={localQuery}
            onChange={(e) => setLocalQuery(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleLocalSearch()}
          />
          <button onClick={handleLocalSearch} disabled={loading}>
            {loading ? (
              <span className="btn-loading-wrap">
                <LoadingSpinner variant="dots" size="small" />
                <span>Поиск...</span>
              </span>
            ) : (
              'Искать'
            )}
          </button>
        </div>
      </section>

      {popupPlayer && (
        <PlayerActionsModal
          player={popupPlayer}
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
            В базе {total > 0 && <span className="count">({total})</span>}
            {(players ?? []).filter((p) => p?.online).length > 0 && (
              <span className="online-count"> · онлайн: {(players ?? []).filter((p) => p?.online).length}</span>
            )}
          </h2>
          {!loading && (players ?? []).length > 0 && (
            <div className="players-toolbar">
              <label className="checkbox-label">
                <input type="checkbox" checked={onlyOnline} onChange={(e) => setOnlyOnline(e.target.checked)} />
                Только онлайн
              </label>
              <label className="checkbox-label">
                <input type="checkbox" checked={onlyBanned} onChange={(e) => setOnlyBanned(e.target.checked)} />
                С банами
              </label>
              <select value={sort} onChange={(e) => setSort(e.target.value as typeof sort)}>
                <option value="online">Сортировка: онлайн</option>
                <option value="playtime">Сортировка: часы</option>
                <option value="bans">Сортировка: баны</option>
                <option value="updated">Сортировка: обновлено</option>
              </select>
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
          <LoadingOverlay text="Загрузка базы..." />
        ) : (players ?? []).length === 0 ? (
          <div className="empty">
            В базе никого нет. На главной найдите игроков по CF — они автоматически добавятся сюда.
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
                {(players ?? []).map((p) => (
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
                    <td className="server-cell">
                      <HistoryPopover cftoolsId={p.cftools_id} displayName={p.display_name}>
                        <span className="server-cell-text" title={p.last_server_identifier || ''}>{p.last_server_identifier || '—'}</span>
                      </HistoryPopover>
                    </td>
                    <td>{formatPlaytime(p.playtime_sec ?? 0)}</td>
                    <td>{p.bans_count ?? 0}</td>
                    <td className="date-cell">{formatDate(p.updated_at)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="player-grid">
            {(players ?? []).map((p) => (
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
                    {p.last_server_identifier && <span className="badge server">{p.last_server_identifier}</span>}
                    {(p.bans_count ?? 0) > 0 && (
                      <span className="badge bans">{p.bans_count} бан(ов)</span>
                    )}
                  </div>
                </div>
                <div className="player-card-stats">
                  {p.last_server_identifier && (
                    <>
                      <span className="server-badge" title="Сервер">{p.last_server_identifier}</span>
                      <span>•</span>
                    </>
                  )}
                  <span title="Время в игре">{formatPlaytime(p.playtime_sec ?? 0)}</span>
                  <span>•</span>
                  <span>{(p.sessions_count ?? 0)} сессий</span>
                  <span>•</span>
                  <span>{(p.linked_accounts_count ?? 0)} связей</span>
                </div>
                <div className="player-card-date">
                  Обновлено: {formatDate(p?.updated_at)}
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  )
}
