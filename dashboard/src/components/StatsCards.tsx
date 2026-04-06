import { TrendingUp, CheckCircle2, Monitor, Clock } from 'lucide-react'
import { useStats } from '../hooks/useStats'

export function StatsCards() {
  const { stats, loading } = useStats()

  const formatLastDeploy = (dt: string | null) => {
    if (!dt) return '—'
    const diff = Math.floor((Date.now() - new Date(dt).getTime()) / 1000)
    if (diff < 60) return `${diff}s atrás`
    if (diff < 3600) return `${Math.floor(diff / 60)} min atrás`
    return `${Math.floor(diff / 3600)}h atrás`
  }

  const accentTheme: Record<
    string,
    { border: string; text: string; iconBg: string }
  > = {
    cyan: {
      border: 'border-l-cyan-500',
      text: 'text-cyan-500',
      iconBg: 'bg-cyan-500/10 text-cyan-500',
    },
    emerald: {
      border: 'border-l-emerald-500',
      text: 'text-emerald-500',
      iconBg: 'bg-emerald-500/10 text-emerald-500',
    },
    violet: {
      border: 'border-l-violet-500',
      text: 'text-violet-500',
      iconBg: 'bg-violet-500/10 text-violet-500',
    },
    amber: {
      border: 'border-l-amber-500',
      text: 'text-amber-500',
      iconBg: 'bg-amber-500/10 text-amber-500',
    },
  }

  const cards = [
    {
      icon: TrendingUp,
      label: 'Total Deploys',
      value: loading ? '—' : String(stats?.total_deploys ?? 0),
      accent: 'cyan' as const,
    },
    {
      icon: CheckCircle2,
      label: 'Taxa de Sucesso',
      value: loading ? '—' : `${(stats?.success_rate ?? 0).toFixed(1)}%`,
      accent: 'emerald' as const,
    },
    {
      icon: Monitor,
      label: 'Apps Ativas',
      value: loading ? '—' : String(stats?.active_apps ?? 0),
      accent: 'violet' as const,
    },
    {
      icon: Clock,
      label: 'Último Deploy',
      value: loading ? '—' : formatLastDeploy(stats?.last_deploy_at ?? null),
      accent: 'amber' as const,
    },
  ]

  return (
    <div className="grid grid-cols-4 gap-5">
      {cards.map((stat) => {
        const theme = accentTheme[stat.accent]
        return (
          <div
            key={stat.label}
            className={`bg-[#161922] border border-[#1e2130] rounded-xl pl-4 pr-5 py-4 border-l-4 ${theme.border}`}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-[11px] text-gray-500 font-semibold uppercase tracking-wide">
                  {stat.label}
                </p>
                <p className="text-2xl font-bold text-white mt-1">{stat.value}</p>
              </div>
              <div className={`p-2.5 rounded-lg ${theme.iconBg}`}>
                <stat.icon className="w-5 h-5" />
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
