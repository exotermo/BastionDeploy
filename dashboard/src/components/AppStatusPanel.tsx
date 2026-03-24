import { Circle } from 'lucide-react'
import { useAppsStatus } from '../hooks/useAppsStatus'

export function AppStatusPanel() {
  const { apps, loading } = useAppsStatus()

  return (
    <div className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl p-6"
         style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}>
      <h3 className="text-lg font-semibold text-white mb-6">App Status</h3>
      {loading ? (
        <p className="text-gray-500 text-sm">Carregando...</p>
      ) : apps.length === 0 ? (
        <p className="text-gray-500 text-sm">Nenhuma app encontrada.</p>
      ) : (
        <div className="space-y-4">
          {apps.map((app) => (
            <div key={app.name}
                 className="flex items-center justify-between p-3 rounded-lg bg-[#ffffff05] border border-[#ffffff08] hover:border-[#ffffff15] transition-all">
              <div className="flex items-center gap-3">
                <Circle className={`w-3 h-3 ${app.status === 'UP' ? 'fill-[#00ff88] text-[#00ff88]' : 'fill-[#ff4444] text-[#ff4444]'}`} />
                <span className="text-white font-medium" style={{ fontFamily: 'JetBrains Mono, monospace' }}>
                  {app.name}
                </span>
              </div>
              <span className={`px-2 py-1 rounded text-xs font-medium ${
                app.status === 'UP'
                  ? 'bg-[#00ff8815] text-[#00ff88] border border-[#00ff8830]'
                  : 'bg-[#ff444415] text-[#ff4444] border border-[#ff444430]'
              }`}>{app.status}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}