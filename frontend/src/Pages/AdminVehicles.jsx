import { useEffect, useMemo, useState } from 'react'
import {
  useGetAllVehiclesQuery,
  useGetVehicleByIdQuery,
  useCreateVehicleMutation,
  useUpdateVehicleMutation,
  useDeleteVehicleMutation,
} from '../Redux/features/vehicles/vehicleApi'

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
  sort: (
    <>
      <path d="M6 7h8" />
      <path d="M6 12h5" />
      <path d="M6 17h2" />
      <path d="M17 5v14" />
      <path d="M14 15l3 3 3-3" />
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
  eye: (
    <>
      <path d="M2.5 12S5.7 5.5 12 5.5 21.5 12 21.5 12 18.3 18.5 12 18.5 2.5 12 2.5 12z" />
      <circle cx="12" cy="12" r="2.6" />
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
  suv: (
    <>
      <path d="M3.5 16V11l2-4.5h9l3 4.5v5" />
      <rect x="2.5" y="16" width="19" height="4" rx="1.3" />
      <circle cx="7" cy="20" r="1.6" />
      <circle cx="17" cy="20" r="1.6" />
    </>
  ),
  image: (
    <>
      <rect x="3" y="4.5" width="18" height="15" rx="1.5" />
      <circle cx="8.5" cy="9.5" r="1.5" />
      <path d="M3 16l5-5 4 4 3-3 6 6" />
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

const categoryMeta = {
  Normal: { bg: '#DBEAFE', text: '#1D4ED8', dot: '#1D4ED8' },
  Moderate: { bg: '#FEF3C7', text: '#B45309', dot: '#B45309' },
  Premium: { bg: '#F3E8FF', text: '#7E22CE', dot: '#7E22CE' },
}

const emptyForm = {
  name: '',
  brand: '',
  model: '',
  type: '',
  category: 'Normal',
  price: '',
  stock: '',
  description: '',
  image_url: '',
}

const PER_PAGE = 10

function formatPrice(amount) {
  const n = Number(amount)
  if (Number.isNaN(n)) return '—'
  return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(n)
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('en-IN', { day: '2-digit', month: 'short', year: 'numeric' })
}

function typeIconFor(type) {
  const key = type?.toLowerCase()
  return icons[key] ? key : 'car'
}

// ---------------------------------------------------------------------------
// Shared building blocks
// ---------------------------------------------------------------------------

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

function TextField({ label, type = 'text', value, onChange, placeholder, min }) {
  return (
    <label className="block">
      <span className="FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">{label}</span>
      <input
        type={type}
        value={value}
        min={min}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="mt-1.5 w-full rounded-xl border border-slate-200 bg-white px-3.5 py-2.5 FleetOps-body text-[14px] text-[#0B0E14] placeholder:text-slate-400 transition-colors focus:border-[#0B0E14] focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/10"
      />
    </label>
  )
}

function TextAreaField({ label, value, onChange, placeholder }) {
  return (
    <label className="block">
      <span className="FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">{label}</span>
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        rows={3}
        className="mt-1.5 w-full resize-none rounded-xl border border-slate-200 bg-white px-3.5 py-2.5 FleetOps-body text-[14px] text-[#0B0E14] placeholder:text-slate-400 transition-colors focus:border-[#0B0E14] focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/10"
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

function Thumb({ src, alt }) {
  const [failed, setFailed] = useState(false)
  if (!src || failed) {
    return (
      <div className="flex h-11 w-11 flex-none items-center justify-center rounded-lg bg-[#F3F5F8] text-[#5B6472]">
        {NodeIcon('image', 'h-4 w-4')}
      </div>
    )
  }
  return (
    <img
      src={src}
      alt={alt}
      onError={() => setFailed(true)}
      className="h-11 w-11 flex-none rounded-lg border border-slate-100 object-cover"
    />
  )
}

// ---------------------------------------------------------------------------
// Modals
// ---------------------------------------------------------------------------

function VehicleFormModal({ mode, initial, onCancel, onSave, isSaving, errorMessage }) {
  const [form, setForm] = useState(initial)
  const set = (key) => (val) => setForm((f) => ({ ...f, [key]: val }))
  const canSave = form.name.trim().length > 0 && form.brand.trim().length > 0 && !isSaving

  const handleSubmit = (e) => {
    e.preventDefault()
    if (!canSave) return
    onSave({
      ...form,
      name: form.name.trim(),
      brand: form.brand.trim(),
      model: form.model.trim(),
      type: form.type.trim(),
      description: form.description.trim(),
      image_url: form.image_url.trim(),
      price: Number(form.price) || 0,
      stock: Number(form.stock) || 0,
    })
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-[#0B0E14]/50 px-4 backdrop-blur-sm"
      onClick={onCancel}
    >
      <div
        className="max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-[28px] border border-slate-200 bg-white p-8 shadow-[0_30px_60px_-30px_rgba(15,23,42,0.35)]"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-start justify-between">
          <div>
            <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">
              {mode === 'create' ? 'add vehicle' : 'edit vehicle'}
            </span>
            <h2 className="mt-1 FleetOps-display text-xl font-semibold text-[#0B0E14]">
              {mode === 'create' ? 'Add a new vehicle' : `Edit ${initial.name}`}
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

        {errorMessage && (
          <p className="mt-4 rounded-xl bg-red-50 px-3.5 py-2.5 FleetOps-body text-[13px] text-[#DC2626]">{errorMessage}</p>
        )}

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <TextField label="Name" value={form.name} onChange={set('name')} placeholder="Fortuner" />

          <div className="grid grid-cols-2 gap-4">
            <TextField label="Brand" value={form.brand} onChange={set('brand')} placeholder="Toyota" />
            <TextField label="Model" value={form.model} onChange={set('model')} placeholder="GR Sport" />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <TextField label="Type" value={form.type} onChange={set('type')} placeholder="SUV" />
            <SelectInput
              label="Category"
              value={form.category}
              onChange={set('category')}
              options={[
                { value: 'Normal', label: 'Normal' },
                { value: 'Moderate', label: 'Moderate' },
                { value: 'Premium', label: 'Premium' },
              ]}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <TextField label="Price (₹)" type="number" min="0" value={form.price} onChange={set('price')} placeholder="4899999" />
            <TextField label="Stock" type="number" min="0" value={form.stock} onChange={set('stock')} placeholder="15" />
          </div>

          <TextField label="Image URL" value={form.image_url} onChange={set('image_url')} placeholder="https://…" />
          <TextAreaField
            label="Description"
            value={form.description}
            onChange={set('description')}
            placeholder="Toyota Fortuner GR Sport 4x4 Automatic Diesel"
          />

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
              {isSaving ? 'Saving…' : mode === 'create' ? 'Add vehicle' : 'Save changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

function VehicleViewModal({ id, onClose }) {
  const { data: vehicle, isLoading, error } = useGetVehicleByIdQuery(id, { skip: !id })

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-[#0B0E14]/50 px-4 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        className="max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-[28px] border border-slate-200 bg-white p-8 shadow-[0_30px_60px_-30px_rgba(15,23,42,0.35)]"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-start justify-between">
          <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">vehicle details</span>
          <button
            type="button"
            onClick={onClose}
            className="flex h-8 w-8 flex-none items-center justify-center rounded-full text-[#5B6472] transition-colors hover:bg-[#F3F5F8] hover:text-[#0B0E14]"
            aria-label="Close"
          >
            {NodeIcon('close', 'h-4 w-4')}
          </button>
        </div>

        {isLoading && <p className="mt-6 FleetOps-body text-[13.5px] text-[#5B6472]">Loading…</p>}
        {error && <p className="mt-6 FleetOps-body text-[13.5px] text-[#DC2626]">Couldn&rsquo;t load this vehicle.</p>}

        {vehicle && (
          <>
            <div className="mt-5 overflow-hidden rounded-2xl bg-[#F3F5F8]">
              {vehicle.image_url ? (
                <img src={vehicle.image_url} alt={vehicle.name} className="aspect-[4/3] w-full object-cover" />
              ) : (
                <div className="flex aspect-[4/3] w-full items-center justify-center text-[#5B6472]">
                  {NodeIcon('image', 'h-8 w-8')}
                </div>
              )}
            </div>

            <h2 className="mt-5 FleetOps-display text-xl font-semibold text-[#0B0E14]">
              {vehicle.brand} {vehicle.name}
            </h2>
            <p className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">{vehicle.model}</p>
            <p className="mt-3 FleetOps-display text-xl font-semibold text-[#0B0E14]">{formatPrice(vehicle.price)}</p>

            {vehicle.description && (
              <p className="mt-3 FleetOps-body text-[13.5px] leading-relaxed text-[#5B6472]">{vehicle.description}</p>
            )}

            <div className="mt-5 rounded-2xl border border-slate-200 px-4">
              {[
                ['Type', vehicle.type],
                ['Category', vehicle.category],
                ['Stock', `${vehicle.stock} units`],
                ['Added', formatDate(vehicle.created_at)],
              ].map(([label, value]) => (
                <div key={label} className="flex items-center justify-between border-b border-slate-100 py-3 last:border-0">
                  <span className="FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">{label}</span>
                  <span className="FleetOps-body text-[13.5px] font-medium text-[#0B0E14]">{value}</span>
                </div>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  )
}

function DeleteConfirm({ vehicle, onCancel, onConfirm, isDeleting }) {
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
        <h2 className="mt-4 FleetOps-display text-lg font-semibold text-[#0B0E14]">Delete {vehicle.name}?</h2>
        <p className="mt-2 FleetOps-body text-[13.5px] text-[#5B6472]">
          This can&rsquo;t be undone. The vehicle will be removed from the catalog.
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
            disabled={isDeleting}
            className="flex-1 rounded-full bg-[#DC2626] px-5 py-2.5 FleetOps-body text-[13.5px] font-semibold text-white transition-colors hover:bg-[#B91C1C] disabled:cursor-not-allowed disabled:opacity-50"
          >
            {isDeleting ? 'Deleting…' : 'Delete'}
          </button>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Pagination — falls back to Prev/Next-only if the API doesn't return a total
// ---------------------------------------------------------------------------

function Pagination({ page, totalPages, hasNext, onChange }) {
  const pagerBtn = 'flex h-9 w-9 items-center justify-center rounded-full FleetOps-mono text-[13px] transition-colors'

  if (totalPages) {
    if (totalPages <= 1) return null
    const pages = Array.from({ length: totalPages }, (_, i) => i + 1)
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

  // no total from the API — plain prev/next
  if (page === 1 && !hasNext) return null
  return (
    <div className="mt-8 flex items-center justify-center gap-3">
      <button
        type="button"
        onClick={() => onChange(Math.max(1, page - 1))}
        disabled={page === 1}
        className="flex items-center gap-1.5 rounded-full border border-slate-200 px-4 py-2 FleetOps-body text-[13px] font-medium text-[#0B0E14] transition-colors hover:bg-[#F3F5F8] disabled:cursor-not-allowed disabled:opacity-40"
      >
        {NodeIcon('chevronLeft', 'h-3.5 w-3.5')}
        Prev
      </button>
      <span className={`${pagerBtn} bg-[#0B0E14] text-white`}>{page}</span>
      <button
        type="button"
        onClick={() => onChange(page + 1)}
        disabled={!hasNext}
        className="flex items-center gap-1.5 rounded-full border border-slate-200 px-4 py-2 FleetOps-body text-[13px] font-medium text-[#0B0E14] transition-colors hover:bg-[#F3F5F8] disabled:cursor-not-allowed disabled:opacity-40"
      >
        Next
        {NodeIcon('chevronRight', 'h-3.5 w-3.5')}
      </button>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Page
// ---------------------------------------------------------------------------

export default function AdminVehicles() {
  const [searchInput, setSearchInput] = useState('')
  const [search, setSearch] = useState('')
  const [typeFilter, setTypeFilter] = useState('all')
  const [categoryFilter, setCategoryFilter] = useState('all')
  const [sortBy, setSortBy] = useState('created_at')
  const [order, setOrder] = useState('desc')
  const [page, setPage] = useState(1)

  const [formState, setFormState] = useState(null) // { mode: 'create'|'edit', initial, id? }
  const [viewId, setViewId] = useState(null)
  const [deleteTarget, setDeleteTarget] = useState(null)
  const [formError, setFormError] = useState('')

  // debounce search so every keystroke doesn't fire a request
  useEffect(() => {
    const t = setTimeout(() => {
      setSearch(searchInput)
      setPage(1)
    }, 400)
    return () => clearTimeout(t)
  }, [searchInput])

  const { data, isLoading, isFetching, error } = useGetAllVehiclesQuery({
    page,
    limit: PER_PAGE,
    search,
    type: typeFilter === 'all' ? undefined : typeFilter,
    category: categoryFilter === 'all' ? undefined : categoryFilter,
    sortBy,
    order,
  })

  console.log("vehicles data : ", data);

  const [createVehicle, { isLoading: isCreating }] = useCreateVehicleMutation()
  const [updateVehicle, { isLoading: isUpdating }] = useUpdateVehicleMutation()
  const [deleteVehicle, { isLoading: isDeleting }] = useDeleteVehicleMutation()

  // --- unwrap the response --------------------------------------------------
  // supports either { data: [...], total } or a bare array — see the note
  // at the top of the file if your API shape differs.
  const vehicles = data?.data ?? (Array.isArray(data) ? data : []) ?? []
  const total = data?.total ?? data?.count ?? null
  const totalPages = total ? Math.max(1, Math.ceil(total / PER_PAGE)) : null
  const hasNext = vehicles.length === PER_PAGE

  const stats = useMemo(
    () => ({
      shown: total ?? vehicles.length,
      premium: vehicles.filter((v) => v.category === 'Premium').length,
      totalStock: vehicles.reduce((sum, v) => sum + (Number(v.stock) || 0), 0),
    }),
    [vehicles, total]
  )

  const withPageReset = (setter) => (val) => {
    setter(val)
    setPage(1)
  }

  // --- CRUD handlers ---------------------------------------------------------

  async function handleCreate(payload) {
    setFormError('')
    try {
      await createVehicle(payload).unwrap()
      setFormState(null)
    } catch (err) {
      setFormError(err?.data?.message || 'Could not create this vehicle. Please try again.')
    }
  }

  async function handleUpdate(id, payload) {
    setFormError('')
    try {
      await updateVehicle({ id, ...payload }).unwrap()
      setFormState(null)
    } catch (err) {
      setFormError(err?.data?.message || 'Could not save changes. Please try again.')
    }
  }

  async function handleDelete(id) {
    try {
      await deleteVehicle(id).unwrap()
      setDeleteTarget(null)
    } catch (err) {
      setDeleteTarget(null)
    }
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
            onClick={() => {
              setFormError('')
              setFormState({ mode: 'create', initial: emptyForm })
            }}
            className="flex items-center gap-1.5 rounded-full bg-[#FF5A1F] px-5 py-2.5 FleetOps-body text-[13.5px] font-semibold text-white shadow-[0_12px_24px_-8px_rgba(255,90,31,0.55)] transition-all hover:-translate-y-0.5"
          >
            {NodeIcon('plus', 'h-4 w-4')}
            Add vehicle
          </button>
        </div>

        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
          <StatCard label={total ? 'Total vehicles' : 'Vehicles (this page)'} value={stats.shown} />
          <StatCard label="Premium (this page)" value={stats.premium} accent="#7E22CE" />
          <StatCard label="Stock units (this page)" value={stats.totalStock} accent="#15803D" />
        </div>

        <div className="mt-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <label className="relative block w-full sm:max-w-xs">
            <span className="sr-only">Search vehicles</span>
            <span className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-[#5B6472]">
              {NodeIcon('search', 'h-4 w-4')}
            </span>
            <input
              type="text"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              placeholder="Search by name, brand, or model"
              className="w-full rounded-full border border-slate-200 bg-white py-2.5 pl-10 pr-4 FleetOps-body text-[13.5px] text-[#0B0E14] placeholder:text-slate-400 transition-colors focus:border-[#0B0E14] focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/10"
            />
          </label>

          <div className="flex flex-wrap items-center gap-3">
            <SelectField
              icon="filter"
              label="Filter by type"
              value={typeFilter}
              onChange={withPageReset(setTypeFilter)}
              options={[
                { value: 'all', label: 'All types' },
                { value: 'Car', label: 'Car' },
                { value: 'SUV', label: 'SUV' },
                { value: 'Bike', label: 'Bike' },
                { value: 'Van', label: 'Van' },
                { value: 'Truck', label: 'Truck' },
              ]}
            />
            <SelectField
              icon="filter"
              label="Filter by category"
              value={categoryFilter}
              onChange={withPageReset(setCategoryFilter)}
              options={[
                { value: 'all', label: 'All categories' },
                { value: 'Normal', label: 'Normal' },
                { value: 'Moderate', label: 'Moderate' },
                { value: 'Premium', label: 'Premium' },
              ]}
            />
            <SelectField
              icon="sort"
              label="Sort by"
              value={sortBy}
              onChange={withPageReset(setSortBy)}
              options={[
                { value: 'name', label: 'Name' },
                { value: 'model', label: 'Model' },
                { value: 'created_at', label: 'Date added' },
              ]}
            />
            <SelectField
              icon="sort"
              label="Order"
              value={order}
              onChange={withPageReset(setOrder)}
              options={[
                { value: 'asc', label: 'Ascending' },
                { value: 'desc', label: 'Descending' },
              ]}
            />
          </div>
        </div>

        {isLoading ? (
          <div className="mt-6 rounded-2xl border border-slate-200 bg-white py-16 text-center">
            <p className="FleetOps-body text-[13.5px] text-[#5B6472]">Loading vehicles…</p>
          </div>
        ) : error ? (
          <div className="mt-6 rounded-2xl border border-red-100 bg-red-50 py-16 text-center">
            <p className="FleetOps-body text-[13.5px] text-[#DC2626]">Couldn&rsquo;t load vehicles. Please try again.</p>
          </div>
        ) : vehicles.length > 0 ? (
          <div className={`mt-6 overflow-hidden rounded-2xl border border-slate-200 bg-white ${isFetching ? 'opacity-60' : ''}`}>
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead>
                  <tr className="border-b border-slate-100 bg-[#F9FAFB]">
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Image</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Name</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Brand</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Model</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Type</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Category</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Price</th>
                    <th className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">Stock</th>
                    <th className="px-5 py-3" />
                  </tr>
                </thead>
                <tbody>
                  {vehicles.map((v) => {
                    const cat = categoryMeta[v.category] || { bg: '#F1F5F9', text: '#475569', dot: '#475569' }
                    return (
                      <tr key={v.id} className="border-b border-slate-100 transition-colors last:border-0 hover:bg-[#F9FAFB]">
                        <td className="px-5 py-3">
                          <Thumb src={v.image_url} alt={v.name} />
                        </td>
                        <td className="whitespace-nowrap px-5 py-3 FleetOps-display text-[14px] font-semibold text-[#0B0E14]">
                          {v.name}
                        </td>
                        <td className="whitespace-nowrap px-5 py-3 FleetOps-body text-[13px] text-[#5B6472]">{v.brand}</td>
                        <td className="whitespace-nowrap px-5 py-3 FleetOps-body text-[13px] text-[#5B6472]">{v.model}</td>
                        <td className="whitespace-nowrap px-5 py-3">
                          <span className="inline-flex items-center gap-1.5 FleetOps-mono text-[11px] text-[#5B6472]">
                            {NodeIcon(typeIconFor(v.type), 'h-4 w-4')}
                            {v.type}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-5 py-3">
                          <span
                            className="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 FleetOps-mono text-[10px] font-medium"
                            style={{ backgroundColor: cat.bg, color: cat.text }}
                          >
                            <span className="h-1.5 w-1.5 rounded-full" style={{ backgroundColor: cat.dot }} />
                            {v.category}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-5 py-3 FleetOps-body text-[13px] font-medium text-[#0B0E14]">
                          {formatPrice(v.price)}
                        </td>
                        <td className="whitespace-nowrap px-5 py-3 FleetOps-mono text-[12px] text-[#5B6472]">{v.stock}</td>
                        <td className="whitespace-nowrap px-5 py-3">
                          <div className="flex items-center justify-end gap-1.5">
                            <button
                              type="button"
                              onClick={() => setViewId(v.id)}
                              className="flex h-8 w-8 items-center justify-center rounded-lg text-[#5B6472] transition-colors hover:bg-[#F3F5F8] hover:text-[#0B0E14]"
                              aria-label={`View ${v.name}`}
                            >
                              {NodeIcon('eye', 'h-4 w-4')}
                            </button>
                            <button
                              type="button"
                              onClick={() => {
                                setFormError('')
                                setFormState({ mode: 'edit', initial: { ...emptyForm, ...v }, id: v.id })
                              }}
                              className="flex h-8 w-8 items-center justify-center rounded-lg text-[#5B6472] transition-colors hover:bg-[#F3F5F8] hover:text-[#0B0E14]"
                              aria-label={`Edit ${v.name}`}
                            >
                              {NodeIcon('edit', 'h-4 w-4')}
                            </button>
                            <button
                              type="button"
                              onClick={() => setDeleteTarget(v)}
                              className="flex h-8 w-8 items-center justify-center rounded-lg text-[#5B6472] transition-colors hover:bg-red-50 hover:text-[#DC2626]"
                              aria-label={`Delete ${v.name}`}
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
                setSearchInput('')
                setSearch('')
                setTypeFilter('all')
                setCategoryFilter('all')
                setPage(1)
              }}
              className="mt-1 rounded-full bg-[#0B0E14] px-4 py-2 FleetOps-body text-[13px] font-medium text-white transition-colors hover:bg-[#1a2030]"
            >
              Reset filters
            </button>
          </div>
        )}

        <Pagination page={page} totalPages={totalPages} hasNext={hasNext} onChange={setPage} />
      </main>

      {formState && (
        <VehicleFormModal
          mode={formState.mode}
          initial={formState.initial}
          isSaving={isCreating || isUpdating}
          errorMessage={formError}
          onCancel={() => setFormState(null)}
          onSave={(payload) => (formState.mode === 'create' ? handleCreate(payload) : handleUpdate(formState.id, payload))}
        />
      )}

      {viewId && <VehicleViewModal id={viewId} onClose={() => setViewId(null)} />}

      {deleteTarget && (
        <DeleteConfirm
          vehicle={deleteTarget}
          isDeleting={isDeleting}
          onCancel={() => setDeleteTarget(null)}
          onConfirm={() => handleDelete(deleteTarget.id)}
        />
      )}
    </div>
  )
}