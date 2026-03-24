import { Sidebar } from './components/Sidebar'
import { StatsCards } from './components/StatsCards'
import { DeploymentsTable } from './components/DeploymentsTable'
import { DeployChart } from './components/DeployChart'
import { AppStatusPanel } from './components/AppStatusPanel'
import { AIAssistant } from './components/AIAssistant'

export default function App() {
  return (
    <div className="flex h-screen bg-[#0f1117] overflow-hidden">
      <Sidebar />
      <div className="flex-1 overflow-y-auto">
        <div className="p-8 max-w-[1800px]">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-white mb-2">Dashboard</h1>
            <p className="text-gray-400">Monitor your deployments and apps in real-time</p>
          </div>
          <div className="mb-8">
            <StatsCards />
          </div>
          <div className="mb-8">
            <AIAssistant />
          </div>
          <div className="mb-8">
            <DeploymentsTable />
          </div>
          <div className="grid grid-cols-2 gap-6">
            <DeployChart />
            <AppStatusPanel />
          </div>
        </div>
      </div>
    </div>
  )
}
