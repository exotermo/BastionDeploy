import { useState, useEffect } from 'react'
import { API_BASE_URL } from '../config/api'

interface Stats {
  total_deploys: number
  success_rate: number
  active_apps: number
  last_deploy_at: string | null
}

export function useStats() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch(`${API_BASE_URL}/api/v1/stats`)
      .then(r => r.json())
      .then(data => { setStats(data); setLoading(false) })
      .catch(() => setLoading(false))

    // Atualiza a cada 10 segundos
    const interval = setInterval(() => {
      fetch(`${API_BASE_URL}/api/v1/stats`)
        .then(r => r.json())
        .then(data => setStats(data))
    }, 10000)

    return () => clearInterval(interval)
  }, [])

  return { stats, loading }
}