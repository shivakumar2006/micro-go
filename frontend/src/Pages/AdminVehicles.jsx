import { useMemo, useState } from 'react'

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
  search: (
    <>
      <circle cx="11" cy="11" r="7" />
      <path d="M21 21l-4.35-4.35" />
    </>
  ),
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
  chevronLeft: <path d="M15 6l-6 6 6 6" />,
  chevronRight: <path d="M9 6l6 6-6 6" />,
  plus: (
    <>
      <path d="M12 5v14" />
      <path d="M5 12h14" />
    </>
  ),
  edit: (
    <>
      <path d="M4 20.5l.9-3.6 11-11a2 2 0 0 1 2.8 0l1.8 1.8a2 2 0 0 1 0 2.8l-11 11-3.6.9z" />
      <path d="M14.5 7l2.5 2.5" />
    </>
  ),
  trash: (
    <>
      <path d="M4 7h16" />
      <path d="M9 7V5a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
      <path d="M6 7l1 13a1.5 1.5 0 0 0 1.5 1.4h7a1.5 1.5 0 0 0 1.5-1.4L18 7" />
      <path d="M10 11v6" />
      <path d="M14 11v6" />
    </>
  ),
  close: <path d="M6 6l12 12M18 6L6 18" />,
  bike: (
    <>
      <circle cx="5.5" cy="17" r="3" />
      <circle cx="18.5" cy="17" r="3" />
      <path d="M5.5 17l4-9h4.5l3 6" />
      <path d="M9.5 8h3.5" />
      <path d="M13 12h5.5" />
    </>
  ),
  car: (
    <>
      <path d="M4 16.5V13l1.5-4.2A2 2 0 0 1 7.4 7.5h9.2a2 2 0 0 1 1.9 1.3L20 13v3.5" />
      <rect x="3" y="13" width="18" height="5" rx="1.5" />
      <circle cx="7.5" cy="18.5" r="1.6" />
      <circle cx="16.5" cy="18.5" r="1.6" />
    </>
  ),
  van: (
    <>
      <rect x="2" y="9" width="11" height="7" rx="1" />
      <path d="M13 11.5h3.5L19 14v2h-6z" />
      <circle cx="6.2" cy="18" r="1.5" />
      <circle cx="16.5" cy="18" r="1.5" />
    </>
  ),
  inbox: (
    <>
      <path d="M3.5 12h5l1.5 3h4l1.5-3h5" />
      <path d="M5 6.5h14L21 12v6a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18v-6l2-5.5z" />
    </>
  ),
}

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

// ---------------------------------------------------------------------------
// Mock data — same shape the API will return
// ---------------------------------------------------------------------------

const statusMeta = {
  available: { label: 'Available', dot: '#16A34A', bg: '#16A34A14', text: '#15803D' },
  'on-trip': { label: 'On trip', dot: '#FF5A1F', bg: '#FF5A1F14', text: '#C2410C' },
  maintenance: { label: 'Maintenance', dot: '#D97706', bg: '#D9770614', text: '#B45309' },
  offline: { label: 'Offline', dot: '#94A3B8', bg: '#94A3B814', text: '#5B6472' },
}

