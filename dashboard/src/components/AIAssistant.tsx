import { Bot, AlertTriangle } from 'lucide-react'

export function AIAssistant() {
  return (
    <div className="bg-[#161922] border border-[#1e2130] rounded-xl p-5">
      <div className="flex items-start gap-4">
        <div className="flex-shrink-0 w-9 h-9 bg-[#1a1c27] rounded-lg flex items-center justify-center mt-0.5">
          <Bot className="w-5 h-5 text-gray-500" />
        </div>
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <h3 className="text-sm font-semibold text-white">AI Assistant</h3>
            <span className="text-[10px] text-gray-600 bg-[#1a1c27] px-2 py-0.5 rounded-md">
              Em breve
            </span>
          </div>
          <p className="text-[13px] text-gray-500">
            O assistente IA analisará logs de deploy com falha e sugerirá correções automáticas.
            Integrado com LLMs para diagnóstico inteligente.
          </p>
          <div className="mt-3 flex items-center gap-2 text-[11px] text-gray-700">
            <AlertTriangle className="w-3 h-3" />
            Funcionalidade em desenvolvimento
          </div>
        </div>
      </div>
    </div>
  )
}
