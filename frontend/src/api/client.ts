const API_BASE = '/api/v1';

const TOKEN_KEY = 'dayzsmartcf_token';

function authHeaders(): HeadersInit {
  const t = typeof localStorage !== 'undefined' ? localStorage.getItem(TOKEN_KEY) : null;
  const h: HeadersInit = { 'Content-Type': 'application/json' };
  if (t) (h as Record<string, string>)['Authorization'] = `Bearer ${t}`;
  return h;
}

async function apiFetch(url: string, init?: RequestInit): Promise<Response> {
  return fetch(url, {
    ...init,
    headers: { ...authHeaders(), ...init?.headers },
  });
}

export interface Player {
  id: number;
  cftools_id: string;
  display_name: string;
  avatar?: string;
  is_bot: boolean;
  account_status: number;
  playtime_sec: number;
  sessions_count: number;
  bans_count: number;
  linked_accounts_count: number;
  last_activity_at?: string;
  last_seen_at?: string;
  online: boolean;
  raw_status?: string;
  raw_overview?: string;
  raw_structure?: string;
  raw_play_state?: string;
  raw_bans?: string;
  raw_battleye?: string;
  steam64?: string;
  steam_avatar?: string;
  steam_persona?: string;
  steam_vac_bans?: number;
  steam_game_bans?: number;
  created_at: string;
  updated_at: string;
  nicknames?: string[];
  /** Связанные аккаунты с флагами: confirmed = CFTools подтвердил связь, trusted = отмечен доверенным */
  linked_accounts?: { cftools_id: string; confirmed: boolean; trusted: boolean }[];
  linked_cftools_ids?: string[];
  server_ids?: string[];
  last_server_identifier?: string;
}

export interface Group {
  id: number;
  name: string;
  members?: Member[];
  created_at: string;
  updated_at: string;
}

export interface Member {
  group_id: number;
  player_id: number;
  cftools_id: string;
  alias: string;
  player?: Player;
  created_at: string;
}

export interface PlayersResponse {
  players: Player[];
  total?: number;
  count?: number;
}

/** Результат GET /cftools/states — только данные из CF, без записи в БД */
export interface ProfileStateItem {
  cftools_id: string;
  display_name: string;
  avatar?: string;
  online: boolean;
  server_name?: string;
}

export interface CFToolsStatesResponse {
  states: ProfileStateItem[];
  count: number;
}

export interface HealthResponse {
  status: string;
}

// App auth
export interface LoginResponse {
  token: string;
  user: { id: number; username: string; role: string };
  expires_in: number;
}

export async function login(username: string, password: string): Promise<LoginResponse> {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error((err as { error?: string }).error || 'Login failed');
  }
  return res.json();
}

