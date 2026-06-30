'use client'

import { useState, useEffect } from 'react'
import './globals.css'

export default function RootLayout({ children }) {
  const [apiUrl, setApiUrl] = useState('')

  useEffect(() => {
    setApiUrl(localStorage.getItem('swiftac_api_url') || 'https://swiftac-api.onrender.com')
  }, [])

  return (
    <html lang="en">
      <body>
        <nav className="navbar">
          <div className="nav-brand">
            <h1>SwiftAntiCheat</h1>
          </div>
          <div className="nav-links">
            <a href="/">Dashboard</a>
            <a href="/scans">Scans</a>
            <a href="/scan">Live Scan</a>
          </div>
          <div className="nav-api">
            <input
              type="text"
              placeholder="API URL"
              value={apiUrl}
              onChange={(e) => {
                setApiUrl(e.target.value)
                localStorage.setItem('swiftac_api_url', e.target.value)
              }}
              className="api-input"
            />
          </div>
        </nav>
        <main className="container">
          {children}
        </main>
      </body>
    </html>
  )
}
