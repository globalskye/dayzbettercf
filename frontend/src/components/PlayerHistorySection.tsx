import { useState, useMemo } from 'react'
import {
  formatHistoryLine,
  formatDuration,
  computeHistoryStats,
  getDateKey,
  getWeekKey,
  type HistoryRecord,
} from '../api/client'
import { ActivityHeatmap } from './ActivityHeatmap'
import './PlayerHistorySection.css'

type PeriodFilter = 'all' | 'day' | 'week'
type SortOrder = 'newest' | 'oldest'

export function PlayerHistorySection({
  history,
  loading,
  error,
  compact = false,
  showHeatmap = true,
}: {
  history: HistoryRecord[]
  loading: boolean
  error?: string | null
  compact?: boolean
  /** В модалках не показывать хитмап — только на странице отслеживания в выпадающем окне */
  showHeatmap?: boolean
}) {
  const [period, setPeriod] = useState<PeriodFilter>('all')
  const [dayKey, setDayKey] = useState<string>(() => {
    const d = new Date()
    return d.toISOString().slice(0, 10)
  })
  const [weekKey, setWeekKey] = useState<string>(() => {
    const d = new Date()
    const day = d.getUTCDay()
    const diff = day === 0 ? -6 : 1 - day
    d.setUTCDate(d.getUTCDate() + diff)
    return d.toISOString().slice(0, 10)
  })
  const [sortOrder, setSortOrder] = useState<SortOrder>('newest')

  const stats = useMemo(() => computeHistoryStats(history), [history])

  const dayKeys = useMemo(() => {
    const set = new Set<string>()
    history.forEach((h) => {
      const dk = getDateKey(h.ts)
      if (dk) set.add(dk)
    })
    return Array.from(set).sort().reverse().slice(0, 31)
  }, [history])

  const weekKeys = useMemo(() => {
    const set = new Set<string>()
    history.forEach((h) => {
      const wk = getWeekKey(h.ts)
      if (wk) set.add(wk)
    })
    return Array.from(set).sort().reverse().slice(0, 12)
  }, [history])

  const filteredAndSorted = useMemo(() => {
    let list = [...history]
    if (period === 'day') {
      list = list.filter((h) => getDateKey(h.ts) === dayKey)
    } else if (period === 'week') {
      list = list.filter((h) => getWeekKey(h.ts) === weekKey)
    }
    if (sortOrder === 'oldest') list.reverse()
    return list
  }, [history, period, dayKey, weekKey, sortOrder])

  const periodOnlineSec = useMemo(() => {
    if (period === 'all') return stats.totalOnlineSec
    if (period === 'day') return stats.byDay[dayKey] ?? 0
    return stats.byWeek[weekKey] ?? 0
  }, [stats, period, dayKey, weekKey])

  if (error) {
    return (
      <section className="info-card info-card-history">
        <h3>История онлайна</h3>
        <p className="history-error">{error}</p>
      </section>
    )
  }

  return (
    <section className="info-card info-card-history">
      <h3>История онлайна</h3>
      <p className="history-hint">Только смена онлайн/оффлайн. Собирается для отслеживаемых.</p>

      {loading ? (
        <p className="loading-small">Загрузка...</p>
      ) : history.length === 0 ? (
        <p className="history-empty">Нет записей. Добавьте игрока в отслеживание для сбора истории.</p>
      ) : (
        <>
          <div className="history-stats-row">
            <span className="history-total">Общий онлайн: {formatDuration(stats.totalOnlineSec) || '0'}</span>
          </div>

          {showHeatmap && <ActivityHeatmap byDay={stats.byDay} history={history} />}

          <div className="history-controls">
            <label className="history-control-label">
              Период:
              <select
                value={period}
                onChange={(e) => setPeriod(e.target.value as PeriodFilter)}
                className="history-select"
              >
                <option value="all">Все</option>
                <option value="day">День</option>
                <option value="week">Неделя</option>
              </select>
            </label>
            {period === 'day' && (
              <label className="history-control-label">
                День:
                <select
                  value={dayKey}
                  onChange={(e) => setDayKey(e.target.value)}
                  className="history-select"
                >
                  {dayKeys.length === 0 && <option value={dayKey}>{dayKey}</option>}
                  {dayKeys.map((dk) => (
                    <option key={dk} value={dk}>
                      {dk}
                    </option>
                  ))}
                </select>
              </label>
            )}
            {period === 'week' && (
              <label className="history-control-label">
                Неделя с:
                <select
                  value={weekKey}
                  onChange={(e) => setWeekKey(e.target.value)}
                  className="history-select"
                >
                  {weekKeys.length === 0 && <option value={weekKey}>{weekKey}</option>}
                  {weekKeys.map((wk) => (
                    <option key={wk} value={wk}>
                      {wk}
                    </option>
                  ))}
                </select>
              </label>
            )}
            {(period === 'day' || period === 'week') && (
              <span className="history-period-total">
                За выбранный период: {formatDuration(periodOnlineSec) || '0'}
              </span>
            )}
            <label className="history-control-label">
              Сортировка:
              <select
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value as SortOrder)}
                className="history-select"
              >
                <option value="newest">Сначала новые</option>
                <option value="oldest">Сначала старые</option>
              </select>
            </label>
          </div>

          <div className="online-graph-inline">
            {[...history].reverse().slice(-80).map((h, i) => (
              <div
                key={i}
                className={`graph-bar-inline ${h.online ? 'online' : 'offline'}`}
                title={`${h.ts} — ${h.online ? 'online' : 'offline'}`}
              />
            ))}
          </div>

          <ul className={`history-log ${compact ? 'history-log-compact' : ''}`}>
            {filteredAndSorted.slice(0, compact ? 40 : 80).map((h, i) => (
              <li key={i} className={h.online ? 'history-online' : 'history-offline'}>
                {formatHistoryLine(h)}
              </li>
            ))}
          </ul>
        </>
      )}
    </section>
  )
}
