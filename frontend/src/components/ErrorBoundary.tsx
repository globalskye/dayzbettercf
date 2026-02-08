import { Component, type ReactNode } from 'react'
import './ErrorBoundary.css'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error?: Error
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="error-boundary">
          <h2>Ошибка приложения</h2>
          <pre className="error-boundary-message">
            {this.state.error?.message ?? 'Неизвестная ошибка'}
          </pre>
          <button
            type="button"
            onClick={() => this.setState({ hasError: false, error: undefined })}
          >
            Попробовать снова
          </button>
        </div>
      )
    }
    return this.props.children
  }
}
