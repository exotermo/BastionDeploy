import { useState } from 'react'
import { useStats } from '../hooks/useStats'

export function ConfigPage() {
  const { stats } = useStats()

  return (
    <>
      <div className="mb-6">
        <h1 className="text-xl font-bold text-white">Configurações</h1>
        <p className="text-[13px] text-gray-500 mt-0.5">Informações do sistema</p>
      </div>

      <div className="space-y-4">
        {/* System */}
        <div className="bg-[#161922] border border-[#1e2130] rounded-xl p-5">
          <h3 className="text-sm font-semibold text-white mb-4">Sistema</h3>
          <div className="grid grid-cols-2 gap-x-8 gap-y-3 text-sm">
            <div className="flex justify-between border-b border-[#1e2130] pb-2">
              <span className="text-gray-500">Versão</span>
              <span className="text-gray-300 font-mono">0.2.0</span>
            </div>
            <div className="flex justify-between border-b border-[#1e2130] pb-2">
              <span className="text-gray-500">Total Deploys</span>
              <span className="text-gray-300 font-mono">{stats?.total_deploys ?? '—'}</span>
            </div>
            <div className="flex justify-between border-b border-[#1e2130] pb-2">
              <span className="text-gray-500">Taxa de Sucesso</span>
              <span className="text-gray-300 font-mono">{stats?.success_rate?.toFixed(1) ?? '—'}%</span>
            </div>
            <div className="flex justify-between border-b border-[#1e2130] pb-2">
              <span className="text-gray-500">Apps Ativas</span>
              <span className="text-gray-300 font-mono">{stats?.active_apps ?? '—'}</span>
            </div>
          </div>
        </div>

        {/* Services */}
        <div className="bg-[#161922] border border-[#1e2130] rounded-xl p-5">
          <h3 className="text-sm font-semibold text-white mb-4">Serviços</h3>
          <div className="space-y-2">
            {[
              { name: 'API', desc: 'Go + Gin', status: 'online' },
              { name: 'Agent', desc: 'Go + Redis Worker', status: 'online' },
              { name: 'PostgreSQL', desc: 'Postgres 16', status: 'online' },
              { name: 'Redis', desc: 'Redis 7', status: 'online' },
            ].map((svc) => (
              <div
                key={svc.name}
                className="flex items-center justify-between py-2.5 border-b border-[#1e2130] last:border-0"
              >
                <div>
                  <span className="text-sm text-gray-200 font-medium">{svc.name}</span>
                  <span className="text-[11px] text-gray-600 ml-2">{svc.desc}</span>
                </div>
                <span className="text-xs text-emerald-500 bg-emerald-500/10 px-2 py-0.5 rounded-md font-medium">
                  {svc.status}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Info */}
        <div className="bg-[#161922] border border-[#1e2130] rounded-xl p-5">
          <h3 className="text-sm font-semibold text-white mb-4">Documentação</h3>
          <div className="text-sm text-gray-500 space-y-2">
            <p>Setup completo: <code className="text-gray-400">setup.md</code> na raiz do projeto</p>
            <p>CLI: <code className="text-gray-400">./cli/exoctl setup</code></p>
            <p>Source: <a href="https://github.com/exotermo/BastionDeploy" className="text-cyan-500 underline">exotermo/BastionDeploy</a></p>
          </div>
        </div>
      </div>
    </>
  )
}
