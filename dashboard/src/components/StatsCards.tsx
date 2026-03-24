import { TrendingUp, CheckCircle2, Package, Clock } from 'lucide-react'

export function StatsCards() {
  const stats = [
    { icon: TrendingUp,   label: 'Total Deploys', value: '42',       color: '#00d4ff' },
    { icon: CheckCircle2, label: 'Success Rate',  value: '94%',      color: '#00ff88' },
    { icon: Package,      label: 'Active Apps',   value: '7',        color: '#00d4ff' },
    { icon: Clock,        label: 'Last Deploy',   value: '2 min ago',color: '#ffaa00' },
  ]

  return (
    <div className="grid grid-cols-4 gap-6">
      {stats.map((stat) => (
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