const initialVehicles = [
  ['FL-104', 'van', 'available', 'Rohit Sharma', 'Sector 14, Gurugram', 82, 2],
  ['FL-089', 'bike', 'on-trip', 'Aman Verma', 'Cyber Hub, Gurugram', 46, 5],
  ['FL-231', 'car', 'maintenance', 'Unassigned', 'Service Center, Udyog Vihar', 12, 210],
  ['FL-057', 'van', 'on-trip', 'Priya Nair', 'Sector 62, Noida', 63, 11],
  ['FL-142', 'bike', 'available', 'Karan Mehta', 'MG Road, Bengaluru', 91, 1],
  ['FL-018', 'car', 'available', 'Sana Sheikh', 'Andheri West, Mumbai', 74, 8],
  ['FL-276', 'van', 'offline', 'Unassigned', 'Depot, Manesar', 0, 960],
  ['FL-063', 'bike', 'on-trip', 'Vikram Rao', 'Koramangala, Bengaluru', 38, 3],
  ['FL-199', 'car', 'available', 'Neha Gupta', 'Powai, Mumbai', 88, 6],
  ['FL-021', 'van', 'maintenance', 'Unassigned', 'Service Center, Udyog Vihar', 24, 300],
  ['FL-158', 'bike', 'available', 'Farhan Ali', 'Golf Course Road, Gurugram', 95, 4],
  ['FL-092', 'car', 'on-trip', 'Ritu Singh', 'HSR Layout, Bengaluru', 52, 14],
].map(([plate, type, status, driver, location, battery, lastActiveMinutes], i) => ({
  id: i + 1,
  plate,
  type,
  status,
  driver,
  location,
  battery,
  lastActiveMinutes,
}))

const emptyForm = { plate: '', type: 'bike', status: 'available', driver: '', location: '', battery: 100 }

