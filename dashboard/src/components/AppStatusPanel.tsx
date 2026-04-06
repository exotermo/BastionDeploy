import { Server, ExternalLink } from 'lucide-react'
import { useAppsStatus } from '../hooks/useAppsStatus'

export function AppStatusPanel() {
  const { apps, loading } = useAppsStatus()

  return (
    <div className="bg-[#161922] border border-[#1e2130] rounded-xl p-5">
      <h3 className="text-sm font-semibold text-white mb-1">Apps</h3>
      <p className="text-[11px] text-gray-600 mb-4">Status dos aplicativos</p>
      {loading ? (
        <p className="text-gray-700 text-sm">Carregando...</p>
      ) : apps.length === 0 ? (
        <p className="text-gray-700 text-sm">Nenhuma app registrada.</p>
      ) : (
        <div className="space-y-2">
          {apps.map((app) => (
            <div
              key={app.name}
              className="flex items-center justify-between px-4 py-3 rounded-lg bg-[#111319] border border-[#1e2130] hover:border-[#2a2d42] transition-all"
            >
              <div className="flex items-center gap-3">
                <div
                  className={`w-2 h-2 rounded-full mt-0.5 ${
                    app.status === 'UP' ? 'bg-emerald-500' : 'bg-red-500'
                  }`}
                />
                <div>
                  <span className="text-sm text-gray-200 font-mono font-medium">
                    {app.name}
                  </span>
                  <p className="text-[10px] text-gray-600">{app.uptime}</p>
                </div>
              </div>
              <span
                className={`text-[10px] font-bold px-2 py-0.5 rounded-md ${
                  app.status === 'UP'
                    ? 'bg-emerald-500/10 text-emerald-500'
                    : 'bg-red-500/10 text-red-500'
                }`}
              >
                {app.status}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
