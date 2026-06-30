'use client'

import { useState } from 'react'

export default function Dashboard() {
  const [reportCode, setReportCode] = useState('')
  const [lookupResult, setLookupResult] = useState(null)
  const [lookupError, setLookupError] = useState('')
  const [lookingUp, setLookingUp] = useState(false)

  async function handleLookup(e) {
    e.preventDefault()
    if (!reportCode.trim()) return
    setLookingUp(true)
    setLookupError('')
    setLookupResult(null)

    try {
      const apiUrl = localStorage.getItem('swiftac_api_url') || 'https://swiftac-api.onrender.com'
      const res = await fetch(`${apiUrl}/api/reports/${reportCode.trim()}`)
      if (res.status === 404) {
        setLookupError('Report not found. Make sure the player typed the code correctly.')
        setLookingUp(false)
        return
      }
      if (!res.ok) throw new Error(`API returned ${res.status}`)
      const data = await res.json()
      window.location.href = `/scan?code=${reportCode.trim()}`
    } catch (e) {
      setLookupError(e.message)
    }
    setLookingUp(false)
  }

  return (
    <div>
      <h1 style={{ marginBottom: '0.5rem' }}>SwiftAntiCheat</h1>
      <p style={{ color: '#888', marginBottom: '2rem' }}>
        Enter a report code to view scan results.
      </p>

      <div className="card" style={{ textAlign: 'center', padding: '3rem', maxWidth: '500px', margin: '0 auto 2rem' }}>
        <h2 style={{ marginBottom: '1rem' }}>Look Up Report</h2>
        <p style={{ color: '#888', marginBottom: '1.5rem', fontSize: '0.9rem' }}>
          Ask the player for their report code and enter it below.
        </p>
        <form onSubmit={handleLookup} style={{ display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
          <input
            type="text"
            value={reportCode}
            onChange={e => setReportCode(e.target.value.toUpperCase())}
            placeholder="e.g. SWIFT-A1B2-C3D4"
            style={{
              padding: '0.8rem 1rem', borderRadius: '6px', border: '1px solid #3a3a5a',
              background: '#1a1a2e', color: 'white', fontSize: '1rem',
              fontFamily: 'monospace', width: '280px', letterSpacing: '1px'
            }}
          />
          <button
            type="submit"
            disabled={lookingUp || !reportCode.trim()}
            style={{
              padding: '0.8rem 1.5rem', borderRadius: '6px', border: 'none',
              background: lookingUp ? '#555' : '#7c4dff', color: 'white',
              fontWeight: 600, cursor: lookingUp ? 'default' : 'pointer', fontSize: '1rem'
            }}
          >
            {lookingUp ? '...' : 'Look Up'}
          </button>
        </form>
        {lookupError && (
          <p style={{ color: '#ff5252', fontSize: '0.85rem', marginTop: '1rem' }}>{lookupError}</p>
        )}
      </div>

      <div className="card" style={{ padding: '2rem' }}>
        <h2>How It Works</h2>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '1.5rem', marginTop: '1rem' }}>
          <div>
            <div style={{ fontSize: '1.5rem', marginBottom: '0.5rem', color: '#7c4dff' }}>1️⃣</div>
            <h3 style={{ fontSize: '1rem', marginBottom: '0.3rem' }}>Run Command</h3>
            <p style={{ color: '#888', fontSize: '0.85rem' }}>
              Use <code style={{ background: '#2a2a4a', padding: '0.1rem 0.3rem', borderRadius: '4px' }}>/swiftanticheat &lt;player&gt;</code> in-game
            </p>
          </div>
          <div>
            <div style={{ fontSize: '1.5rem', marginBottom: '0.5rem', color: '#7c4dff' }}>2️⃣</div>
            <h3 style={{ fontSize: '1rem', marginBottom: '0.3rem' }}>Player Scans</h3>
            <p style={{ color: '#888', fontSize: '0.85rem' }}>
              Player clicks the link, downloads the scanner, and runs it
            </p>
          </div>
          <div>
            <div style={{ fontSize: '1.5rem', marginBottom: '0.5rem', color: '#7c4dff' }}>3️⃣</div>
            <h3 style={{ fontSize: '1rem', marginBottom: '0.3rem' }}>Enter Code</h3>
            <p style={{ color: '#888', fontSize: '0.85rem' }}>
              Player gives you their report code — enter it above to see results
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
