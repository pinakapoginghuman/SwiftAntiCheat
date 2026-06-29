'use client'

import { useState, useEffect } from 'react'

export default function Dashboard() {
  const [stats, setStats] = useState(null)
  const [recentScans, setRecentScans] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchData()
  }, [])

  async function fetchData() {
    setLoading(true)
    setError('')
    const apiUrl = localStorage.getItem('swiftac_api_url') || 'http://localhost:3000'

    try {
      const [scansRes] = await Promise.all([
        fetch(`${apiUrl}/api/scans?limit=10`),
      ])

      if (!scansRes.ok) throw new Error(`API returned ${scansRes.status}`)

      const scans = await scansRes.json()
      setRecentScans(scans)

      setStats({
        total: scans.length,
        completed: scans.filter(s => s.status === 'completed').length,
        pending: scans.filter(s => s.status === 'pending').length,
        players: [...new Set(scans.map(s => s.playerName))].length,
      })
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <div className="loading">Loading dashboard...</div>
  if (error) return <div className="error">Failed to connect to API: {error}</div>

  return (
    <div>
      <h1 style={{ marginBottom: '1rem' }}>Dashboard</h1>

      {stats && (
        <div className="stats-grid">
          <div className="stat-box">
            <div className="value">{stats.total}</div>
            <div className="label">Total Scans</div>
          </div>
          <div className="stat-box">
            <div className="value">{stats.completed}</div>
            <div className="label">Completed</div>
          </div>
          <div className="stat-box">
            <div className="value">{stats.pending}</div>
            <div className="label">Pending</div>
          </div>
          <div className="stat-box">
            <div className="value">{stats.players}</div>
            <div className="label">Players Scanned</div>
          </div>
        </div>
      )}

      <div className="card">
        <h2>Recent Scans</h2>
        {recentScans.length === 0 ? (
          <p style={{ color: '#888', textAlign: 'center', padding: '1rem' }}>
            No scans yet. Use /swiftanticheat &lt;player&gt; in-game to generate one.
          </p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Player</th>
                <th>Staff</th>
                <th>Status</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {recentScans.map(scan => (
                <tr key={scan.id}>
                  <td style={{ fontWeight: 600 }}>{scan.playerName}</td>
                  <td>{scan.staffName}</td>
                  <td>
                    <span className={`status-badge status-${scan.status}`}>
                      {scan.status}
                    </span>
                  </td>
                  <td>{new Date(scan.createdAt * 1000).toLocaleString()}</td>
                  <td>
                    <a
                      href={`/scan?id=${scan.id}`}
                      style={{ color: '#7c4dff', textDecoration: 'none' }}
                    >
                      View
                    </a>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
