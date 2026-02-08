import { createContext, useContext, useState, useCallback } from 'react'
import { PlayerModal } from '../components/PlayerModal'

type OpenPlayerModal = (playerId: string) => void

const PlayerModalContext = createContext<OpenPlayerModal | null>(null)

export function PlayerModalProvider({ children }: { children: React.ReactNode }) {
  const [playerId, setPlayerId] = useState<string | null>(null)

  const openPlayerModal = useCallback((id: string) => {
    setPlayerId(id)
  }, [])

  return (
    <PlayerModalContext.Provider value={openPlayerModal}>
      {children}
      {playerId && (
        <PlayerModal
          playerId={playerId}
          onClose={() => setPlayerId(null)}
          onOpenPlayer={setPlayerId}
        />
      )}
    </PlayerModalContext.Provider>
  )
}

export function usePlayerModal() {
  const ctx = useContext(PlayerModalContext)
  return ctx
}
