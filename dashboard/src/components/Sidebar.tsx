import { Rocket, LayoutDashboard, GitBranch, Monitor, Settings, LogOut } from 'lucide-react'

export function Sidebar() {
  const navItems = [
    { icon: LayoutDashboard, label: 'Dashboard', active: true },
    { icon: GitBranch, label: 'Deploys', active: false },
    { icon: Monitor, label: 'Apps', active: false },
    { icon: Settings, label: 'Configurações', active: false },
  ]

  return (
    <aside className="w-60 bg-[#13151c] border-r border-[#1e2130] flex flex-col shrink-0">
      <div className="p-5 border-b border-[#1e2130]">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-cyan-500 rounded-lg flex items-center justify-center">
            <Rocket className="w-4 h-4 text-white" />
          </div>
          <div className="leading-none">
            <span className="text-base font-bold text-white block">BastionDeploy</span>
            <span className="text-[10px] text-gray-500 font-medium tracking-wide">v0.2.0</span>
          </div>
        </div>
      </div>

      <nav className="flex-1 p-3">
        <div className="text-[10px] text-gray-600 font-semibold uppercase tracking-wider px-3 mb-2">
          Navegação
        </div>
        <ul className="space-y-0.5">
          {navItems.map((item) => (
            <li key={item.label}>
              <button
                className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all ${
                  item.active
                    ? 'bg-cyan-500/10 text-cyan-400 font-medium'
                    : 'text-gray-500 hover:text-gray-200 hover:bg-[#1a1c27]'
                }`}
              >
                <item.icon className="w-4 h-4 shrink-0" />
                <span>{item.label}</span>
              </button>
            </li>
          ))}
        </ul>
      </nav>

      <div className="p-3 border-t border-[#1e2130]">
        <div className="flex items-center gap-2.5 px-3 py-2.5 rounded-lg bg-[#1a1c27]">
          <div className="w-8 h-8 bg-[#222639] rounded-full flex items-center justify-center text-[11px] text-gray-400 font-semibold">
            X
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm text-gray-200 truncate">xande</p>
            <p className="text-[11px] text-gray-600">Admin</p>
          </div>
          <LogOut className="w-3.5 h-3.5 text-gray-600 hover:text-gray-300 cursor-pointer shrink-0" />
        </div>
      </div>
    </aside>
  )
}
