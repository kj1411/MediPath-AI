import { useState } from 'react'
import { Plus, X, Search, Pill } from 'lucide-react'
import './DrugInput.css'

export default function DrugInput({ onAnalyze, loading, error }) {
  const [drugs, setDrugs] = useState([])
  const [input, setInput] = useState('')

  const addDrug = () => {
    const drug = input.trim().toLowerCase()
    if (drug && !drugs.includes(drug)) {
      setDrugs([...drugs, drug])
      setInput('')
    }
  }

  const removeDrug = (drug) => {
    setDrugs(drugs.filter((d) => d !== drug))
  }

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      addDrug()
    }
  }

  const handleSubmit = () => {
    if (drugs.length > 0) {
      onAnalyze(drugs)
    }
  }

  return (
    <section className="drug-input-section">
      <div className="hero">
        <div className="hero-icon">
          <Pill size={48} strokeWidth={1.5} />
        </div>
        <h1>Understand Your Medications</h1>
        <p>
          Enter the medications from your prescription. We&apos;ll help you
          understand what condition they treat and provide a simple care guide.
        </p>
      </div>

      <div className="input-card">
        <label className="input-label">Your Medications</label>
        <div className="input-row">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Type a medication name..."
            className="drug-input"
            disabled={loading}
          />
          <button
            onClick={addDrug}
            className="add-btn"
            disabled={!input.trim() || loading}
          >
            <Plus size={20} />
            Add
          </button>
        </div>

        {drugs.length > 0 && (
          <div className="drug-tags">
            {drugs.map((drug) => (
              <span key={drug} className="drug-tag">
                <Pill size={14} />
                {drug}
                <button
                  onClick={() => removeDrug(drug)}
                  className="tag-remove"
                  disabled={loading}
                  aria-label={`Remove ${drug}`}
                >
                  <X size={14} />
                </button>
              </span>
            ))}
          </div>
        )}

        {error && (
          <div className="error-message">
            <span>⚠️</span> {error}
          </div>
        )}

        <button
          onClick={handleSubmit}
          className="analyze-btn"
          disabled={drugs.length === 0 || loading}
        >
          {loading ? (
            <>
              <span className="spinner" />
              Analyzing...
            </>
          ) : (
            <>
              <Search size={20} />
              Analyze My Medications
            </>
          )}
        </button>
      </div>
    </section>
  )
}
