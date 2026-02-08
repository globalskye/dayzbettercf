import { Outlet, Link, useNavigate } from 'react-router-dom'
import { useEffect, useState, useRef } from 'react'
import { fetchHealth } from '../api/client'
import { useAuth } from '../context/AuthContext'
import './Layout.css'

const RADIO_STATIONS = [
  ['https://ep256.hostingradio.ru:8052/europaplus256.mp3', 'Европа плюс', 'europaplus'],
  ['https://radiorecord.hostingradio.ru/rr_main96.aacp', 'Радио Рекорд', 'radiorecord'],
  ['https://nashe1.hostingradio.ru/nashe-256', 'Наше радио', 'nashe'],
  ['https://online.kissfm.ua/KissFM_HD', 'Kiss FM', 'kissfm'],
]

export function Layout() {
  const [health, setHealth] = useState<string | null>(null)
  const [navSearch, setNavSearch] = useState('')
  const [copied, setCopied] = useState<string | null>(null)
  const radioContainerRef = useRef<HTMLDivElement>(null)
  const radioLoadedRef = useRef(false)
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopied(id)
      setTimeout(() => setCopied(null), 2000)
    })
  }

  useEffect(() => {
    fetchHealth()
      .then((r) => setHealth(r?.status ?? 'unknown'))
      .catch(() => setHealth('offline'))
    const t = setInterval(() => {
      fetchHealth().then((r) => setHealth(r?.status ?? 'unknown')).catch(() => setHealth('offline'))
    }, 30000)
    return () => clearInterval(t)
  }, [])

  useEffect(() => {
    if (!radioContainerRef.current || radioLoadedRef.current) return
    radioLoadedRef.current = true
    const root = document.getElementById('root')
    if (!root) return
    const link = document.createElement('link')
    link.href = 'https://www.radiobells.com/script/style.css'
    link.rel = 'stylesheet'
    link.type = 'text/css'
    document.head.appendChild(link)
    const config = document.createElement('script')
    config.textContent = [
      'var rad_backcolor="#434242";',
      'var rad_logo="black";',
      'var rad_autoplay=false;',
      'var rad_width="responsive";',
      'var rad_width_px=330;',
      'var rad_stations=' + JSON.stringify(RADIO_STATIONS) + ';',
    ].join(' ')
    document.body.appendChild(config)
    const script = document.createElement('script')
    script.src = 'https://www.radiobells.com/script/v2_1.js'
    script.charset = 'UTF-8'
    script.async = true
    document.body.appendChild(script)
  }, [])

  const handleNavSearch = () => {
    const q = navSearch.trim()
    if (!q) return
    navigate(`/?q=${encodeURIComponent(q)}`)
    setNavSearch('')
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="layout">
      <header className="layout-header">
        <Link to="/" className="layout-logo">
          DayZ Smart CF
        </Link>
        <nav className="layout-nav">
          <div className="layout-search-wrap">
            <input
              type="text"
              className="layout-search-input"
              placeholder="Поиск по CF (ник, identifier…)"
              value={navSearch}
              onChange={(e) => setNavSearch(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleNavSearch()}
              title="Поиск по CFTools — откроет главную с результатами"
            />
            <button
              type="button"
              className="layout-search-btn"
              onClick={handleNavSearch}
              disabled={!navSearch.trim()}
              title="Искать на главной"
            >
              Найти
            </button>
          </div>
          <Link to="/" className="layout-link">Главная</Link>
          <Link to="/base" className="layout-link">База</Link>
          <Link to="/groups" className="layout-link">Группы</Link>
          <Link to="/tracked" className="layout-link">Отслеживание</Link>
          {user?.role === 'admin' && (
            <Link to="/settings" className="layout-link">Настройки</Link>
          )}
          {user && (
            <span className="layout-user">
              {user.username} ({user.role})
              <button type="button" className="layout-logout" onClick={handleLogout}>Выход</button>
            </span>
          )}
          <span className={`layout-status ${health === 'ok' ? 'online' : 'offline'}`}>
          {health ?? '...'}
          </span>
        </nav>
      </header>
      <main className="layout-main">
        <Outlet />
      </main>
      <footer className="layout-footer">
        <div className="layout-footer-inner">
          <section className="layout-radio">
            <span className="layout-radio-label">По рофлу</span>
            <div ref={radioContainerRef} id="radiobells_container">
              <a href="https://www.radiobells.com/" target="_blank" rel="noopener noreferrer" className="layout-radio-link">
                Онлайн радио
              </a>
            </div>
          </section>
          <section className="layout-plans">
            <h3 className="layout-plans-title">Планы на будущее</h3>
            <ul className="layout-plans-list">
              <li>Жёсткий поиск людей по банам — фильтры по причине, серверу, дате</li>
              <li>Поиск по Steam-аккаунтам и парсинг друзей — вычисление связей, альты, граф</li>
              
            </ul>
          </section>
          <section className="layout-support">
            <h3 className="layout-support-title">Хочешь поддержать проект?</h3>
            <div className="layout-support-grid">
              <div
                className="layout-support-item layout-support-clickable"
                onClick={() => copyToClipboard('5S1CD1CJTsK9hrbkHQUKXZhvZKx47g5uYGvQ2RqMXEuu', 'sol')}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => e.key === 'Enter' && copyToClipboard('5S1CD1CJTsK9hrbkHQUKXZhvZKx47g5uYGvQ2RqMXEuu', 'sol')}
                title="Копировать"
              >
                <span className="layout-support-label">SOL (Solana)</span>
                <code className="layout-support-address">
                  5S1CD1CJTsK9hrbkHQUKXZhvZKx47g5uYGvQ2RqMXEuu
                </code>
                {copied === 'sol' && <span className="layout-support-copied">Скопировано</span>}
              </div>
              <div
                className="layout-support-item layout-support-clickable"
                onClick={() => copyToClipboard('TWaLfNRe8rX3iRVq3AaJTXSV9pCtfmzapF', 'usdt')}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => e.key === 'Enter' && copyToClipboard('TWaLfNRe8rX3iRVq3AaJTXSV9pCtfmzapF', 'usdt')}
                title="Копировать"
              >
                <span className="layout-support-label">USDT (TRC20)</span>
                <code className="layout-support-address">
                  TWaLfNRe8rX3iRVq3AaJTXSV9pCtfmzapF
                </code>
                {copied === 'usdt' && <span className="layout-support-copied">Скопировано</span>}
              </div>
              <div
                className="layout-support-item layout-support-card layout-support-clickable"
                onClick={() => copyToClipboard('9112380128681774', 'card')}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => e.key === 'Enter' && copyToClipboard('9112380128681774', 'card')}
                title="Копировать номер карты"
              >
                <span className="layout-support-label">Карта (СБП / перевод)</span>
                <div className="layout-support-card-number">9112 3801 2868 1774</div>
                <div className="layout-support-card-name">ALIAKSEI SAMALIOTAU</div>
                {copied === 'card' && <span className="layout-support-copied">Скопировано</span>}
              </div>
            </div>
          </section>
        </div>
      </footer>
    </div>
  )
}
