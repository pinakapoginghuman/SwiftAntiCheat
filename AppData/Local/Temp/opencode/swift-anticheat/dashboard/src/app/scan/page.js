'use client'

import { useState, useEffect, Suspense } from 'react'
import { useSearchParams } from 'next/navigation'

const SCANNER_DOWNLOAD_URL = 'https://github.com/pinakapoginghuman/SwiftAntiCheat/releases/download/v1.2.0/swiftac-scanner.exe'

function ScanContent() {
  const searchParams = useSearchParams()
  const codeParam = searchParams.get('code')

  const [reportCode, setReportCode] = useState(codeParam || '')
  const [lookupResult, setLookupResult] = useState(null)
  const [lookupError, setLookupError] = useState('')
  const [lookingUp, setLookingUp] = useState(false)
  const [autoLookedUp, setAutoLookedUp] = useState(false)

  useEffect(() => {
    if (codeParam && !autoLookedUp) {
      setAutoLookedUp(true)
      doLookup(codeParam)
    }
  }, [codeParam])

  async function doLookup(code) {
    setLookingUp(true)
    setLookupError('')
    setLookupResult(null)

    try {
      const apiUrl = localStorage.getItem('swiftac_api_url') || 'https://swiftac-api.onrender.com'
      const res = await fetch(`${apiUrl}/api/reports/${code}`)
      if (res.status === 404) {
        setLookupError('Report not found. Make sure the code is correct.')
        setLookingUp(false)
        return
      }
      if (!res.ok) throw new Error(`API returned ${res.status}`)
      const data = await res.json()
      setLookupResult(data)
    } catch (e) {
      setLookupError(e.message)
    }
    setLookingUp(false)
  }

  async function handleLookup(e) {
    e.preventDefault()
    if (!reportCode.trim()) return
    doLookup(reportCode.trim())
  }

  if (lookupResult) {
    return <ResultsView results={lookupResult.results} reportCode={lookupResult.reportCode} />
  }

  return (
    <div>
      <div className="card" style={{ textAlign: 'center', padding: '3rem', borderColor: '#3a3a5a', maxWidth: '700px', margin: '0 auto' }}>
        <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>🛡️</div>
        <h1 style={{ marginBottom: '0.5rem' }}>SwiftAntiCheat Scanner</h1>
        <p style={{ color: '#888', marginBottom: '2rem' }}>
          A staff member has requested you to complete a scan. Follow the steps below.
        </p>

        <div style={{ background: '#1a1a2e', padding: '1.5rem', borderRadius: '12px', marginBottom: '2rem', textAlign: 'left', fontSize: '0.95rem' }}>
          <p style={{ color: '#7c4dff', fontWeight: 700, marginBottom: '1rem', fontSize: '1.1rem' }}>📋 Step-by-Step Instructions</p>
          <ol style={{ color: '#c0c0d0', paddingLeft: '1.2rem', lineHeight: '2.4', margin: 0 }}>
            <li>
              <strong style={{ color: '#e0e0e0' }}>Click the Download button below</strong>
              <span style={{ display: 'block', fontSize: '0.8rem', color: '#888' }}>
                It will download <code style={{ background: '#2a2a4a', padding: '0.1rem 0.4rem', borderRadius: '4px' }}>swiftac-scanner.exe</code>
              </span>
            </li>
            <li>
              <strong style={{ color: '#e0e0e0' }}>Double-click the downloaded file</strong>
              <span style={{ display: 'block', fontSize: '0.8rem', color: '#888' }}>
                Windows may show a smart screen warning — click "Run anyway"
              </span>
            </li>
            <li>
              <strong style={{ color: '#e0e0e0' }}>Wait a few seconds</strong>
              <span style={{ display: 'block', fontSize: '0.8rem', color: '#888' }}>
                A spinner animates while it scans — <strong>do not close the window</strong>
              </span>
            </li>
            <li>
              <strong style={{ color: '#e0e0e0' }}>Copy your Report Code</strong>
              <span style={{ display: 'block', fontSize: '0.8rem', color: '#888' }}>
                The scanner will show <strong style={{ color: '#7c4dff' }}>SWIFT-XXXX-XXXX</strong>
              </span>
            </li>
            <li>
              <strong style={{ color: '#e0e0e0' }}>Send it to staff on Discord</strong>
              <span style={{ display: 'block', fontSize: '0.8rem', color: '#888' }}>
                They&apos;ll enter it on the dashboard to see results
              </span>
            </li>
          </ol>
        </div>

        <a
          href={SCANNER_DOWNLOAD_URL}
          target="_blank"
          rel="noopener noreferrer"
          style={{ display: 'inline-block', padding: '1rem 3rem', background: '#7c4dff', color: 'white', borderRadius: '8px', textDecoration: 'none', fontWeight: 700, fontSize: '1.1rem', marginBottom: '1rem' }}
        >
          ⬇ Download Scanner
        </a>

        <p style={{ color: '#555', fontSize: '0.8rem', marginBottom: '2rem' }}>
          Windows only · v1.1.0 · ~8 MB
        </p>

        <div style={{ borderTop: '1px solid #2a2a4a', paddingTop: '1.5rem' }}>
          <p style={{ color: '#888', fontSize: '0.85rem', marginBottom: '0.5rem' }}>
            Already have a report code? Look up your results:
          </p>
          <form onSubmit={handleLookup} style={{ display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
            <input
              type="text"
              value={reportCode}
              onChange={e => setReportCode(e.target.value.toUpperCase())}
              placeholder="Enter report code..."
              style={{
                padding: '0.6rem 1rem', borderRadius: '6px', border: '1px solid #3a3a5a',
                background: '#1a1a2e', color: 'white', fontSize: '0.9rem',
                fontFamily: 'monospace', width: '220px'
              }}
            />
            <button
              type="submit"
              disabled={lookingUp || !reportCode.trim()}
              style={{
                padding: '0.6rem 1.2rem', borderRadius: '6px', border: 'none',
                background: lookingUp ? '#555' : '#7c4dff', color: 'white',
                fontWeight: 600, cursor: lookingUp ? 'default' : 'pointer'
              }}
            >
              {lookingUp ? 'Searching...' : 'Look Up'}
            </button>
          </form>
          {lookupError && (
            <p style={{ color: '#ff5252', fontSize: '0.85rem', marginTop: '0.5rem' }}>{lookupError}</p>
          )}
        </div>
      </div>

      <div style={{ textAlign: 'center', marginTop: '1.5rem' }}>
        <p style={{ color: '#555', fontSize: '0.8rem' }}>
          SwiftAntiCheat  ·  Your privacy is protected  ·  No personal data is stored
        </p>
      </div>
    </div>
  )
}

export default function ScanPage() {
  return (
    <Suspense fallback={<div className="loading">Loading...</div>}>
      <ScanContent />
    </Suspense>
  )
}

function ResultsView({ results, reportCode }) {
  const flags = results?.flags || []
  const highFlags = flags.filter(f => f.severity === 'high')
  const mediumFlags = flags.filter(f => f.severity === 'medium')

  if (!results) {
    return <div className="error">No results data found for this report code.</div>
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <div>
          <h1>Scan Results</h1>
          <p style={{ color: '#888', fontSize: '0.85rem' }}>
            Report Code: <strong style={{ color: '#7c4dff' }}>{reportCode}</strong>
          </p>
        </div>
      </div>

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
        <div className="stat-box">
          <div className="value">{results.startup_programs?.length || 0}</div>
          <div className="label">Startup Entries</div>
        </div>
        <div className="stat-box">
          <div className="value">{(results.suspicious_services?.length || 0) + (results.installed_programs?.length || 0)}</div>
          <div className="label">Services + Programs</div>
        </div>
      </div>

      {flags.length > 0 && (
        <>
          {highFlags.length > 0 && (
            <div className="card">
              <h2 style={{ color: '#ff5252' }}>🔥 High Severity Flags</h2>
              {flags.filter(f => f.severity === 'high').map((flag, i) => (
                <div key={i} className="flag-item">
                  <div>
                    <div className="flag-name">{flag.name || flag.type}</div>
                    <div className="flag-detail">{flag.detail}</div>
                  </div>
                  <span className="severity-high flag-severity">HIGH</span>
                </div>
              ))}
            </div>
          )}

          {mediumFlags.length > 0 && (
            <div className="card">
              <h2 style={{ color: '#ffab00' }}>⚠️ Medium Severity Flags</h2>
              {flags.filter(f => f.severity === 'medium').map((flag, i) => (
                <div key={i} className="flag-item">
                  <div>
                    <div className="flag-name">{flag.name || flag.type}</div>
                    <div className="flag-detail">{flag.detail}</div>
                  </div>
                  <span className="severity-medium flag-severity">MEDIUM</span>
                </div>
              ))}
            </div>
          )}
        </>
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

      {results.startup_programs?.length > 0 && (
        <div className="card">
          <h2>Startup Programs</h2>
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Command</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {results.startup_programs.map((entry, i) => (
                <tr key={i}>
                  <td style={{ fontWeight: 600 }}>{entry.name}</td>
                  <td style={{ fontSize: '0.8rem', color: '#888', maxWidth: '300px', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                    {entry.command}
                  </td>
                  <td>
                    {entry.suspicious
                      ? <span className="severity-high flag-severity">SUSPICIOUS</span>
                      : <span style={{ color: '#888' }}>Clean</span>
                    }
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {results.suspicious_services?.length > 0 && (
        <div className="card">
          <h2>Suspicious Services</h2>
          <ul style={{ color: '#b0b0c0', lineHeight: '2' }}>
            {results.suspicious_services.map((svc, i) => (
              <li key={i}>{svc}</li>
            ))}
          </ul>
        </div>
      )}

      {results.installed_programs?.length > 0 && (
        <div className="card">
          <h2>Suspicious Installed Programs</h2>
          <ul style={{ color: '#b0b0c0', lineHeight: '2' }}>
            {results.installed_programs.map((prog, i) => (
              <li key={i}>{prog}</li>
            ))}
          </ul>
        </div>
      )}

      {results.windows_artifacts?.event_log_cleared && (
        <div className="card" style={{ borderColor: '#ff5252' }}>
          <h2 style={{ color: '#ff5252' }}>🚨 Event Log Cleared</h2>
          <p style={{ color: '#b0b0c0' }}>The Windows event log was cleared — a common anti-forensic technique used to hide cheat traces.</p>
        </div>
      )}

      {results.windows_artifacts?.dns_cache?.length > 0 && (
        <div className="card">
          <h2>DNS Cache Entries</h2>
          <ul style={{ color: '#b0b0c0', lineHeight: '2' }}>
            {results.windows_artifacts.dns_cache.map((dns, i) => (
              <li key={i}>{dns}</li>
            ))}
          </ul>
        </div>
      )}

      {results.windows_artifacts?.prefetch_files?.length > 0 && (
        <div className="card">
          <h2>Prefetch Entries</h2>
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Last Run</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {results.windows_artifacts.prefetch_files.map((pf, i) => (
                <tr key={i}>
                  <td style={{ fontWeight: 600 }}>{pf.name}</td>
                  <td style={{ color: '#888', fontSize: '0.85rem' }}>{pf.last_run}</td>
                  <td>
                    {pf.suspicious
                      ? <span className="severity-medium flag-severity">SUSPICIOUS</span>
                      : <span style={{ color: '#888' }}>Clean</span>
                    }
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
                <tr><td style={{ color: '#888' }}>RAM</td><td>{results.system_info.total_ram ? (parseInt(results.system_info.total_ram) / 1073741824).toFixed(1) + ' GB' : 'Unknown'}</td></tr>
                <tr><td style={{ color: '#888' }}>MAC</td><td style={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{results.system_info.mac_address}</td></tr>
                <tr><td style={{ color: '#888' }}>Motherboard</td><td style={{ fontSize: '0.8rem' }}>{results.system_info.motherboard_serial}</td></tr>
                <tr><td style={{ color: '#888' }}>HWID</td><td style={{ fontSize: '0.7rem', fontFamily: 'monospace' }}>{results.hwid_hash?.slice(0, 16)}...</td></tr>
                <tr><td style={{ color: '#888' }}>Boot Time</td><td style={{ fontSize: '0.8rem' }}>{results.system_info.boot_time}</td></tr>
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

          {results.windows_artifacts?.recent_documents?.length > 0 && (
            <div style={{ marginTop: '1rem' }}>
              <h3>Recent Documents</h3>
              <ul style={{ fontSize: '0.8rem', color: '#b0b0c0', lineHeight: '1.8' }}>
                {results.windows_artifacts.recent_documents.map((doc, i) => (
                  <li key={i}>{doc}</li>
                ))}
              </ul>
            </div>
          )}

          {results.windows_artifacts?.suspicious_registry_keys?.length > 0 && (
            <div style={{ marginTop: '1rem' }}>
              <h3>Suspicious Registry Entries</h3>
              <ul style={{ fontSize: '0.75rem', color: '#b0b0c0', lineHeight: '1.8', wordBreak: 'break-all' }}>
                {results.windows_artifacts.suspicious_registry_keys.map((rk, i) => (
                  <li key={i}>{rk}</li>
                ))}
              </ul>
            </div>
          )}

          {results.windows_artifacts?.suspicious_run_keys?.length > 0 && (
            <div style={{ marginTop: '1rem' }}>
              <h3>Suspicious Auto-Run Entries</h3>
              <ul style={{ fontSize: '0.8rem', color: '#b0b0c0', lineHeight: '1.8' }}>
                {results.windows_artifacts.suspicious_run_keys.map((rk, i) => (
                  <li key={i}>{rk}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
