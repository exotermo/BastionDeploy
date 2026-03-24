import { GitBranch, User, Clock } from 'lucide-react'

interface Deployment {
  app: string
  branch: string
  commit: string
  status: 'SUCCESS' | 'FAILED' | 'PENDING' | 'RUNNING'
  triggeredBy: string
  time: string
}

const deployments: Deployment[] = [
  { app: 'meu-bot',      branch: 'main', commit: 'abc123', status: 'SUCCESS', triggeredBy: 'xande', time: '2 min ago'  },
  { app: 'api-finance',  branch: 'dev',  commit: 'def456', status: 'FAILED',  triggeredBy: 'xande', time: '1 hour ago' },
  { app: 'dashboard',    branch: 'main', commit: 'ghi789', status: 'PENDING', triggeredBy: 'xande', time: 'just now'   },
  { app: 'telegram-bot', branch: 'main', commit: 'jkl012', status: 'RUNNING', triggeredBy: 'xande', time: '30s ago'    },
]

const statusStyle: Record<Deployment['status'], string> = {
  SUCCESS: 'bg-[#00ff8815] text-[#00ff88] border-[#00ff8830]',
  FAILED:  'bg-[#ff444415] text-[#ff4444] border-[#ff444430]',
  PENDING: 'bg-[#ffaa0015] text-[#ffaa00] border-[#ffaa0030]',
  RUNNING: 'bg-[#00d4ff15] text-[#00d4ff] border-[#00d4ff30] animate-pulse',
}

export function DeploymentsTable() {
  return (
    <div className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl overflow-hidden"
         style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}>
      <div className="p-6 border-b border-[#ffffff10]">
        <h2 className="text-xl font-semibold text-white">Recent Deployments</h2>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[#ffffff10]">
              {['App','Branch','Commit','Status','Triggered by','Time'].map(h => (
                <th key={h} className="text-left px-6 py-4 text-sm font-medium text-gray-400">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {deployments.map((d, i) => (
              <tr key={i} className="border-b border-[#ffffff08] hover:bg-[#ffffff05] transition-colors">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2">
                    <div className="w-8 h-8 bg-[#00d4ff15] rounded-lg flex items-center justify-center">
                      <GitBranch className="w-4 h-4 text-[#00d4ff]" />
                    </div>
                    <span className="text-white font-medium">{d.app}</span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span className="text-gray-300" style={{ fontFamily: 'JetBrains Mono, monospace' }}>{d.branch}</span>
                </td>
                <td className="px-6 py-4">
                  <code className="text-gray-400 text-sm bg-[#ffffff08] px-2 py-1 rounded" style={{ fontFamily: 'JetBrains Mono, monospace' }}>{d.commit}</code>
                </td>
                <td className="px-6 py-4">
                  <span className={`px-3 py-1 rounded-full text-xs font-medium border ${statusStyle[d.status]}`}>{d.status}</span>
                </td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2 text-gray-300">
                    <User className="w-4 h-4" />{d.triggeredBy}
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2 text-gray-400">
                    <Clock className="w-4 h-4" />{d.time}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
