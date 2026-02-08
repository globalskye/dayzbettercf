import { useState } from 'react'
import {
  addToGroup,
  addTracked,
  removeTracked,
  type Player,
  type Group,
} from '../api/client'
import '../pages/HomePage.css'

export function PlayerActionsModal({
  player,
  groups,
  isTracked,
  onClose,
  onSuccess,
}: {
  player: Player
  groups: Group[]
  isTracked: boolean
  onClose: () => void
  onSuccess: () => void
}) {
  const [groupId, setGroupId] = useState<number | null>(null)
  const [alias, setAlias] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleAddToGroup = () => {
    if (!groupId) return
    setLoading(true)
    setError(null)
    addToGroup(groupId, player.cftools_id, alias.trim() || undefined)
      .then(() => {
        setGroupId(null)
        setAlias('')
        onSuccess()
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  const handleToggleTracked = () => {
    setLoading(true)
    setError(null)
    const fn = isTracked ? removeTracked : addTracked
    fn(player.cftools_id)
      .then(onSuccess)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>{player.display_name}</h3>
          <button type="button" className="modal-close" onClick={onClose}>×</button>
        </div>
        {error && <div className="error-msg">{error}</div>}
        <div className="modal-body">
          <div className="modal-section">
            <h4>Добавить в группу</h4>
            <div className="modal-row">
              <select
                value={groupId ?? ''}
                onChange={(e) => setGroupId(e.target.value ? Number(e.target.value) : null)}
              >
                <option value="">Выбери группу...</option>
                {groups.map((g) => (
                  <option key={g.id} value={g.id}>{g.name}</option>
                ))}
              </select>
              <input
                type="text"
                placeholder="Алиас (как запомнил)"
                value={alias}
                onChange={(e) => setAlias(e.target.value)}
              />
              <button
                type="button"
                onClick={handleAddToGroup}
                disabled={!groupId || loading}
              >
                Добавить
              </button>
            </div>
          </div>
          <div className="modal-section">
            <h4>Отслеживание</h4>
            <button
              type="button"
              onClick={handleToggleTracked}
              disabled={loading}
              className={isTracked ? 'btn-tracked' : ''}
            >
              {isTracked ? 'Убрать из отслеживания' : 'Добавить в отслеживание'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
