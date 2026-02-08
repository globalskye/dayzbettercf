import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { usePlayerModal } from '../context/PlayerModalContext'
import {
  fetchGroups,
  fetchGroup,
  createGroup,
  deleteGroup,
  removeFromGroup,
  updateMemberAlias,
  syncBatch,
  formatDate,
  type Group,
  type Member,
} from '../api/client'
import { LoadingSpinner, LoadingOverlay } from '../components/LoadingSpinner'
import './GroupsPage.css'

export function GroupsPage() {
  const openPlayerModal = usePlayerModal()
  const [groups, setGroups] = useState<Group[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [newName, setNewName] = useState('')
  const [adding, setAdding] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [membersSort, setMembersSort] = useState<'online' | 'playtime' | 'bans' | 'name'>('online')

  const load = () => {
    setLoading(true)
    fetchGroups(membersSort, false)
      .then((r) => setGroups(Array.isArray(r?.groups) ? r.groups : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  const [loadingOnlineId, setLoadingOnlineId] = useState<number | null>(null)
  const [refreshingData, setRefreshingData] = useState(false)
  const refreshGroupOnline = (groupId: number) => {
    setLoadingOnlineId(groupId)
    fetchGroup(groupId, membersSort)
      .then((g) => setGroups((prev) => prev.map((gr) => (gr.id === g.id ? g : gr))))
      .catch((e) => setError(e.message))
      .finally(() => setLoadingOnlineId(null))
  }
  const handleRefreshGroupData = () => {
    if (!selected) return
    const ids = (selected.members ?? []).map((m) => m.cftools_id).filter(Boolean)
    if (ids.length === 0) return
    setRefreshingData(true)
    setError(null)
    syncBatch(ids, false)
      .then(() => refreshGroupOnline(selected.id))
      .catch((e) => setError(e.message))
      .finally(() => setRefreshingData(false))
  }

  useEffect(() => {
    load()
  }, [membersSort])

  // Как в old: при открытии группы — один запрос за онлайном по этой группе
  useEffect(() => {
    if (selectedId == null) return
    refreshGroupOnline(selectedId)
  }, [selectedId, membersSort])

  const handleCreate = () => {
    if (!newName.trim()) return
    setAdding(true)
    setError(null)
    createGroup(newName.trim())
      .then((g) => {
        setGroups((prev) => [...prev, g])
        setNewName('')
        setSelectedId(g.id)
      })
      .catch((e) => setError(e.message))
      .finally(() => setAdding(false))
  }

  const handleDelete = (id: number) => {
    if (!confirm('Удалить группу?')) return
    deleteGroup(id)
      .then(() => {
        setGroups((prev) => prev.filter((g) => g.id !== id))
        if (selectedId === id) setSelectedId(null)
      })
      .catch((e) => setError(e.message))
  }

  const selected = (groups ?? []).find((g) => g.id === selectedId)

  return (
    <div className="groups-page">
      <Link to="/" className="back-link">← Назад</Link>
      {loading ? (
        <LoadingOverlay text="Загрузка групп..." />
      ) : (
        <>
      <div className="groups-toolbar">
        <input
          type="text"
          className="groups-new-input"
          placeholder="Название новой группы"
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
        />
        <button onClick={handleCreate} disabled={adding || !newName.trim()} className="groups-new-btn">
          {adding ? (
            <span className="btn-loading-wrap">
              <LoadingSpinner variant="dots" size="small" />
              <span>Создание...</span>
            </span>
          ) : (
            'Создать группу'
          )}
        </button>
        <select
          value={membersSort}
          onChange={(e) => setMembersSort(e.target.value as typeof membersSort)}
          className="groups-sort-select"
        >
          <option value="online">Сорт: онлайн</option>
          <option value="playtime">Сорт: часы</option>
          <option value="bans">Сорт: баны</option>
          <option value="name">Сорт: имя</option>
        </select>
      </div>
      <div className="groups-layout">
        <aside className="groups-sidebar">
          <ul className="groups-list">
            {(groups ?? []).map((g) => (
              <li key={g.id} className={selectedId === g.id ? 'active' : ''}>
                <button
                  type="button"
                  className="group-btn"
                  onClick={() => setSelectedId(g.id)}
                >
                  {g.name}
                  <span className="group-count">({g.members?.length ?? 0})</span>
                </button>
                <button
                  type="button"
                  className="group-delete-btn"
                  onClick={(e) => {
                    e.stopPropagation()
                    handleDelete(g.id)
                  }}
                  title="Удалить группу"
                >
                  ×
                </button>
              </li>
            ))}
          </ul>
        </aside>
        <main className="groups-main">
          {error && <div className="error-msg">{error}</div>}
          {selected ? (
            <GroupDetail
              group={selected}
              loadingOnline={loadingOnlineId === selected.id}
              refreshingData={refreshingData}
              onRefresh={load}
              onRefreshOnline={() => refreshGroupOnline(selected.id)}
              onRefreshData={handleRefreshGroupData}
              onOpenPlayer={openPlayerModal ?? undefined}
            />
          ) : (
            <p className="groups-empty">Выбери группу или создай новую</p>
          )}
        </main>
      </div>
        </>
      )}
    </div>
  )
}

function GroupDetail({
  group,
  loadingOnline,
  refreshingData,
  onRefresh,
  onRefreshOnline,
  onRefreshData,
  onOpenPlayer,
}: {
  group: Group
  loadingOnline?: boolean
  refreshingData?: boolean
  onRefresh: () => void
  onRefreshOnline: () => void
  onRefreshData: () => void
  onOpenPlayer?: (id: string) => void
}) {
  const members = Array.isArray(group?.members) ? group.members : []
  const onlineCount = members.filter((m) => m?.player?.online).length
  const anyLoading = loadingOnline || refreshingData

  return (
    <>
      <div className="group-header">
        <div className="group-header-top">
          <h2>{group.name}</h2>
          <button
            type="button"
            className="btn-online-refresh"
            onClick={onRefreshOnline}
            disabled={anyLoading || members.length === 0}
            title="Обновить онлайн из CF"
          >
            {loadingOnline ? (
              <span className="btn-loading-wrap">
                <LoadingSpinner variant="dots" size="small" />
                <span>Обновление...</span>
              </span>
            ) : (
              'Обновить онлайн'
            )}
          </button>
          <button
            type="button"
            className="btn-data-refresh"
            onClick={onRefreshData}
            disabled={anyLoading || members.length === 0}
            title="Подтянуть из CF ники, время игры, сервер и т.д."
          >
            {refreshingData ? (
              <span className="btn-loading-wrap">
                <LoadingSpinner variant="dots" size="small" />
                <span>Обновление данных...</span>
              </span>
            ) : (
              'Обновить последние данные'
            )}
          </button>
        </div>
        <p className="group-stats">
          {members.length > 0 ? `${onlineCount} / ${members.length} онлайн` : 'Нет участников'}
        </p>
      </div>
      <div className="group-members">
        {members.length > 0 ? (
          <table className="members-table">
            <thead>
              <tr>
                <th>Имя</th>
                <th>Последний ник</th>
                <th>Последняя игра</th>
                <th>Сервер</th>
                <th className="member-cell-actions" />
              </tr>
            </thead>
            <tbody>
              {members.map((m) => (
                <MemberRow
                  key={m.player_id}
                  member={m}
                  groupId={group.id}
                  onRefresh={onRefresh}
                  onOpenPlayer={onOpenPlayer}
                />
              ))}
            </tbody>
          </table>
        ) : (
          <p className="empty-hint">Добавь игроков через поиск — кнопка «В группу»</p>
        )}
      </div>
    </>
  )
}

function MemberRow({
  member,
  groupId,
  onRefresh,
  onOpenPlayer,
}: {
  member: Member
  groupId: number
  onRefresh: () => void
  onOpenPlayer?: (id: string) => void
}) {
  const [removing, setRemoving] = useState(false)
  const [editingAlias, setEditingAlias] = useState(false)
  const [aliasInput, setAliasInput] = useState(member.alias ?? '')
  const [savingAlias, setSavingAlias] = useState(false)
  const p = member.player

  const displayName = (member.alias && member.alias.trim()) ? member.alias.trim() : (p?.display_name ?? '')
  const nick = (p?.nicknames?.length ? p.nicknames[p.nicknames.length - 1] : null) ?? ''
  const server = p?.last_server_identifier ?? ''

  const handleRemove = () => {
    setRemoving(true)
    removeFromGroup(groupId, member.cftools_id)
      .then(onRefresh)
      .finally(() => setRemoving(false))
  }

  const handleEditAlias = () => {
    setEditingAlias(true)
    setAliasInput(member.alias ?? '')
  }

  const handleSaveAlias = () => {
    setSavingAlias(true)
    updateMemberAlias(groupId, member.cftools_id, aliasInput.trim())
      .then(() => {
        onRefresh()
        setEditingAlias(false)
      })
      .catch(() => {})
      .finally(() => setSavingAlias(false))
  }

  const handleCancelAlias = () => {
    setEditingAlias(false)
    setAliasInput(member.alias ?? '')
  }

  const lastPlaytime = p?.updated_at ? formatDate(p.updated_at) : '—'
  const serverLabel = p?.online ? (server || 'Online') : (server || 'Offline')

  return (
    <tr className="member-row">
      <td className="member-cell-name">
        <button
          type="button"
          className="member-link"
          onClick={() => onOpenPlayer?.(member.cftools_id)}
        >
          {p?.avatar ? (
            <img src={p.avatar} alt="" className="member-avatar" />
          ) : (
            <div className="member-avatar-placeholder" />
          )}
          <div className="member-info">
            {editingAlias ? (
              <div className="member-alias-edit" onClick={(e) => e.stopPropagation()}>
                <input
                  type="text"
                  value={aliasInput}
                  onChange={(e) => setAliasInput(e.target.value)}
                  placeholder="Алиас"
                  autoFocus
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') handleSaveAlias()
                    if (e.key === 'Escape') handleCancelAlias()
                  }}
                />
                <button type="button" onClick={handleSaveAlias} disabled={savingAlias}>✓</button>
                <button type="button" onClick={handleCancelAlias}>✕</button>
              </div>
            ) : (
              <>
                <span className="member-name">{displayName}</span>
                <button
                  type="button"
                  className="member-edit-alias"
                  onClick={(e) => {
                    e.stopPropagation()
                    handleEditAlias()
                  }}
                  title="Изменить алиас"
                >
                  ✎
                </button>
              </>
            )}
          </div>
        </button>
      </td>
      <td className="member-cell-nick">
        {nick || '—'}
      </td>
      <td className="member-cell-playtime">
        {lastPlaytime}
      </td>
      <td className="member-cell-server">
        {p?.online && <span className="badge online">online</span>}
        <span className="member-server" title={server || ''}>{serverLabel}</span>
      </td>
      <td className="member-cell-actions">
        <button
          type="button"
          className="member-remove"
          onClick={handleRemove}
          disabled={removing}
          title="Удалить из группы"
        >
          ×
        </button>
      </td>
    </tr>
  )
}
