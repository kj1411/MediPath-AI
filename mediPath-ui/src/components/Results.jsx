import {
  Heart,
  BookOpen,
  ListChecks,
  ArrowLeft,
  ShieldCheck,
} from 'lucide-react'
import './Results.css'

export default function Results({ data, onReset }) {
  const { disease, confidence, explanation, routine } = data
  const pct = Math.round(confidence * 100)

  const getConfidenceLevel = (val) => {
    if (val >= 0.7) return 'high'
    if (val >= 0.4) return 'medium'
    return 'low'
  }

  const confidenceLevel = getConfidenceLevel(confidence)

  const renderTextBlock = (text) =>
    text
      .split('\n')
      .filter(Boolean)
      .map((line, i) => <p key={i}>{line}</p>)

  return (
    <section className="results-section">
      <div className="results-header">
        <button onClick={onReset} className="back-btn">
          <ArrowLeft size={18} />
          New Analysis
        </button>
        <h2>Your Health Insight</h2>
        <p>Based on your medications, here&apos;s what we found</p>
      </div>

      <div className="result-card disease-card">
        <div className="disease-info">
          <div className="disease-icon">
            <Heart size={32} />
          </div>
          <div>
            <span className="card-label">Predicted Condition</span>
            <h3 className="disease-name">{disease}</h3>
          </div>
        </div>
        <div className={`confidence-badge ${confidenceLevel}`}>
          <div className="confidence-bar">
            <div
              className="confidence-fill"
              style={{ width: `${pct}%` }}
            />
          </div>
          <span className="confidence-text">{pct}% confidence</span>
        </div>
      </div>

      <div className="result-card explanation-card">
        <div className="card-header">
          <BookOpen size={24} />
          <h3>What This Means</h3>
        </div>
        <div className="card-body">{renderTextBlock(explanation)}</div>
      </div>

      <div className="result-card routine-card">
        <div className="card-header">
          <ListChecks size={24} />
          <h3>Your Care Routine</h3>
        </div>
        <div className="card-body">{renderTextBlock(routine)}</div>
      </div>

      <div className="disclaimer">
        <ShieldCheck size={18} />
        <p>
          This information is for educational purposes only and does not replace
          professional medical advice. Always consult your healthcare provider
          for diagnosis and treatment.
        </p>
      </div>
    </section>
  )
}
