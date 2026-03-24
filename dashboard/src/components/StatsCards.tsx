import { TrendingUp, CheckCircle2, Package, Clock } from 'lucide-react'
import { useStats } from '../hooks/useStats'

export function StatsCards() {
  const { stats, loading } = useStats()

  const formatLastDeploy = (dt: string | null) => {
    if (!dt) return '—'
    const diff = Math.floor((Date.now() - new Date(dt).getTime()) / 1000)
    if (diff < 60) return `${diff}s ago`
    if (diff < 3600) return `${Math.floor(diff / 60)} min ago`
    return `${Math.floor(diff / 3600)}h ago`
  }

  const cards = [
    {
      icon: TrendingUp,
      label: 'Total Deploys',
      value: loading ? '...' : String(stats?.total_deploys ?? 0),
      color: '#00d4ff',
    },
    {
      icon: CheckCircle2,
      label: 'Success Rate',
      value: loading ? '...' : `${(stats?.success_rate ?? 0).toFixed(1)}%`,
      color: '#00ff88',
    },
    {
      icon: Package,
      label: 'Active Apps',
      value: loading ? '...' : String(stats?.active_apps ?? 0),
      color: '#00d4ff',
    },
    {
      icon: Clock,
      label: 'Last Deploy',
      value: loading ? '...' : formatLastDeploy(stats?.last_deploy_at ?? null),
      color: '#ffaa00',
    },
  ]

  return (
    <div className="grid grid-cols-4 gap-6">
      {cards.map((stat) => (
        <div
          key={stat.label}
          className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl p-6 hover:border-[#ffffff20] transition-all"
          style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}
        >
          <div className="flex items-start justify-between">
            <div>
              <p className="text-gray-400 text-sm mb-2">{stat.label}</p>
              <p className="text-3xl font-bold text-white">{stat.value}</p>
            </div>
            <div className="p-3 rounded-lg" style={{ backgroundColor: `${stat.color}15` }}>
              <stat.icon className="w-6 h-6" style={{ color: stat.color }} />
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}