function formatLastActive(minutes) {
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes} min ago`
  if (minutes < 1440) return `${Math.round(minutes / 60)} hr ago`
  return `${Math.round(minutes / 1440)} d ago`
}

function getInitials(name) {
  if (!name || name === 'Unassigned') return '—'
  return name
    .split(' ')
    .filter(Boolean)
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()
}

// ---------------------------------------------------------------------------
// Shared building blocks
// ---------------------------------------------------------------------------

function Avatar({ initials }) {
  return (
    <span className="flex h-6 w-6 flex-none items-center justify-center rounded-full bg-[#0B0E14] FleetOps-mono text-[10px] font-semibold text-white">
      {initials}
    </span>
  )
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

function TextField({ label, type = 'text', value, onChange, placeholder, min, max }) {
  return (
    <label className="block">
      <span className="FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">{label}</span>
      <input
        type={type}
        value={value}
        min={min}
        max={max}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="mt-1.5 w-full rounded-xl border border-slate-200 bg-white px-3.5 py-2.5 FleetOps-body text-[14px] text-[#0B0E14] placeholder:text-slate-400 transition-colors focus:border-[#0B0E14] focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/10"
      />
    </label>
  )
}

function SelectInput({ label, value, onChange, options }) {
  return (
    <label className="block">
      <span className="FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">{label}</span>
      <div className="relative mt-1.5">
        <select
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="w-full appearance-none rounded-xl border border-slate-200 bg-white px-3.5 py-2.5 pr-8 FleetOps-body text-[14px] text-[#0B0E14] transition-colors focus:border-[#0B0E14] focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/10"
        >
          {options.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
        <span className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-[#5B6472]">
          {NodeIcon('chevronDown', 'h-3.5 w-3.5')}
        </span>
      </div>
    </label>
  )
}

// ---------------------------------------------------------------------------
// Modals
// ---------------------------------------------------------------------------

function VehicleFormModal({ mode, initial, onCancel, onSave }) {
  const [form, setForm] = useState(initial)
  const set = (key) => (val) => setForm((f) => ({ ...f, [key]: val }))
  const canSave = form.plate.trim().length > 0

  const handleSubmit = (e) => {
    e.preventDefault()
    if (!canSave) return
    onSave({
      ...form,
      plate: form.plate.trim(),
      driver: form.driver.trim() || 'Unassigned',
      location: form.location.trim(),
      battery: Math.min(100, Math.max(0, Number(form.battery) || 0)),
    })
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-[#0B0E14]/50 px-4 backdrop-blur-sm"
      onClick={onCancel}
    >
      <div
        className="w-full max-w-md rounded-[28px] border border-slate-200 bg-white p-8 shadow-[0_30px_60px_-30px_rgba(15,23,42,0.35)]"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-start justify-between">
          <div>
            <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">
              {mode === 'create' ? 'add vehicle' : 'edit vehicle'}
            </span>
            <h2 className="mt-1 FleetOps-display text-xl font-semibold text-[#0B0E14]">
              {mode === 'create' ? 'Add a new vehicle' : `Edit ${initial.plate}`}
            </h2>
          </div>
          <button
            type="button"
            onClick={onCancel}
            className="flex h-8 w-8 flex-none items-center justify-center rounded-full text-[#5B6472] transition-colors hover:bg-[#F3F5F8] hover:text-[#0B0E14]"
            aria-label="Close"
          >
            {NodeIcon('close', 'h-4 w-4')}
          </button>
        </div>

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <TextField label="Plate / ID" value={form.plate} onChange={set('plate')} placeholder="FL-104" />

          <div className="grid grid-cols-2 gap-4">
            <SelectInput
              label="Type"
              value={form.type}
              onChange={set('type')}
              options={[
                { value: 'bike', label: 'Bike' },
                { value: 'car', label: 'Car' },
                { value: 'van', label: 'Van' },
              ]}
            />
            <SelectInput
              label="Status"
              value={form.status}
              onChange={set('status')}
              options={[
                { value: 'available', label: 'Available' },
                { value: 'on-trip', label: 'On trip' },
                { value: 'maintenance', label: 'Maintenance' },
                { value: 'offline', label: 'Offline' },
              ]}
            />
          </div>

          <TextField label="Driver" value={form.driver} onChange={set('driver')} placeholder="Unassigned" />
          <TextField label="Location" value={form.location} onChange={set('location')} placeholder="Sector 14, Gurugram" />
          <TextField label="Battery %" type="number" min="0" max="100" value={form.battery} onChange={set('battery')} placeholder="0–100" />

          <div className="flex items-center justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={onCancel}
              className="rounded-full border border-slate-200 px-5 py-2.5 FleetOps-body text-[13.5px] font-medium text-[#0B0E14] transition-colors hover:bg-[#F3F5F8]"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={!canSave}
              className="rounded-full bg-[#FF5A1F] px-5 py-2.5 FleetOps-body text-[13.5px] font-semibold text-white shadow-[0_12px_24px_-8px_rgba(255,90,31,0.55)] transition-all hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:translate-y-0"
            >
              {mode === 'create' ? 'Add vehicle' : 'Save changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

function DeleteConfirm({ vehicle, onCancel, onConfirm }) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-[#0B0E14]/50 px-4 backdrop-blur-sm"
      onClick={onCancel}
    >
      <div
        className="w-full max-w-sm rounded-[28px] border border-slate-200 bg-white p-8 text-center shadow-[0_30px_60px_-30px_rgba(15,23,42,0.35)]"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-red-50 text-[#DC2626]">
          {NodeIcon('trash', 'h-5 w-5')}
        </div>
        <h2 className="mt-4 FleetOps-display text-lg font-semibold text-[#0B0E14]">Delete {vehicle.plate}?</h2>
        <p className="mt-2 FleetOps-body text-[13.5px] text-[#5B6472]">
          This can&rsquo;t be undone. The vehicle will be removed from the fleet.
        </p>
        <div className="mt-6 flex items-center justify-center gap-3">
          <button
            type="button"
            onClick={onCancel}
            className="flex-1 rounded-full border border-slate-200 px-5 py-2.5 FleetOps-body text-[13.5px] font-medium text-[#0B0E14] transition-colors hover:bg-[#F3F5F8]"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className="flex-1 rounded-full bg-[#DC2626] px-5 py-2.5 FleetOps-body text-[13.5px] font-semibold text-white transition-colors hover:bg-[#B91C1C]"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------------------

function Pagination({ page, totalPages, onChange }) {
  if (totalPages <= 1) return null
  const pages = Array.from({ length: totalPages }, (_, i) => i + 1)
  const pagerBtn = 'flex h-9 w-9 items-center justify-center rounded-full FleetOps-mono text-[13px] transition-colors'

  return (
    <div className="mt-8 flex items-center justify-center gap-1.5">
      <button
        type="button"
        onClick={() => onChange(Math.max(1, page - 1))}
        disabled={page === 1}
        className={`${pagerBtn} text-[#5B6472] hover:bg-[#F3F5F8] disabled:cursor-not-allowed disabled:opacity-30 disabled:hover:bg-transparent`}
        aria-label="Previous page"
      >
        {NodeIcon('chevronLeft', 'h-4 w-4')}
      </button>
      {pages.map((p) => (
        <button
          key={p}
          type="button"
          onClick={() => onChange(p)}
          className={p === page ? `${pagerBtn} bg-[#0B0E14] text-white` : `${pagerBtn} text-[#5B6472] hover:bg-[#F3F5F8]`}
        >
          {p}
        </button>
      ))}
      <button
        type="button"
        onClick={() => onChange(Math.min(totalPages, page + 1))}
        disabled={page === totalPages}
        className={`${pagerBtn} text-[#5B6472] hover:bg-[#F3F5F8] disabled:cursor-not-allowed disabled:opacity-30 disabled:hover:bg-transparent`}
        aria-label="Next page"
      >
        {NodeIcon('chevronRight', 'h-4 w-4')}
      </button>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Page
// ---------------------------------------------------------------------------

const PER_PAGE = 10

export default function AdminVehicles() {
  const [vehicles, setVehicles] = useState(initialVehicles)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')
  const [typeFilter, setTypeFilter] = useState('all')
  const [page, setPage] = useState(1)

  const [formState, setFormState] = useState(null) // { mode: 'create'|'edit', initial, id? }
  const [deleteTarget, setDeleteTarget] = useState(null)

  const stats = useMemo(
    () => ({
      total: vehicles.length,
      available: vehicles.filter((v) => v.status === 'available').length,
      onTrip: vehicles.filter((v) => v.status === 'on-trip').length,
      needsAttention: vehicles.filter((v) => v.status === 'maintenance' || v.status === 'offline').length,
    }),
    [vehicles]
  )

  const filtered = useMemo(() => {
    let list = vehicles
    if (search.trim()) {
      const q = search.trim().toLowerCase()
      list = list.filter((v) => v.plate.toLowerCase().includes(q) || v.driver.toLowerCase().includes(q))
    }
    if (statusFilter !== 'all') list = list.filter((v) => v.status === statusFilter)
    if (typeFilter !== 'all') list = list.filter((v) => v.type === typeFilter)
    return list
  }, [vehicles, search, statusFilter, typeFilter])

  const totalPages = Math.max(1, Math.ceil(filtered.length / PER_PAGE))
  const currentPage = Math.min(page, totalPages)
  const paginated = filtered.slice((currentPage - 1) * PER_PAGE, currentPage * PER_PAGE)

  const withPageReset = (setter) => (val) => {
    setter(val)
    setPage(1)
  }

  // --- CRUD handlers -------------------------------------------------------
  // Swap the setVehicles(...) lines below for the matching API call once a
  // backend is wired up; the rest of the page doesn't need to change.

  function handleCreate(data) {
    const nextId = vehicles.reduce((max, v) => Math.max(max, v.id), 0) + 1
    setVehicles((vs) => [{ id: nextId, ...data, lastActiveMinutes: 0 }, ...vs]) // 1. POST /vehicles
    setFormState(null)
  }

  function handleUpdate(id, data) {
    setVehicles((vs) => vs.map((v) => (v.id === id ? { ...v, ...data } : v))) // 2. PUT /vehicles/:id
    setFormState(null)
  }

  function handleDelete(id) {
    setVehicles((vs) => vs.filter((v) => v.id !== id)) // 3. DELETE /vehicles/:id
    setDeleteTarget(null)
  }

  return (
    <div className="min-h-screen bg-[#F3F5F8] FleetOps-body text-[#0B0E14]">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@500;600;700&family=Manrope:wght@400;500;600;700&family=IBM+Plex+Mono:wght@400;500&display=swap');
        .FleetOps-display { font-family: 'Space Grotesk', ui-sans-serif, sans-serif; }
        .FleetOps-body { font-family: 'Manrope', ui-sans-serif, sans-serif; }
        .FleetOps-mono { font-family: 'IBM Plex Mono', ui-monospace, monospace; }
      `}</style>

      <main className="mx-auto max-w-7xl px-6 py-10">
        <div className="flex flex-col gap-1 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">admin</span>
            <h1 className="mt-1 FleetOps-display text-2xl font-semibold tracking-tight text-[#0B0E14] sm:text-3xl">
              Manage vehicles
            </h1>
          </div>
          <button
            type="button"
            onClick={() => setFormState({ mode: 'create', initial: emptyForm })}
            className="flex items-center gap-1.5 rounded-full bg-[#FF5A1F] px-5 py-2.5 FleetOps-body text-[13.5px] font-semibold text-white shadow-[0_12px_24px_-8px_rgba(255,90,31,0.55)] transition-all hover:-translate-y-0.5"
          >
            {NodeIcon('plus', 'h-4 w-4')}
            Add vehicle
          </button>
        </div>

        <div className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-4">
          <StatCard label="Total fleet" value={stats.total} />
          <StatCard label="Available" value={stats.available} accent="#15803D" />
          <StatCard label="On trip" value={stats.onTrip} accent="#C2410C" />
          <StatCard label="Needs attention" value={stats.needsAttention} accent="#B45309" />
        </div>

        <div className="mt-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <label className="relative block w-full sm:max-w-xs">
            <span className="sr-only">Search vehicles</span>
            <span className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-[#5B6472]">
              {NodeIcon('search', 'h-4 w-4')}
            </span>
            <input
              type="text"
              value={search}
              onChange={(e) => withPageReset(setSearch)(e.target.value)}
              placeholder="Search by plate or driver"
              className="w-full rounded-full border border-slate-200 bg-white py-2.5 pl-10 pr-4 FleetOps-body text-[13.5px] text-[#0B0E14] placeholder:text-slate-400 transition-colors focus:border-[#0B0E14] focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/10"
            />
          </label>

          <div className="flex flex-wrap items-center gap-3">
            <SelectField
              icon="filter"
              label="Filter by status"
              value={statusFilter}
              onChange={withPageReset(setStatusFilter)}
              options={[
                { value: 'all', label: 'All statuses' },
                { value: 'available', label: 'Available' },
                { value: 'on-trip', label: 'On trip' },
                { value: 'maintenance', label: 'Maintenance' },
                { value: 'offline', label: 'Offline' },
              ]}
            />
            <SelectField
              icon="van"
              label="Filter by type"
              value={typeFilter}
              onChange={withPageReset(setTypeFilter)}
              options={[
                { value: 'all', label: 'All types' },
                { value: 'bike', label: 'Bike' },
                { value: 'car', label: 'Car' },
                { value: 'van', label: 'Van' },
              ]}
            />
          </div>
        </div>

        {paginated.length > 0 ? (
          <div className="mt-6 overflow-hidden rounded-2xl border border-slate-200 bg-white">
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead>
                  <tr className="border-b border-slate-100 bg-[#F9FAFB]">
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Vehicle</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Status</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Driver</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Location</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Battery</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Last active</th>
                    <th className="px-5 py-3" />
                  </tr>
                </thead>
                <tbody>
                  {paginated.map((v) => {
                    const s = statusMeta[v.status]
                    return (
                      <tr key={v.id} className="border-b border-slate-100 transition-colors last:border-0 hover:bg-[#F9FAFB]">
                        <td className="whitespace-nowrap px-5 py-4">
                          <div className="flex items-center gap-3">
                            <div className="flex h-9 w-9 flex-none items-center justify-center rounded-lg bg-[#0B0E14] text-white">
                              {NodeIcon(v.type, 'h-4 w-4')}
                            </div>
                            <div>
                              <p className="FleetOps-display text-[14px] font-semibold text-[#0B0E14]">{v.plate}</p>
                              <p className="FleetOps-mono text-[10px] uppercase tracking-wide text-[#5B6472]">{v.type}</p>
                            </div>
                          </div>
                        </td>
                        <td className="whitespace-nowrap px-5 py-4">
                          <span
                            className="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 FleetOps-mono text-[10px] font-medium"
                            style={{ backgroundColor: s.bg, color: s.text }}
                          >
                            <span className="h-1.5 w-1.5 rounded-full" style={{ backgroundColor: s.dot }} />
                            {s.label}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-5 py-4">
                          <div className="flex items-center gap-2 FleetOps-body text-[13px] text-[#5B6472]">
                            <Avatar initials={getInitials(v.driver)} />
                            {v.driver}
                          </div>
                        </td>
                        <td className="whitespace-nowrap px-5 py-4 FleetOps-body text-[13px] text-[#5B6472]">{v.location}</td>
                        <td className="whitespace-nowrap px-5 py-4">
                          <div className="flex items-center gap-2">
                            <div className="h-1.5 w-14 overflow-hidden rounded-full bg-slate-100">
                              <div
                                className="h-full rounded-full"
                                style={{ width: `${v.battery}%`, backgroundColor: v.battery > 30 ? '#0B0E14' : '#D97706' }}
                              />
                            </div>
                            <span className="FleetOps-mono text-[11px] text-[#5B6472]">{v.battery}%</span>
                          </div>
                        </td>
                        <td className="whitespace-nowrap px-5 py-4 FleetOps-mono text-[11px] text-[#5B6472]">
                          {formatLastActive(v.lastActiveMinutes)}
                        </td>
                        <td className="whitespace-nowrap px-5 py-4">
                          <div className="flex items-center justify-end gap-1.5">
                            <button
                              type="button"
                              onClick={() => setFormState({ mode: 'edit', initial: v, id: v.id })}
                              className="flex h-8 w-8 items-center justify-center rounded-lg text-[#5B6472] transition-colors hover:bg-[#F3F5F8] hover:text-[#0B0E14]"
                              aria-label={`Edit ${v.plate}`}
                            >
                              {NodeIcon('edit', 'h-4 w-4')}
                            </button>
                            <button
                              type="button"
                              onClick={() => setDeleteTarget(v)}
                              className="flex h-8 w-8 items-center justify-center rounded-lg text-[#5B6472] transition-colors hover:bg-red-50 hover:text-[#DC2626]"
                              aria-label={`Delete ${v.plate}`}
                            >
                              {NodeIcon('trash', 'h-4 w-4')}
                            </button>
                          </div>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          </div>
        ) : (
          <div className="mt-6 flex flex-col items-center gap-3 rounded-2xl border border-dashed border-slate-300 bg-white py-16 text-center">
            {NodeIcon('inbox', 'h-8 w-8 text-[#5B6472]')}
            <p className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">No vehicles match your filters</p>
            <p className="FleetOps-body text-[13.5px] text-[#5B6472]">Try a different search term, reset the filters, or add a new vehicle.</p>
            <button
              type="button"
              onClick={() => {
                setSearch('')
                setStatusFilter('all')
                setTypeFilter('all')
                setPage(1)
              }}
              className="mt-1 rounded-full bg-[#0B0E14] px-4 py-2 FleetOps-body text-[13px] font-medium text-white transition-colors hover:bg-[#1a2030]"
            >
              Reset filters
            </button>
          </div>
        )}

        <Pagination page={currentPage} totalPages={totalPages} onChange={setPage} />
      </main>

      {formState && (
        <VehicleFormModal
          mode={formState.mode}
          initial={formState.initial}
          onCancel={() => setFormState(null)}
          onSave={(data) => (formState.mode === 'create' ? handleCreate(data) : handleUpdate(formState.id, data))}
        />
      )}

      {deleteTarget && (
        <DeleteConfirm
          vehicle={deleteTarget}
          onCancel={() => setDeleteTarget(null)}
          onConfirm={() => handleDelete(deleteTarget.id)}
        />
      )}
    </div>
  )
}