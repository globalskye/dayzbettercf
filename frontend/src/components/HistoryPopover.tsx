import { useEffect, useRef, useState } from 'react'
import { fetchPlayerHistory, formatHistoryLine, type HistoryRecord } from '../api/client'
import './HistoryPopover.css'

export function HistoryPopover({
  cftoolsId,
  displayName,
  children,
  onOpenChange,
}: {
  cftoolsId: string
  displayName?: string
  children: React.ReactNode
  onOpenChange?: (open: boolean) => void
}) {
  const [open, setOpen] = useState(false)
  const [history, setHistory] = useState<HistoryRecord[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const popRef = useRef<HTMLDivElement>(null)
  const triggerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open || !cftoolsId) return
    onOpenChange?.(true)
    setLoading(true)
    setError(null)
    fetchPlayerHistory(cftoolsId, 150)
      .then((r) => setHistory(Array.isArray(r?.history) ? r.history : []))
      .catch((e) => setError(e?.message ?? 'Не удалось загрузить историю'))
      .finally(() => setLoading(false))
  }, [open, cftoolsId, onOpenChange])

  useEffect(() => {
    if (!open) return
    const close = (e: MouseEvent) => {
      if (
        popRef.current?.contains(e.target as Node) ||
        triggerRef.current?.contains(e.target as Node)
      ) return
      setOpen(false)
      onOpenChange?.(false)
    }
    document.addEventListener('mousedown', close)
    return () => document.removeEventListener('mousedown', close)
  }, [open, onOpenChange])

  const toggle = () => {
    setOpen((v) => !v)
    if (!open) onOpenChange?.(true)
    else onOpenChange?.(false)
  }

  const sorted = [...history].sort(
    (a, b) => new Date(b.ts).getTime() - new Date(a.ts).getTime()
  )

  return (
    <div className="history-popover-wrap">
      <div
        ref={triggerRef}
        role="button"
        tabIndex={0}
        className="history-popover-trigger"
        onClick={toggle}
        onKeyDown={(e) => e.key === 'Enter' && toggle()}
        title="История онлайна (клик)"
      >
        {children}
      </div>
      {open && (
        <div ref={popRef} className="history-popover-dropdown">
          <div className="history-popover-header">
            <span className="history-popover-title">
              История — {displayName || cftoolsId}
            </span>
            <button
              type="button"
              className="history-popover-close"
              onClick={() => { setOpen(false); onOpenChange?.(false) }}
              title="Закрыть"
            >
              ×
            </button>
          </div>
          <div className="history-popover-body">
            {loading ? (
              <p className="history-popover-loading">Загрузка...</p>
            ) : error ? (
              <p className="history-popover-error">{error}</p>
            ) : sorted.length === 0 ? (
              <p className="history-popover-empty">Нет записей</p>
            ) : (
              <ul className="history-popover-list">
                {sorted.map((h, i) => (
                  <li
                    key={`${h.ts}-${i}`}
                    className={h.online ? 'history-popover-online' : 'history-popover-offline'}
                  >
                    {formatHistoryLine(h)}
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
