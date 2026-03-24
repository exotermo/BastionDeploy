import { Rocket, LayoutDashboard, GitBranch, Package, ScrollText, Settings, User } from 'lucide-react'

export function Sidebar() {
  const navItems = [
    { icon: LayoutDashboard, label: 'Dashboard', active: true },
    { icon: GitBranch, label: 'Deployments', active: false },
    { icon: Package, label: 'Apps', active: false },
    { icon: ScrollText, label: 'Logs', active: false },
    { icon: Settings, label: 'Settings', active: false },
  ]

  return (
    <div className="w-64 bg-[#0a0c10] border-r border-[#ffffff10] flex flex-col h-screen">
      <div className="p-6 border-b border-[#ffffff10]">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-gradient-to-br from-[#00d4ff] to-[#0088ff] rounded-lg flex items-center justify-center">
            <Rocket className="w-5 h-5 text-white" />
          </div>
          <span className="text-xl font-semibold text-white">BastionDeploy</span>
        </div>
      </div>
      <nav className="flex-1 p-4">
        <ul className="space-y-1">
          {navItems.map((item) => (
            <li key={item.label}>
              <button
                className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${
                  item.active
                    ? 'bg-[#00d4ff15] text-[#00d4ff] border border-[#00d4ff30]'
                    : 'text-gray-400 hover:text-white hover:bg-[#ffffff08]'
                }`}
              >
                <item.icon className="w-5 h-5" />
                <span className="font-medium">{item.label}</span>
              </button>
            </li>
          ))}
        </ul>
      </nav>
      <div className="p-4 border-t border-[#ffffff10]">
        <div className="flex items-center gap-3 px-4 py-3 rounded-lg bg-[#ffffff05] border border-[#ffffff08]">
          <div className="w-10 h-10 bg-gradient-to-br from-[#00d4ff] to-[#0088ff] rounded-full flex items-center justify-center">
            <User className="w-5 h-5 text-white" />
          </div>
          <div>
            <p className="text-white font-medium text-sm">xande</p>
            <p className="text-gray-500 text-xs">admin@bastiondeploy</p>
          </div>
        </div>
      </div>
    </div>
  )
}
