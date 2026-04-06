import { useState, useEffect } from 'react'
import { API_BASE_URL, apiRequest } from '../config/api'

interface Stats {
  total_deploys: number
  success_rate: number
  active_apps: number
  last_deploy_at: string | null
}

const POLL_INTERVAL = 10000

export function useStats() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    apiRequest(`${API_BASE_URL}/api/v1/stats`)
      .then(r => r.json())
      .then(data => { setStats(data); setLoading(false) })
      .catch(() => setLoading(false))

    const interval = setInterval(() => {
      apiRequest(`${API_BASE_URL}/api/v1/stats`)
        .then(r => r.json())
        .then(data => setStats(data))
    }, POLL_INTERVAL)

    return () => clearInterval(interval)
  }, [])

  return { stats, loading }
}