export async function fetchAuthMe(token: string): Promise<{ id: number; username: string; role: string }> {
  const res = await fetch(`${API_BASE}/auth/me`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error('Not authenticated');
  return res.json();
}

export async function fetchHealth(): Promise<HealthResponse> {
  const res = await fetch('/health');
  if (!res.ok) throw new Error('API offline');
  return res.json();
}

export type PlayersFilters = {
  limit?: number
  offset?: number
  online?: boolean
  banned?: boolean
  sort?: 'online' | 'playtime' | 'bans' | 'updated'
}

export async function fetchPlayers(filters?: PlayersFilters): Promise<PlayersResponse> {
  const params = new URLSearchParams()
  params.set('limit', String(filters?.limit ?? 50))
  params.set('offset', String(filters?.offset ?? 0))
  if (filters?.online) params.set('online', '1')
  if (filters?.banned) params.set('banned', '1')
  if (filters?.sort) params.set('sort', filters.sort)
  const res = await apiFetch(`${API_BASE}/players?${params}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function searchPlayersLocal(q: string, filters?: PlayersFilters): Promise<PlayersResponse> {
  const params = new URLSearchParams({ q })
  if (filters?.limit) params.set('limit', String(filters.limit))
  if (filters?.online) params.set('online', '1')
  if (filters?.banned) params.set('banned', '1')
  if (filters?.sort) params.set('sort', filters.sort)
  const res = await apiFetch(`${API_BASE}/players/search?${params}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export interface CFtoolsSearchResult {
  cftools_id: string;
  display_name: string;
  avatar?: string;
  identifier?: string;
}

export async function searchCFtools(q: string): Promise<{ results: CFtoolsSearchResult[]; count: number }> {
  const res = await apiFetch(`${API_BASE}/players/cftools-search?q=${encodeURIComponent(q)}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

/** Запрос состояний по CF (GlobalQuery + playState), без записи в БД — как в old. */
export async function getProfileStates(q: string): Promise<CFToolsStatesResponse> {
  const res = await apiFetch(`${API_BASE}/cftools/states?q=${encodeURIComponent(q)}`);
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error((err as { error?: string }).error || 'CF states failed');
  }
  return res.json();
}

/** Поиск по CF с синхронизацией в базу (search-cf). Используется для кнопки «В базу». */
export async function searchByCF(q: string, light: boolean): Promise<PlayersResponse> {
  const params = new URLSearchParams({ q });
  params.set('light', light ? '1' : '0');
  const res = await apiFetch(`${API_BASE}/players/search-cf?${params}`);
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error((err as { error?: string }).error || 'Search failed');
  }
  return res.json();
}

export async function syncBatch(cftoolsIds: string[], light = true): Promise<PlayersResponse> {
  const params = light ? '?light=1' : '';
  const res = await apiFetch(`${API_BASE}/players/sync-batch${params}`, {
    method: 'POST',
    body: JSON.stringify({ cftools_ids: cftoolsIds }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error((err as { error?: string }).error || 'Sync failed');
  }
  return res.json();
}

export async function fetchPlayer(id: string): Promise<Player> {
  const res = await apiFetch(`${API_BASE}/players/${encodeURIComponent(id)}`);
  if (!res.ok) {
    if (res.status === 404) throw new Error('Player not found');
    throw new Error(await res.text());
  }
  return res.json();
}

/** light=true — только статус/онлайн/overview; light=false — полный профиль (Steam, баны и т.д.) */
export async function syncPlayer(id: string, light = true): Promise<Player> {
  const q = light ? '?light=1' : '';
  const res = await apiFetch(`${API_BASE}/players/${encodeURIComponent(id)}/sync${q}`, {
    method: 'POST',
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export function formatPlaytime(sec?: number): string {
  const s = sec ?? 0;
  if (s < 60) return `${s}s`;
  if (s < 3600) return `${Math.floor(s / 60)}m`;
  if (s < 86400) return `${(s / 3600).toFixed(1)}h`;
  return `${(s / 86400).toFixed(1)}d`;
}

export function formatDate(s?: string): string {
  if (!s) return '—';
  try {
    const d = new Date(s);
    return d.toLocaleString();
  } catch {
    return s;
  }
}

export interface AuthSettings {
  configured: boolean
}

export interface AuthSettingsBody {
  cdn_auth: string
  cf_clearance?: string
  session?: string
  user_info?: string
  acsrf?: string
}

export async function fetchAuthSettings(): Promise<AuthSettings> {
  const res = await apiFetch(`${API_BASE}/settings/auth`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export interface AuthCheckResult {
  ok: boolean
  error?: string
}

export async function checkAuth(): Promise<AuthCheckResult> {
  const res = await apiFetch(`${API_BASE}/settings/auth/check`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function updateAuthSettings(body: AuthSettingsBody): Promise<void> {
  const res = await apiFetch(`${API_BASE}/settings/auth`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error((err as { error?: string }).error || 'Update failed')
  }
}

/** Очистить все данные базы (игроки, группы, история). Только админ. Учётные записи не удаляются. */
export async function wipeDatabase(): Promise<{ message: string }> {
  const res = await apiFetch(`${API_BASE}/settings/db/wipe`, { method: 'POST' })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error((err as { error?: string }).error || 'Очистка не удалась')
  }
  return res.json()
}

// --- Admin: users & request logs ---
export interface AdminUser {
  id: number
  username: string
  role: string
  created_at: string
}

export interface RequestLogEntry {
  id: number
  user_id: number
  method: string
  path: string
  created_at: string
}

export async function fetchAdminUsers(): Promise<{ users: AdminUser[] }> {
  const res = await apiFetch(`${API_BASE}/admin/users`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function adminCreateUser(username: string, password: string, role: string): Promise<{ users: AdminUser[] }> {
  const res = await apiFetch(`${API_BASE}/admin/users`, {
    method: 'POST',
    body: JSON.stringify({ username, password, role }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error((err as { error?: string }).error || 'Ошибка создания')
  }
  return res.json()
}

export async function adminUpdateUser(id: number, data: { role?: string; password?: string }): Promise<{ users: AdminUser[] }> {
  const res = await apiFetch(`${API_BASE}/admin/users/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error((err as { error?: string }).error || 'Ошибка обновления')
  }
  return res.json()
}

export async function adminDeleteUser(id: number): Promise<void> {
  const res = await apiFetch(`${API_BASE}/admin/users/${id}`, { method: 'DELETE' })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error((err as { error?: string }).error || 'Ошибка удаления')
  }
}

export async function fetchUserRequestLogs(userId: number, limit?: number): Promise<{ logs: RequestLogEntry[] }> {
  const q = limit != null ? `?limit=${limit}` : ''
  const res = await apiFetch(`${API_BASE}/admin/users/${userId}/logs${q}`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

/** Список групп. enrich=false — без запросов онлайна в CF (быстро). */
export async function fetchGroups(sort?: string, enrich = false): Promise<{ groups: Group[] }> {
  const params = new URLSearchParams()
  if (sort) params.set('sort', sort)
  if (!enrich) params.set('enrich', '0')
  const url = `${API_BASE}/groups${params.toString() ? `?${params}` : ''}`
  const res = await apiFetch(url)
  if (!res.ok) throw new Error(await res.text())
  const data = await res.json()
  return { groups: Array.isArray(data?.groups) ? data.groups : [] }
}

/** Одна группа с актуальным онлайном из CF (для кнопки «Просмотр онлайна»). */
export async function fetchGroup(id: number, sort?: string): Promise<Group> {
  const params = sort ? `?sort=${encodeURIComponent(sort)}` : ''
  const res = await apiFetch(`${API_BASE}/groups/${id}${params}`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function createGroup(name: string): Promise<Group> {
  const res = await apiFetch(`${API_BASE}/groups/create/${encodeURIComponent(name)}`, { method: 'POST' })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function deleteGroup(id: number): Promise<void> {
  const res = await apiFetch(`${API_BASE}/groups/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(await res.text())
}

export async function addToGroup(groupId: number, cftoolsId: string, alias?: string): Promise<Group> {
  const url = `${API_BASE}/groups/${groupId}/add/${encodeURIComponent(cftoolsId)}` + (alias ? `?alias=${encodeURIComponent(alias)}` : '')
  const res = await apiFetch(url, { method: 'POST' })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function updateMemberAlias(groupId: number, cftoolsId: string, alias: string): Promise<Group> {
  const res = await apiFetch(`${API_BASE}/groups/${groupId}/members/${encodeURIComponent(cftoolsId)}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ alias }),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function removeFromGroup(groupId: number, cftoolsId: string): Promise<void> {
  const res = await apiFetch(`${API_BASE}/groups/${groupId}/remove/${encodeURIComponent(cftoolsId)}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(await res.text())
}

// Tracked players (до 10, постоянное обновление статусов)
export async function fetchTracked(sort?: string): Promise<{ players: Player[] }> {
  const url = sort ? `${API_BASE}/tracked?sort=${encodeURIComponent(sort)}` : `${API_BASE}/tracked`
  const res = await apiFetch(url)
  if (!res.ok) throw new Error(await res.text())
  const data = await res.json()
  return { players: Array.isArray(data?.players) ? data.players : [] }
}

export async function addTracked(cftoolsId: string): Promise<{ players: Player[] }> {
  const res = await apiFetch(`${API_BASE}/tracked/add/${encodeURIComponent(cftoolsId)}`, { method: 'POST' })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error((err as { error?: string }).error || 'Add failed')
  }
  return res.json()
}

export async function removeTracked(cftoolsId: string): Promise<void> {
  const res = await apiFetch(`${API_BASE}/tracked/remove/${encodeURIComponent(cftoolsId)}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(await res.text())
}

export interface HistoryRecord {
  ts: string
  online: boolean
  server_name?: string
  playtime_sec: number
  sessions_count: number
  display_name?: string
  /** Длительность сессии при уходе оффлайн (сек) */
  session_duration_sec?: number
  /** Сколько был оффлайн перед этим заходом (сек) */
  offline_duration_sec?: number
}

/** Форматирует длительность для лога: "45 мин", "2ч 30мин" */
export function formatDuration(sec: number): string {
  if (sec <= 0) return ''
  if (sec < 60) return `${sec} сек`
  if (sec < 3600) return `${Math.floor(sec / 60)} мин`
  const h = Math.floor(sec / 3600)
  const m = Math.floor((sec % 3600) / 60)
  if (m === 0) return `${h} ч`
  return `${h}ч ${m}мин`
}

/** Одна строка лога онлайна: дата, играет online сервер да / вышел, сессия N */
export function formatHistoryLine(h: HistoryRecord): string {
  const dateStr = formatDate(h.ts)
  if (h.online) {
    const server = h.server_name?.trim() ? ` ${h.server_name}` : ''
    const offlinePart = (h.offline_duration_sec ?? 0) > 0 ? ` (был оффлайн ${formatDuration(h.offline_duration_sec!)})` : ''
    return `${dateStr} — играет online${server} да${offlinePart}`
  }
  const sessionPart = (h.session_duration_sec ?? 0) > 0 ? `, сессия ${formatDuration(h.session_duration_sec!)}` : ''
  return `${dateStr} — вышел${sessionPart}`
}

/** Ключ дня YYYY-MM-DD из ts */
export function getDateKey(ts: string): string {
  try {
    const d = new Date(ts)
    return d.toISOString().slice(0, 10)
  } catch {
    return ''
  }
}

/** По истории строит по дням список сессий: сервер и интервал (HH:mm–HH:mm). День = день окончания сессии (offline). Текущая сессия (online) добавляется на сегодня с концом "сейчас". */
export function getDaySessionsFromHistory(history: HistoryRecord[]): Record<string, { server: string; start: string; end: string }[]> {
  const daySessions: Record<string, { server: string; start: string; end: string }[]> = {}
  const sorted = [...history].sort((a, b) => new Date(a.ts).getTime() - new Date(b.ts).getTime())
  let lastOnline: { server: string; ts: string } | null = null
  for (const h of sorted) {
    if (h.online) {
      lastOnline = { server: h.server_name?.trim() || '—', ts: h.ts }
    } else if (lastOnline) {
      const dayKey = getDateKey(h.ts)
      if (!dayKey) continue
      if (!daySessions[dayKey]) daySessions[dayKey] = []
      const startTime = new Date(lastOnline.ts)
      const endTime = new Date(h.ts)
      const fmt = (d: Date) => d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit', hour12: false })
      daySessions[dayKey].push({
        server: lastOnline.server,
        start: fmt(startTime),
        end: fmt(endTime),
      })
      lastOnline = null
    }
  }
  if (lastOnline) {
    const dayKey = getDateKey(new Date().toISOString())
    if (dayKey) {
      if (!daySessions[dayKey]) daySessions[dayKey] = []
      const fmt = (d: Date) => d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit', hour12: false })
      daySessions[dayKey].push({
        server: lastOnline!.server,
        start: fmt(new Date(lastOnline!.ts)),
        end: 'сейчас',
      })
    }
  }
  return daySessions
}

/** Понедельник недели в YYYY-MM-DD */
export function getWeekKey(ts: string): string {
  try {
    const d = new Date(ts)
    const day = d.getUTCDay()
    const diff = day === 0 ? -6 : 1 - day
    d.setUTCDate(d.getUTCDate() + diff)
    return d.toISOString().slice(0, 10)
  } catch {
    return ''
  }
}

export interface HistoryStats {
  totalOnlineSec: number
  byDay: Record<string, number>
  byWeek: Record<string, number>
}

/** Считает общий онлайн и по дням/неделям из истории (session_duration_sec при уходе оффлайн). Текущая сессия (последняя запись online) добавляется на сегодня. */
export function computeHistoryStats(history: HistoryRecord[]): HistoryStats {
  const byDay: Record<string, number> = {}
  const byWeek: Record<string, number> = {}
  let totalOnlineSec = 0
  for (const h of history) {
    if (h.online) continue
    const sec = h.session_duration_sec ?? 0
    if (sec <= 0) continue
    totalOnlineSec += sec
    const dk = getDateKey(h.ts)
    if (dk) byDay[dk] = (byDay[dk] ?? 0) + sec
    const wk = getWeekKey(h.ts)
    if (wk) byWeek[wk] = (byWeek[wk] ?? 0) + sec
  }
  const sorted = [...history].sort((a, b) => new Date(b.ts).getTime() - new Date(a.ts).getTime())
  const last = sorted[0]
  if (last?.online && last.ts) {
    const startMs = new Date(last.ts).getTime()
    const nowMs = Date.now()
    const sec = Math.max(0, Math.floor((nowMs - startMs) / 1000))
    if (sec > 0) {
      totalOnlineSec += sec
      const todayKey = getDateKey(new Date().toISOString())
      if (todayKey) byDay[todayKey] = (byDay[todayKey] ?? 0) + sec
      const wk = getWeekKey(new Date().toISOString())
      if (wk) byWeek[wk] = (byWeek[wk] ?? 0) + sec
    }
  }
  return { totalOnlineSec, byDay, byWeek }
}

export async function fetchTrackedHistory(cftoolsId: string, limit?: number): Promise<{ player: Player; history: HistoryRecord[] }> {
  const q = limit ? `?limit=${limit}` : ''
  const res = await apiFetch(`${API_BASE}/tracked/${encodeURIComponent(cftoolsId)}/history${q}`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function fetchPlayerHistory(cftoolsId: string, limit?: number): Promise<{ history: HistoryRecord[] }> {
  const q = limit ? `?limit=${limit}` : ''
  const res = await apiFetch(`${API_BASE}/players/${encodeURIComponent(cftoolsId)}/history${q}`)
  if (!res.ok) {
    const text = await res.text()
    throw new Error(res.status === 404 ? 'Player not found' : text || 'Failed to load history')
  }
  return res.json()
}
