import { useMemo, useState } from 'react';
import { useGetPaymentAnalyticQuery } from '../Redux/features/analytics/analytics';

{/*icon */}
const Icon = ({ path, className = 'w-5 h-5' }) => (
  <svg
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="1.6"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
    aria-hidden="true"
  >
    {path}
  </svg>
)

const icons = {
  filter: (
    <>
      <path d="M4 6h16" />
      <circle cx="9" cy="6" r="1.8" fill="currentColor" stroke="none" />
      <path d="M4 12h16" />
      <circle cx="15" cy="12" r="1.8" fill="currentColor" stroke="none" />
      <path d="M4 18h16" />
      <circle cx="11" cy="18" r="1.8" fill="currentColor" stroke="none" />
    </>
  ),
  chevronDown: <path d="M6 9l6 6 6-6" />,
  inbox: (
    <>
      <path d="M3.5 12h5l1.5 3h4l1.5-3h5" />
      <path d="M5 6.5h14L21 12v6a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18v-6l2-5.5z" />
    </>
  ),
}

{/*nodeicon*/}
const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

const statusMeta = {
  PAID: { label: 'Paid', bg: '#DCFCE7', text: '#15803D', dot: '#16A34A' },
  PENDING: { label: 'Pending', bg: '#FEF3C7', text: '#B45309', dot: '#B45309' },
  FAILED: { label: 'Failed', bg: '#FEE2E2', text: '#DC2626', dot: '#DC2626' },
  REFUNDED: { label: 'Refunded', bg: '#F1F5F9', text: '#475569', dot: '#94A3B8' },
}

// sample — the one event from your log, so this page renders standalone
const sampleEvents = [
  {
    OrderID: 81,
    PaymentID: 75,
    UserID: 7,
    Email: 'godvenus419@gmail.com',
    Status: 'PAID',
    CreatedAt: '2026-08-26 13:05:18.541738 +0000 UTC',
  },
]

// reads a field off an event whether it's PascalCase (Go struct dump)
// or snake_case (typical JSON over the wire)
function field(event, pascalKey, snakeKey) {
  return event[pascalKey] ?? event[snakeKey]
}

