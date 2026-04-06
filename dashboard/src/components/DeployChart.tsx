import { useState, useEffect } from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { API_BASE_URL } from '../config/api'

export function DeployChart() {
  const [data, setData] = useState<{ day: string; deploys: number; success: number }[]>([])

  useEffect(() => {
    fetch(`${API_BASE_URL}/api/v1/deploys`)
      .then((r) => r.json())
      .then((d) => {
        const deploys = d.deploys ?? []
        const last7 = Array.from({ length: 7 }, (_, i) => {
          const date = new Date()
          date.setDate(date.getDate() - (6 - i))
          date.setHours(0, 0, 0, 0)
          const dayStr = date.toLocaleDateString('pt-BR', { weekday: 'short' })
          const dayEnd = new Date(date)
          dayEnd.setDate(dayEnd.getDate() + 1)

          const dayDeploys = deploys.filter(
            (p: { created_at: string }) =>
              new Date(p.created_at) >= date && new Date(p.created_at) < dayEnd
          )
          return {
            day: dayStr.charAt(0).toUpperCase() + dayStr.slice(1, 3),
            deploys: dayDeploys.length,
            success: dayDeploys.filter(
              (p: { status: string }) => p.status === 'success'
            ).length,
          }
        })
        setData(last7)
      })
  }, [])

  return (
    <div className="bg-[#161922] border border-[#1e2130] rounded-xl p-5">
      <h3 className="text-sm font-semibold text-white mb-1">Deploys — Últimos 7 dias</h3>
      <p className="text-[11px] text-gray-600 mb-5">Volume total e taxa de sucesso</p>
      <ResponsiveContainer width="100%" height={220}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#1e2130" vertical={false} />
          <XAxis dataKey="day" stroke="#4a4e5c" fontSize={11} tickLine={false} axisLine={false} />
          <YAxis stroke="#4a4e5c" fontSize={11} tickLine={false} axisLine={false} />
          <Tooltip
            contentStyle={{
              backgroundColor: '#161922',
              border: '1px solid #1e2130',
              borderRadius: '8px',
              color: '#e4e6ed',
              fontSize: '12px',
            }}
            cursor={{ fill: '#1a1c27' }}
          />
          <Bar
            dataKey="deploys"
            fill="#7c3aed"
            radius={[6, 6, 0, 0]}
            maxBarSize={36}
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
