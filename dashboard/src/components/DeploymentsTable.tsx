import { GitBranch, User, Clock } from 'lucide-react'
import { useDeploys, Deploy } from '../hooks/useDeploys'

type UIStatus = 'SUCCESS' | 'FAILED' | 'PENDING' | 'RUNNING'

const toUIStatus = (s: Deploy['status']): UIStatus =>
  s.toUpperCase() as UIStatus

const statusStyle: Record<UIStatus, string> = {
  SUCCESS: 'bg-[#00ff8815] text-[#00ff88] border-[#00ff8830]',
  FAILED:  'bg-[#ff444415] text-[#ff4444] border-[#ff444430]',
  PENDING: 'bg-[#ffaa0015] text-[#ffaa00] border-[#ffaa0030]',
  RUNNING: 'bg-[#00d4ff15] text-[#00d4ff] border-[#00d4ff30] animate-pulse',
}

const formatTime = (dt: string) => {
  const diff = Math.floor((Date.now() - new Date(dt).getTime()) / 1000)
  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)} min ago`
  return `${Math.floor(diff / 3600)}h ago`
}

export function DeploymentsTable() {
  const { deploys, loading } = useDeploys()

  return (
    <div className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl overflow-hidden"
         style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}>
      <div className="p-6 border-b border-[#ffffff10]">
        <h2 className="text-xl font-semibold text-white">Recent Deployments</h2>
      </div>
      <div className="overflow-x-auto">
        {loading ? (
          <div className="p-8 text-center text-gray-500">Carregando deploys...</div>
        ) : deploys.length === 0 ? (
          <div className="p-8 text-center text-gray-500">Nenhum deploy encontrado.</div>
        ) : (
          <table className="w-full">
            <thead>
              <tr className="border-b border-[#ffffff10]">
                {['App','Branch','Commit','Status','Triggered by','Time'].map(h => (
                  <th key={h} className="text-left px-6 py-4 text-sm font-medium text-gray-400">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {deploys.map((d) => {
                const uiStatus = toUIStatus(d.status)
                return (
                  <tr key={d.id} className="border-b border-[#ffffff08] hover:bg-[#ffffff05] transition-colors">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2">
                        <div className="w-8 h-8 bg-[#00d4ff15] rounded-lg flex items-center justify-center">
                          <GitBranch className="w-4 h-4 text-[#00d4ff]" />
                        </div>
                        <span className="text-white font-medium">{d.app_name}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-gray-300" style={{ fontFamily: 'JetBrains Mono, monospace' }}>
                        {d.branch.replace('refs/heads/', '')}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <code className="text-gray-400 text-sm bg-[#ffffff08] px-2 py-1 rounded"
                            style={{ fontFamily: 'JetBrains Mono, monospace' }}>
                        {d.commit_sha.slice(0, 7)}
                      </code>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-3 py-1 rounded-full text-xs font-medium border ${statusStyle[uiStatus]}`}>
                        {uiStatus}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2 text-gray-300">
                        <User className="w-4 h-4" />{d.triggered_by}
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2 text-gray-400">
                        <Clock className="w-4 h-4" />{formatTime(d.created_at)}
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}