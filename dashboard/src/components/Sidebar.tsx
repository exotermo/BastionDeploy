interface Props {
  page: string
  onNavigate: (page: string) => void
}

export function Sidebar({ page, onNavigate }: Props) {
  const navItems = [
    { id: 'dashboard', label: 'Dashboard', icon: '📊' },
    { id: 'deploys', label: 'Deploys', icon: '🚀' },
    { id: 'apps', label: 'Apps', icon: '🖥️' },
    { id: 'config', label: 'Configurações', icon: '⚙️' },
  ]

  return (
    <aside className="w-56 bg-[#13151c] border-r border-[#1e2130] flex flex-col shrink-0">
      <div className="p-5 border-b border-[#1e2130]">
        <div className="flex items-center gap-2.5">
          <div className="w-7 h-7 bg-cyan-500 rounded-lg flex items-center justify-center text-white text-xs font-bold">
            E
          </div>
          <span className="text-sm font-bold text-white">BastionDeploy</span>
        </div>
      </div>

      <nav className="flex-1 p-3">
        <ul className="space-y-0.5">
          {navItems.map((item) => (
            <li key={item.id}>
              <button
                onClick={() => onNavigate(item.id)}
                className={`w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg text-sm transition-all ${
                  page === item.id
                    ? 'bg-cyan-500/10 text-cyan-400 font-medium'
                    : 'text-gray-500 hover:text-gray-200 hover:bg-[#1a1c27]'
                }`}
              >
                <span className="text-base">{item.icon}</span>
                <span>{item.label}</span>
              </button>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  )
}