{/*format */}
function formatTimestamp(value) {
  if (!value) return '—'
  // "2026-08-26 13:05:18.541738 +0000 UTC" -> parseable ISO-ish string
  const parsed = new Date(value.replace(' ', 'T').replace(' UTC', 'Z').replace('+0000', '+00:00'))
  if (Number.isNaN(parsed.getTime())) return value
  return parsed.toLocaleString('en-IN', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function StatCard({ label, value, accent }) {
  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-5">
      <p className="FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">{label}</p>
      <p className="mt-1.5 FleetOps-display text-2xl font-semibold" style={{ color: accent || '#0B0E14' }}>
        {value}
      </p>
    </div>
  )
}

function SelectField({ icon, value, onChange, options, label }) {
  return (
    <label className="relative flex items-center">
      <span className="sr-only">{label}</span>
      <span className="pointer-events-none absolute left-3 text-[#5B6472]">{NodeIcon(icon, 'h-4 w-4')}</span>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="appearance-none rounded-full border border-slate-200 bg-white py-2.5 pl-9 pr-8 FleetOps-body text-[13px] font-medium text-[#0B0E14] transition-colors hover:border-slate-300 focus:border-[#0B0E14] focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/10"
      >
        {options.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
      <span className="pointer-events-none absolute right-3 text-[#5B6472]">{NodeIcon('chevronDown', 'h-3.5 w-3.5')}</span>
    </label>
  )
}

export default function Analytics({ events = sampleEvents }) {
  const [statusFilter, setStatusFilter] = useState('all')

  const { data, isLoading, error } = useGetPaymentAnalyticQuery();

  if (isLoading) {
    return <div>
        <p>Loading...</p>
    </div>
  }

  if (error) {
    return <div>
        <p>error...</p>
    </div>
  }

  const rows = events.map((e) => ({
    orderId: field(e, 'OrderID', 'order_id'),
    paymentId: field(e, 'PaymentID', 'payment_id'),
    userId: field(e, 'UserID', 'user_id'),
    email: field(e, 'Email', 'email'),
    status: field(e, 'Status', 'status'),
    createdAt: field(e, 'CreatedAt', 'created_at'),
  }))

  const stats = useMemo(
    () => ({
      total: rows.length,
      paid: rows.filter((r) => r.status === 'PAID').length,
      failed: rows.filter((r) => r.status === 'FAILED').length,
      uniqueUsers: new Set(rows.map((r) => r.userId)).size,
    }),
    [rows]
  )

  const statusOptions = useMemo(() => {
    const seen = new Set(rows.map((r) => r.status))
    return [{ value: 'all', label: 'All statuses' }, ...Array.from(seen).map((s) => ({ value: s, label: statusMeta[s]?.label || s }))]
  }, [rows])

  return (
    <div className="min-h-screen bg-[#F3F5F8] FleetOps-body text-[#0B0E14]">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@500;600;700&family=Manrope:wght@400;500;600;700&family=IBM+Plex+Mono:wght@400;500&display=swap');
        .FleetOps-display { font-family: 'Space Grotesk', ui-sans-serif, sans-serif; }
        .FleetOps-body { font-family: 'Manrope', ui-sans-serif, sans-serif; }
        .FleetOps-mono { font-family: 'IBM Plex Mono', ui-monospace, monospace; }
      `}</style>

      <main className="mx-auto max-w-6xl px-6 py-10">
        <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">analytics</span>
        <h1 className="mt-1 FleetOps-display text-2xl font-semibold tracking-tight text-[#0B0E14] sm:text-3xl">
          Payment events
        </h1>

        <div className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-4">
          <StatCard label="Total events" value={stats.total} />
          <StatCard label="Paid" value={stats.paid} accent="#15803D" />
          <StatCard label="Failed" value={stats.failed} accent="#DC2626" />
          <StatCard label="Unique users" value={stats.uniqueUsers} />
        </div>

        <div className="mt-8 flex items-center justify-between">
          <SelectField
            icon="filter"
            label="Filter by status"
            value={statusFilter}
            onChange={setStatusFilter}
            options={statusOptions}
          />
          <p className="FleetOps-mono text-[11px] text-[#5B6472]">
            {data.length} event{data.length === 1 ? '' : 's'}
          </p>
        </div>

        {data.length > 0 ? (
          <div className="mt-4 overflow-hidden rounded-2xl border border-slate-200 bg-white">
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead>
                  <tr className="border-b border-slate-100 bg-[#F9FAFB]">
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Order ID</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Payment ID</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">User ID</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Email</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Status</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Created At</th>
                  </tr>
                </thead>
                <tbody>
                  {data.map((r, i) => {
                    const s = statusMeta[r.status] || { label: r.status, bg: '#F1F5F9', text: '#475569', dot: '#475569' }
                    return (
                      <tr key={`${r.orderId}-${r.paymentId}-${i}`} className="border-b border-slate-100 transition-colors last:border-0 hover:bg-[#F9FAFB]">
                        <td className="whitespace-nowrap px-5 py-3.5 FleetOps-mono text-[13px] font-medium text-[#0B0E14]">#{r.order_id}</td>
                        <td className="whitespace-nowrap px-5 py-3.5 FleetOps-mono text-[13px] text-[#5B6472]">#{r.payment_id}</td>
                        <td className="whitespace-nowrap px-5 py-3.5 FleetOps-mono text-[13px] text-[#5B6472]">#{r.user_id}</td>
                        <td className="whitespace-nowrap px-5 py-3.5 FleetOps-body text-[13px] text-[#0B0E14]">{r.email}</td>
                        <td className="whitespace-nowrap px-5 py-3.5">
                          <span
                            className="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 FleetOps-mono text-[10px] font-medium"
                            style={{ backgroundColor: s.bg, color: s.text }}
                          >
                            <span className="h-1.5 w-1.5 rounded-full" style={{ backgroundColor: s.dot }} />
                            {s.label}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-5 py-3.5 FleetOps-mono text-[11.5px] text-[#5B6472]">
                          {formatTimestamp(r.created_at)}
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          </div>
        ) : (
          <div className="mt-4 flex flex-col items-center gap-3 rounded-2xl border border-dashed border-slate-300 bg-white py-16 text-center">
            {NodeIcon('inbox', 'h-8 w-8 text-[#5B6472]')}
            <p className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">No events match this filter</p>
          </div>
        )}
      </main>
    </div>
  )
}