import { useState, useEffect } from 'react'
import { API_BASE_URL, apiRequest } from '../config/api'

export interface AppStatus {
  name: string
  status: 'UP' | 'DOWN'
  uptime: string
}

const POLL_INTERVAL = 10000

export function useAppsStatus() {
  const [apps, setApps] = useState<AppStatus[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    apiRequest(`${API_BASE_URL}/api/v1/apps/status`)
      .then(r => r.json())
      .then(data => { setApps(data.apps ?? []); setLoading(false) })
      .catch(() => setLoading(false))

    const interval = setInterval(() => {
      apiRequest(`${API_BASE_URL}/api/v1/apps/status`)
        .then(r => r.json())
        .then(data => setApps(data.apps ?? []))
    }, POLL_INTERVAL)

    return () => clearInterval(interval)
  }, [])

  return { apps, loading }
}
