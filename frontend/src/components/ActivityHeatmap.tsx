import React, { useMemo, useState, useRef } from 'react'
import { formatDuration, getDaySessionsFromHistory, type HistoryRecord } from '../api/client'
import './ActivityHeatmap.css'

const WEEKS = 53
const DAYS_PER_WEEK = 7
const TOTAL_DAYS = WEEKS * DAYS_PER_WEEK // 371

/** Уровень заливки по секундам: 0 — нет, 1–4 — по возрастанию */
function getLevel(sec: number): 0 | 1 | 2 | 3 | 4 {
  if (sec <= 0) return 0
  if (sec < 30 * 60) return 1       // < 30 мин
  if (sec < 60 * 60) return 2       // < 1 ч
  if (sec < 2 * 60 * 60) return 3   // < 2 ч
  return 4                           // 2+ ч
}

function formatDayLabel(date: Date): string {
  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

export function ActivityHeatmap({
  byDay,
  history,
}: {
  byDay: Record<string, number>
  /** Если передан — в тултипе показываются серверы и временные промежутки по дням */
  history?: HistoryRecord[]
}) {
  const [tooltip, setTooltip] = useState<{ content: React.ReactNode; x: number; y: number } | null>(null)
  const wrapRef = useRef<HTMLDivElement>(null)

  const daySessions = useMemo(
    () => (history ? getDaySessionsFromHistory(history) : {}),
    [history]
  )

  const grid = useMemo(() => {
    const today = new Date()
    today.setUTCHours(0, 0, 0, 0)
    const start = new Date(today)
    start.setUTCDate(start.getUTCDate() - TOTAL_DAYS + 1)

    const out: { dateKey: string; sec: number; level: 0 | 1 | 2 | 3 | 4; date: Date }[][] = []
    for (let col = 0; col < WEEKS; col++) {
      const column: { dateKey: string; sec: number; level: 0 | 1 | 2 | 3 | 4; date: Date }[] = []
      for (let row = 0; row < DAYS_PER_WEEK; row++) {
        const dayIndex = col * DAYS_PER_WEEK + row
        const cellDate = new Date(start)
        cellDate.setUTCDate(cellDate.getUTCDate() + dayIndex)
        const dateKey = cellDate.toISOString().slice(0, 10)
        const sec = byDay[dateKey] ?? 0
        const level = getLevel(sec)
        column.push({ dateKey, sec, level, date: cellDate })
      }
      out.push(column)
    }
    return out
  }, [byDay])

  return (
    <div ref={wrapRef} className="activity-heatmap-wrap">
      <div className="activity-heatmap">
        {grid.map((column, colIndex) => (
          <div key={colIndex} className="heatmap-column">
            {column.map((cell) => (
              <div
                key={cell.dateKey}
                className={`heatmap-cell level-${cell.level}`}
                onMouseEnter={(e) => {
                  const wrap = wrapRef.current
                  if (!wrap) return
                  const rect = e.currentTarget.getBoundingClientRect()
                  const wrapRect = wrap.getBoundingClientRect()
                  const sessions = daySessions[cell.dateKey]
                  const content = (
                    <div className="heatmap-tooltip-content">
                      <div className="heatmap-tooltip-title">
                        {formatDayLabel(cell.date)} — {cell.sec > 0 ? formatDuration(cell.sec) : 'нет активности'}
                      </div>
                      {sessions && sessions.length > 0 && (
                        <div className="heatmap-tooltip-sessions">
                          {sessions.map((s, i) => (
                            <div key={i} className="heatmap-tooltip-session">
                              <span className="heatmap-tooltip-server">{s.server}</span>
                              <span className="heatmap-tooltip-time">{s.start} – {s.end}</span>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  )
                  setTooltip({
                    content,
                    x: rect.left - wrapRect.left + rect.width / 2,
                    y: rect.top - wrapRect.top,
                  })
                }}
                onMouseLeave={() => setTooltip(null)}
                title={cell.sec > 0 ? formatDuration(cell.sec) : 'Нет активности'}
              />
            ))}
          </div>
        ))}
      </div>
      {tooltip && (
        <div
          className="heatmap-tooltip"
          style={{
            left: tooltip.x,
            top: tooltip.y,
            transform: 'translate(-50%, -100%) translateY(-6px)',
          }}
          aria-hidden
        >
          {tooltip.content}
        </div>
      )}
      <div className="activity-heatmap-legend">
        <span className="legend-label">Меньше</span>
        <span className="legend-cell level-0" />
        <span className="legend-cell level-1" />
        <span className="legend-cell level-2" />
        <span className="legend-cell level-3" />
        <span className="legend-cell level-4" />
        <span className="legend-label">Больше</span>
      </div>
    </div>
  )
}
