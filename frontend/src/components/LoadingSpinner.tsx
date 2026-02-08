import './LoadingSpinner.css'

type Variant = 'spinner' | 'dots' | 'pulse'

export function LoadingSpinner({
  variant = 'spinner',
  size = 'medium',
  className = '',
}: {
  variant?: Variant
  size?: 'small' | 'medium' | 'large'
  className?: string
}) {
  const sizeClass = `loading--${size}`
  if (variant === 'dots') {
    return (
      <div className={`loading-dots ${sizeClass} ${className}`} aria-hidden>
        <span />
        <span />
        <span />
      </div>
    )
  }
  if (variant === 'pulse') {
    return (
      <div className={`loading-pulse ${sizeClass} ${className}`} aria-hidden>
        <span />
      </div>
    )
  }
  return (
    <div className={`loading-spinner ${sizeClass} ${className}`} aria-hidden>
      <span />
    </div>
  )
}

export function LoadingOverlay({ text }: { text?: string }) {
  return (
    <div className="loading-overlay">
      <LoadingSpinner variant="spinner" size="large" />
      {text && <p className="loading-overlay-text">{text}</p>}
    </div>
  )
}
