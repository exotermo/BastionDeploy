import { useState, useEffect } from 'react'
import { API_BASE_URL, apiRequest } from '../config/api'

export interface Deploy {
  id: string
  app_name: string
  branch: string
  commit_sha: string
  status: 'success' | 'failed' | 'pending' | 'running'
  triggered_by: string
  created_at: string
  updated_at: string
}

const POLL_INTERVAL = 5000

export function useDeploys() {
  const [deploys, setDeploys] = useState<Deploy[]>([])
  const [loading, setLoading] = useState(true)

  const fetchDeploys = () =>
    apiRequest(`${API_BASE_URL}/api/v1/deploys`)
      .then(r => r.json())
      .then(data => setDeploys(data.deploys ?? []))
      .catch(() => {/* ignora erros silenciosamente em polling */})

  useEffect(() => {
    fetchDeploys().then(() => setLoading(false))
    const interval = setInterval(fetchDeploys, POLL_INTERVAL)
    return () => clearInterval(interval)
  }, [])

  return { deploys, loading }
}
