import { useState, useEffect } from 'react'
import { API_BASE_URL } from '../config/api'

export interface AppStatus {
  name: string
  status: 'UP' | 'DOWN'
  uptime: string
}

export function useAppsStatus() {
  const [apps, setApps] = useState<AppStatus[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch(`${API_BASE_URL}/api/v1/apps/status`)
      .then(r => r.json())
      .then(data => { setApps(data.apps ?? []); setLoading(false) })
      .catch(() => setLoading(false))

    const interval = setInterval(() => {
      fetch(`${API_BASE_URL}/api/v1/apps/status`)
        .then(r => r.json())
        .then(data => setApps(data.apps ?? []))
    }, 10000)

    return () => clearInterval(interval)
  }, [])

  return { apps, loading }
}