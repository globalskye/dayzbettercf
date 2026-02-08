import { useEffect, useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { usePlayerModal } from '../context/PlayerModalContext'
import { PlayerActionsModal } from '../components/PlayerActionsModal'
import { ActivityHeatmap } from '../components/ActivityHeatmap'
import {
  fetchTracked,
  fetchGroups,
  removeTracked,
  fetchPlayerHistory,
  computeHistoryStats,
  formatPlaytime,
  formatDate,
  type Player,
  type Group,
  type HistoryRecord,
} from '../api/client'
import './TrackedPage.css'
import './HomePage.css'

const MAX_TRACKED = 10

export function TrackedPage() {
  const [players, setPlayers] = useState<Player[]>([])
  const [groups, setGroups] = useState<Group[]>([])
  const [popupPlayer, setPopupPlayer] = useState<Player | null>(null)
  const openPlayerModal = usePlayerModal()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [sort, setSort] = useState<'online' | 'playtime' | 'bans'>('online')
  const [heatmapPlayer, setHeatmapPlayer] = useState<Player | null>(null)
  const [heatmapHistory, setHeatmapHistory] = useState<HistoryRecord[]>([])
  const [heatmapLoading, setHeatmapLoading] = useState(false)

  useEffect(() => {
    if (!heatmapPlayer) return
    setHeatmapLoading(true)
    fetchPlayerHistory(heatmapPlayer.cftools_id, 500)
      .then((r) => setHeatmapHistory(Array.isArray(r?.history) ? r.history : []))
      .catch(() => setHeatmapHistory([]))
      .finally(() => setHeatmapLoading(false))
  }, [heatmapPlayer?.cftools_id])

  const heatmapByDay = useMemo(
    () => (heatmapHistory.length ? computeHistoryStats(heatmapHistory).byDay : {}),
    [heatmapHistory]
  )

  const load = () => {
    setLoading(true)
    setError(null)
    fetchTracked(sort)
      .then((r) => setPlayers(Array.isArray(r?.players) ? r.players : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    load()
    fetchGroups().then((r) => setGroups(Array.isArray(r?.groups) ? r.groups : [])).catch(() => {})
  }, [sort])

  const handleRemove = (cftoolsId: string) => {
    removeTracked(cftoolsId)
      .then(load)
      .catch((e) => setError(e.message))
  }

  if (loading) return <div className="tracked-page loading">Загрузка...</div>

  return (
    <div className="tracked-page">
      <Link to="/" className="back-link">← Назад</Link>

      <div className="tracked-header">
        <h1>Отслеживаемые игроки</h1>
        <p className="tracked-desc">
          До {MAX_TRACKED} игроков. Та же таблица и инструменты — в группу, профиль по клику.
        </p>
        {!loading && (players ?? []).length > 0 && (
          <div className="players-toolbar">
            <select value={sort} onChange={(e) => setSort(e.target.value as typeof sort)}>
              <option value="online">Сортировка: онлайн</option>
              <option value="playtime">Сортировка: часы</option>
              <option value="bans">Сортировка: баны</option>
            </select>
          </div>
        )}
      </div>

      {error && <div className="error-msg">{error}</div>}

      {(players ?? []).length === 0 ? (
        <div className="empty-hint">
          Добавь игроков с главной — кнопка «+» → «Добавить в отслеживание».
        </div>
      ) : (
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
                <th className="th-activity">Активность</th>
                <th />
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
                  <td className="server-cell" title={p.last_server_identifier || ''}>{p.last_server_identifier || '—'}</td>
                  <td>{formatPlaytime(p.playtime_sec ?? 0)}</td>
                  <td>{p.bans_count ?? 0}</td>
                  <td className="date-cell">{formatDate(p.updated_at)}</td>
                  <td className="td-activity">
                    <button
                      type="button"
                      className="tracked-heatmap-btn"
                      onClick={() => setHeatmapPlayer(heatmapPlayer?.cftools_id === p.cftools_id ? null : p)}
                      title="График активности (GitHub-стиль)"
                    >
                      📊
                    </button>
                  </td>
                  <td>
                    <button
                      type="button"
                      className="tracked-remove-btn"
                      onClick={() => handleRemove(p.cftools_id)}
                      title="Убрать из отслеживания"
                    >
                      ×
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {heatmapPlayer && (
        <div className="tracked-heatmap-dropdown-overlay" onClick={() => setHeatmapPlayer(null)}>
          <div className="tracked-heatmap-dropdown" onClick={(e) => e.stopPropagation()}>
            <div className="tracked-heatmap-dropdown-header">
              <span>Активность — {heatmapPlayer.display_name || heatmapPlayer.cftools_id}</span>
              <button type="button" className="tracked-heatmap-close" onClick={() => setHeatmapPlayer(null)} title="Закрыть">
                ×
              </button>
            </div>
            <div className="tracked-heatmap-dropdown-body">
              {heatmapLoading ? (
                <p className="tracked-heatmap-loading">Загрузка истории...</p>
              ) : (
                <ActivityHeatmap byDay={heatmapByDay} history={heatmapHistory} />
              )}
            </div>
          </div>
        </div>
      )}

      {popupPlayer && (
        <PlayerActionsModal
          player={popupPlayer}
          groups={groups}
          isTracked={true}
          onClose={() => setPopupPlayer(null)}
          onSuccess={() => {
            load()
            setPopupPlayer(null)
          }}
        />
      )}
    </div>
  )
}

