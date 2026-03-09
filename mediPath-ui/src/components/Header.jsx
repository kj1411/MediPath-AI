import { Activity } from 'lucide-react'
import './Header.css'

export default function Header() {
  return (
    <header className="header">
      <div className="header-inner">
        <div className="logo">
          <div className="logo-icon">
            <Activity size={24} />
          </div>
          <span className="logo-text">MediPath</span>
        </div>
        <span className="header-badge">Patient Care Navigator</span>
      </div>
    </header>
  )
}
