import { Bot, AlertCircle } from 'lucide-react'

export function AIAssistant() {
  return (
    <div
      className="bg-gradient-to-br from-[#0a0c10] to-[#0f1117] border border-[#00d4ff30] rounded-xl p-6 relative overflow-hidden"
      style={{ boxShadow: '0 0 40px rgba(0,212,255,0.15), 0 4px 24px rgba(0,0,0,0.4)' }}
    >
      <div className="absolute -top-20 -right-20 w-40 h-40 bg-[#00d4ff] rounded-full blur-[100px] opacity-20" />
      <div className="relative z-10">
        <div className="flex items-start gap-4">
          <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-[#00d4ff] to-[#0088ff] rounded-lg flex items-center justify-center shadow-lg shadow-[#00d4ff30]">
            <Bot className="w-6 h-6 text-white" />
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-3">
              <h3 className="text-white font-semibold">AI Agent</h3>
              <span className="text-xs text-[#00d4ff] bg-[#00d4ff15] px-2 py-0.5 rounded-full border border-[#00d4ff30]">
                Active
              </span>
            </div>
            <div className="bg-[#0f1117] border border-[#ffffff10] rounded-lg p-4 text-sm" style={{ fontFamily: 'JetBrains Mono, monospace' }}>
              <div className="flex items-start gap-2 mb-2">
                <AlertCircle className="w-4 h-4 text-[#ff4444] flex-shrink-0 mt-0.5" />
                <div className="text-gray-300">
                  <span className="text-[#ff4444] font-medium">Deploy failed</span> on{' '}
                  <span className="text-[#00d4ff]">meu-bot</span>
                </div>
              </div>
              <div className="pl-6 space-y-2">
                <p className="text-gray-400">
                  Missing environment variable{' '}
                  <code className="text-[#00ff88] bg-[#00ff8815] px-1.5 py-0.5 rounded">DATABASE_URL</code>
                </p>
                <div className="border-l-2 border-[#00d4ff30] pl-3 mt-3">
                  <p className="text-[#00d4ff] text-xs font-medium mb-1">💡 Suggestion:</p>
                  <p className="text-gray-300">
                    Add it to your <code className="text-gray-400">.env</code> file and redeploy.
                  </p>
                </div>
              </div>
            </div>
            <div className="mt-3 text-xs text-gray-500">Last analysis: 2 minutes ago</div>
          </div>
        </div>
      </div>
    </div>
  )
}
