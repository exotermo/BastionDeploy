import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

const data = [
  { day: 'Mon', deploys: 5 },
  { day: 'Tue', deploys: 8 },
  { day: 'Wed', deploys: 6 },
  { day: 'Thu', deploys: 10 },
  { day: 'Fri', deploys: 7 },
  { day: 'Sat', deploys: 3 },
  { day: 'Sun', deploys: 4 },
]

export function DeployChart() {
  return (
    <div className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl p-6"
         style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}>
      <h3 className="text-lg font-semibold text-white mb-6">Deploys Last 7 Days</h3>
      <ResponsiveContainer width="100%" height={240}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#ffffff10" />
          <XAxis dataKey="day" stroke="#6b7280" style={{ fontSize: '12px' }} />
          <YAxis stroke="#6b7280" style={{ fontSize: '12px' }} />
          <Tooltip
            contentStyle={{ backgroundColor: '#0a0c10', border: '1px solid #ffffff20', borderRadius: '8px', color: '#fff' }}
            cursor={{ fill: '#ffffff08' }}
          />
          <Bar dataKey="deploys" fill="#00d4ff" radius={[8, 8, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
