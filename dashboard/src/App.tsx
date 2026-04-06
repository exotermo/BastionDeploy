import { useState } from 'react'
import { Sidebar } from './components/Sidebar'
import { StatsCards } from './components/StatsCards'
import { DeploymentsTable } from './components/DeploymentsTable'
import { DeployChart } from './components/DeployChart'
import { AppStatusPanel } from './components/AppStatusPanel'
import { AIAssistant } from './components/AIAssistant'

export default function App() {
  return (
    <div className="flex h-screen bg-[#0e1015] overflow-hidden">
      <Sidebar />
      <div className="flex-1 overflow-y-auto">
        <div className="p-7 max-w-[1800px]">
          {/* Header */}
          <div className="mb-7">
            <h1 className="text-xl font-bold text-white">Dashboard</h1>
            <p className="text-[13px] text-gray-500 mt-0.5">
              Monitore deploys e aplicativos em tempo real
            </p>
          </div>

          {/* Stats */}
          <div className="mb-6">
            <StatsCards />
          </div>

          {/* AI Assistant */}
          <div className="mb-6">
            <AIAssistant />
          </div>

          {/* Deployments Table */}
          <div className="mb-6">
            <DeploymentsTable />
          </div>

          {/* Chart + Apps */}
          <div className="grid grid-cols-2 gap-6">
            <DeployChart />
            <AppStatusPanel />
          </div>

          {/* Footer */}
          <div className="mt-10 mb-4 pt-6 border-t border-[#1a1c27] flex items-center justify-between">
            <p className="text-[11px] text-gray-700">
              BastionDeploy — Secure Self-Hosted PaaS
            </p>
            <p className="text-[11px] text-gray-700 font-mono">v0.2.0</p>
          </div>
        </div>
      </div>
    </div>
  )
}
