'use client'

import { useState, useEffect } from 'react'

export default function ScansPage() {
  const [scans, setScans] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [search, setSearch] = useState('')

  useEffect(() => {
    fetchScans()
  }, [])

  async function fetchScans(player) {
    setLoading(true)
    setError('')
    const apiUrl = localStorage.getItem('swiftac_api_url') || 'http://localhost:3000'

    try {
      let url = `${apiUrl}/api/scans?limit=100`
      if (player) url += `&player=${encodeURIComponent(player)}`
      const res = await fetch(url)
      if (!res.ok) throw new Error(`API returned ${res.status}`)
      const data = await res.json()
      setScans(data)
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  function handleSearch(e) {
    e.preventDefault()
    fetchScans(search)
  }

  return (
    <div>
      <h1 style={{ marginBottom: '1rem' }}>All Scans</h1>

      <form className="search-bar" onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Search by player name..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <button type="submit">Search</button>
      </form>

      {loading && <div className="loading">Loading scans...</div>}
      {error && <div className="error">{error}</div>}

      {!loading && !error && (
        <div className="card">
          <table>
            <thead>
              <tr>
                <th>Player</th>
                <th>UUID</th>
                <th>Staff</th>
                <th>Status</th>
                <th>Flags</th>
                <th>Date</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {scans.length === 0 && (
                <tr>
                  <td colSpan={7} style={{ textAlign: 'center', color: '#888', padding: '2rem' }}>
                    No scans found
                  </td>
                </tr>
              )}
              {scans.map(scan => (
                <tr key={scan.id}>
                  <td style={{ fontWeight: 600 }}>{scan.playerName}</td>
                  <td style={{ fontSize: '0.8rem', color: '#888' }}>
                    {scan.playerUUID?.slice(0, 8)}...
                  </td>
                  <td>{scan.staffName}</td>
                  <td>
                    <span className={`status-badge status-${scan.status}`}>
                      {scan.status}
                    </span>
                  </td>
                  <td>
                    {scan.results?.flags?.length > 0
                      ? <span style={{ color: '#ff5252' }}>{scan.results.flags.length}</span>
                      : <span style={{ color: '#00e676' }}>0</span>
                    }
                  </td>
                  <td style={{ fontSize: '0.85rem' }}>
                    {new Date(scan.createdAt * 1000).toLocaleDateString()}
                  </td>
                  <td>
                    <a
                      href={`/scan?id=${scan.id}`}
                      style={{ color: '#7c4dff', textDecoration: 'none' }}
                    >
                      View →
                    </a>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
