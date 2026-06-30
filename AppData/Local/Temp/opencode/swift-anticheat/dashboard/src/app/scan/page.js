'use client'

import { Suspense, useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'

const API_BASE_URL = 'https://swiftac-api.onrender.com'
const SCANNER_DOWNLOAD_URL = 'https://github.com/pinakapoginghuman/SwiftAntiCheat/releases/download/v1.0.0/swiftac-scanner.exe'

function ScanContent() {
  const searchParams = useSearchParams()
  const scanId = searchParams.get('id')
  const [scan, setScan] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!scanId) {
      setError('No scan ID provided')
      setLoading(false)
      return
    }
    fetchScan()
    const interval = setInterval(fetchScan, 5000)
    return () => clearInterval(interval)
  }, [scanId])

  async function fetchScan() {
    setLoading(true)
    setError('')
    const apiUrl = localStorage.getItem('swiftac_api_url') || API_BASE_URL

    try {
      const res = await fetch(`${apiUrl}/api/scans/${scanId}`)
      if (res.status === 404) throw new Error('Scan not found')
      if (!res.ok) throw new Error(`API returned ${res.status}`)
      const data = await res.json()
      setScan(data)
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <div className="loading">Loading scan data...</div>
  if (error) return <div className="error">{error}</div>
  if (!scan) return <div className="error">Scan not found</div>

  const results = scan.results
  const flags = results?.flags || []
  const highFlags = flags.filter(f => f.severity === 'high')
  const mediumFlags = flags.filter(f => f.severity === 'medium')

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <div>
          <h1>Scan: {scan.playerName}</h1>
          <p style={{ color: '#888', fontSize: '0.85rem' }}>
            ID: {scan.id} · Staff: {scan.staffName}
            · {new Date(scan.createdAt * 1000).toLocaleString()}
          </p>
        </div>
        <span className={`status-badge status-${scan.status}`} style={{ fontSize: '0.9rem', padding: '0.4rem 1rem' }}>
          {scan.status.toUpperCase()}
        </span>
      </div>

      {scan.status === 'pending' && (
        <div className="card" style={{ textAlign: 'center', padding: '3rem', borderColor: '#3a3a5a' }}>
          <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>⏳</div>
          <h2 style={{ color: '#ffab00' }}>Scan In Progress</h2>
          <p style={{ color: '#888', marginBottom: '1.5rem' }}>Download and run the scanner to complete your scan.</p>
          <div style={{ background: '#2a2a4a', padding: '1rem', borderRadius: '8px', marginBottom: '1.5rem', textAlign: 'left', fontSize: '0.85rem' }}>
            <p style={{ color: '#7c4dff', fontWeight: 600, marginBottom: '0.5rem' }}>Instructions:</p>
            <ol style={{ color: '#b0b0c0', paddingLeft: '1.2rem', lineHeight: '1.8' }}>
              <li>Click the <strong>Download Scanner</strong> button below</li>
              <li>Open the downloaded file (swiftac-scanner.exe)</li>
              <li>It will scan your system automatically with Scan ID: <strong style={{ color: '#7c4dff' }}>{scanId}</strong></li>
              <li>Wait for "Results uploaded successfully!" message</li>
              <li>Staff will review the results</li>
            </ol>
          </div>
          <a
            href={SCANNER_DOWNLOAD_URL}
            target="_blank"
            rel="noopener noreferrer"
            style={{ display: 'inline-block', padding: '1rem 3rem', background: '#7c4dff', color: 'white', borderRadius: '8px', textDecoration: 'none', fontWeight: 700, fontSize: '1.1rem' }}
          >
            ⬇ Download Scanner
          </a>
          <p style={{ color: '#888', fontSize: '0.8rem', marginTop: '1rem' }}>
            Scan ID: {scanId}
          </p>
        </div>
      )}

      {results && (
        <>
          <div className="stats-grid">
            <div className="stat-box">
              <div className="value" style={{ color: highFlags.length > 0 ? '#ff5252' : '#00e676' }}>
                {highFlags.length}
              </div>
              <div className="label">High Flags</div>
            </div>
            <div className="stat-box">
              <div className="value" style={{ color: mediumFlags.length > 0 ? '#ffab00' : '#00e676' }}>
                {mediumFlags.length}
              </div>
              <div className="label">Medium Flags</div>
            </div>
            <div className="stat-box">
              <div className="value">{results.suspicious_files?.length || 0}</div>
              <div className="label">Suspicious Files</div>
            </div>
            <div className="stat-box">
              <div className="value">{results.minecraft_mods?.length || 0}</div>
              <div className="label">Minecraft Mods</div>
            </div>
          </div>

          {flags.length > 0 && (
            <div className="card">
              <h2>Flags & Detections</h2>
              {flags.map((flag, i) => (
                <div key={i} className="flag-item">
                  <div>
                    <div className="flag-name">{flag.name || flag.type}</div>
                    <div className="flag-detail">{flag.detail}</div>
                  </div>
                  <span className={`severity-${flag.severity} flag-severity`}>
                    {flag.severity.toUpperCase()}
                  </span>
                </div>
              ))}
            </div>
          )}

          {flags.length === 0 && (
            <div className="card" style={{ textAlign: 'center', padding: '2rem', borderColor: '#1b3a2a' }}>
              <div style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>✅</div>
              <h2 style={{ color: '#00e676' }}>No Detections</h2>
              <p style={{ color: '#888' }}>No suspicious files, processes, or modifications found.</p>
            </div>
          )}

          {results.suspicious_files?.length > 0 && (
            <div className="card">
              <h2>Suspicious Files</h2>
              <table>
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Path</th>
                    <th>Size</th>
                    <th>Reason</th>
                  </tr>
                </thead>
                <tbody>
                  {results.suspicious_files.map((file, i) => (
                    <tr key={i}>
                      <td style={{ fontWeight: 600 }}>{file.name}</td>
                      <td style={{ fontSize: '0.8rem', color: '#888', maxWidth: '300px', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                        {file.path}
                      </td>
                      <td>{(file.size / 1024).toFixed(1)} KB</td>
                      <td>
                        <span className={`severity-${file.severity} flag-severity`}>
                          {file.reason || file.severity}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <div className="card">
              <h2>System Info</h2>
              {results.system_info ? (
                <table>
                  <tbody>
                    <tr><td style={{ color: '#888' }}>Hostname</td><td>{results.system_info.hostname}</td></tr>
                    <tr><td style={{ color: '#888' }}>OS</td><td>{results.system_info.os_version}</td></tr>
                    <tr><td style={{ color: '#888' }}>Username</td><td>{results.system_info.username}</td></tr>
                    <tr><td style={{ color: '#888' }}>CPU</td><td style={{ fontSize: '0.8rem' }}>{results.system_info.cpu_info}</td></tr>
                    <tr><td style={{ color: '#888' }}>HWID</td><td style={{ fontSize: '0.7rem', fontFamily: 'monospace' }}>{results.hwid_hash?.slice(0, 16)}...</td></tr>
                  </tbody>
                </table>
              ) : (
                <p style={{ color: '#888' }}>No system info collected</p>
              )}
            </div>

            <div className="card">
              <h2>Minecraft Mods</h2>
              {results.minecraft_mods?.length > 0 ? (
                <table>
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {results.minecraft_mods.map((mod, i) => (
                      <tr key={i}>
                        <td style={{ fontSize: '0.85rem' }}>{mod.name}</td>
                        <td>
                          {mod.suspicious
                            ? <span className="severity-high flag-severity">SUSPICIOUS</span>
                            : <span style={{ color: '#888' }}>Clean</span>
                          }
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              ) : (
                <p style={{ color: '#888' }}>No Minecraft mods detected</p>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  )
}

export default function ScanPage() {
  return (
    <Suspense fallback={<div className="loading">Loading scan data...</div>}>
      <ScanContent />
    </Suspense>
  )
}