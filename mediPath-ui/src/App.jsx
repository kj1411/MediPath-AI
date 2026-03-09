import { useState } from 'react'
import Header from './components/Header'
import DrugInput from './components/DrugInput'
import Results from './components/Results'
import { predictCondition } from './api/predict'
import './App.css'

function App() {
  const [results, setResults] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const handleAnalyze = async (drugs) => {
    setLoading(true)
    setError(null)
    setResults(null)
    try {
      const data = await predictCondition(drugs)
      setResults(data)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleReset = () => {
    setResults(null)
    setError(null)
  }

  return (
    <div className="app">
      <Header />
      <main className="main">
        {!results ? (
          <DrugInput onAnalyze={handleAnalyze} loading={loading} error={error} />
        ) : (
          <Results data={results} onReset={handleReset} />
        )}
      </main>
      <footer className="footer">
        <p>⚕️ MediPath is for educational purposes only. Always consult your healthcare provider for medical advice.</p>
      </footer>
    </div>
  )
}

export default App
