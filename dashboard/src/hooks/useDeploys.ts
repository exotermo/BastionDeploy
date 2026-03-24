import { useState, useEffect } from 'react'
import { API_BASE_URL } from '../config/api'

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

export function useDeploys() {
  const [deploys, setDeploys] = useState<Deploy[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch(`${API_BASE_URL}/api/v1/deploys`)
      .then(r => r.json())
      .then(data => { setDeploys(data.deploys ?? []); setLoading(false) })
      .catch(() => setLoading(false))

    const interval = setInterval(() => {
      fetch(`${API_BASE_URL}/api/v1/deploys`)
        .then(r => r.json())
        .then(data => setDeploys(data.deploys ?? []))
    }, 10000)

    return () => clearInterval(interval)
  }, [])

  return { deploys, loading }
}