import { useAppsStatus } from '../hooks/useAppsStatus'

export function AppsPage() {
  const { apps, loading } = useAppsStatus()

  return (
    <>
      <div className="mb-6">
        <h1 className="text-xl font-bold text-white">Aplicativos</h1>
        <p className="text-[13px] text-gray-500 mt-0.5">Status dos aplicativos implantados</p>
      </div>

      {loading ? (
        <div className="text-center py-12 text-gray-600">Carregando...</div>
      ) : apps.length === 0 ? (
        <div className="bg-[#161922] border border-[#1e2130] rounded-xl p-12 text-center">
          <p className="text-gray-500">Nenhuma app registrada.</p>
        </div>
      ) : (
        <div className="space-y-2">
          {apps.map((app) => (
            <div
              key={app.name}
              className="bg-[#161922] border border-[#1e2130] rounded-xl flex items-center justify-between p-4"
            >
              <div className="flex items-center gap-3">
                <div
                  className={`w-2.5 h-2.5 rounded-full ${
                    app.status === 'UP' ? 'bg-emerald-500' : 'bg-red-500'
                  }`}
                />
                <div>
                  <p className="text-sm text-gray-200 font-mono font-medium">{app.name}</p>
                  <p className="text-[11px] text-gray-600">{app.uptime}</p>
                </div>
              </div>
              <span
                className={`text-xs font-semibold px-2.5 py-1 rounded-md ${
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
    </>
  )
}
