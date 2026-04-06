import { useState } from 'react'
import { Sidebar } from './components/Sidebar'
import { StatsCards } from './components/StatsCards'
import { DeploymentsTable } from './components/DeploymentsTable'
import { DeployChart } from './components/DeployChart'
import { AppStatusPanel } from './components/AppStatusPanel'
import { AIAssistant } from './components/AIAssistant'
import { DeployPage } from './pages/DeployPage'
import { AppsPage } from './pages/AppsPage'
import { ConfigPage } from './pages/ConfigPage'

export default function App() {
  const [page, setPage] = useState('dashboard')

  const renderPage = () => {
    switch (page) {
      case 'deploys':
        return <DeployPage />
      case 'apps':
        return <AppsPage />
      case 'config':
        return <ConfigPage />
      default:
        return (
          <>
            <div className="mb-7">
              <h1 className="text-xl font-bold text-white">Dashboard</h1>
              <p className="text-[13px] text-gray-500 mt-0.5">Visão geral do sistema</p>
            </div>
            <div className="mb-6"><StatsCards /></div>
            <div className="mb-6"><AIAssistant /></div>
            <div className="mb-6"><DeploymentsTable /></div>
            <div className="grid grid-cols-2 gap-6">
              <DeployChart />
              <AppStatusPanel />
            </div>
          </>
        )
    }
  }

  return (
    <div className="flex h-screen bg-[#0e1015] overflow-hidden">
      <Sidebar page={page} onNavigate={setPage} />
      <div className="flex-1 overflow-y-auto">
        <div className="p-7 max-w-[1800px]">
          {renderPage()}
        </div>
      </div>
    </div>
  )
}
