import {useState } from 'react'
import VehicleCard from '../components/VehicleCard';
import { useGetAllVehiclesQuery } from '../Redux/features/vehicles/vehicleApi';
import { useNavigate } from 'react-router-dom';

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
  van: (
    <>
      <rect x="2" y="9" width="11" height="7" rx="1" />
      <path d="M13 11.5h3.5L19 14v2h-6z" />
      <circle cx="6.2" cy="18" r="1.5" />
      <circle cx="16.5" cy="18" r="1.5" />
    </>
  ),
  chevronDown: <path d="M6 9l6 6 6-6" />,
  chevronLeft: <path d="M15 6l-6 6 6 6" />,
  chevronRight: <path d="M9 6l6 6-6 6" />,
  inbox: (
    <>
      <path d="M3.5 12h5l1.5 3h4l1.5-3h5" />
      <path d="M5 6.5h14L21 12v6a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18v-6l2-5.5z" />
    </>
  ),
}

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

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

function Pagination({ page, totalPages, onChange }) {
  if (totalPages <= 1) return null
  const pages = Array.from({ length: totalPages }, (_, i) => i + 1)
  const pagerBtn = 'flex h-9 w-9 items-center justify-center rounded-full FleetOps-mono text-[13px] transition-colors'

  return (
    <div className="mt-10 flex items-center justify-center gap-1.5">
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

export default function Vehicles() {
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')
  const [typeFilter, setTypeFilter] = useState('all')
  const [sortBy, setSortBy] = useState('name-asc')
  const [page, setPage] = useState(1)
  const navigate = useNavigate();

  const { data, isLoading, error } = useGetAllVehiclesQuery({
    page,
    limit: PER_PAGE,
    search,
    sortBy,
    order: sortBy === "name-desc" ? "desc" : "asc",
    type: typeFilter === "all" ? "" : typeFilter,
    category: statusFilter === "all" ? "" : statusFilter,
  });

  if (isLoading) {
    return "Loading..."
  }

  if (error) {
    return "error..."
  }

  const vehicles = data?.data || [];

  const withPageReset = (setter) => (val) => {
    setter(val)
    setPage(1)
  }

  const sortMap = {
    "name-asc": {
        sortBy: "name",
        order: "asc",
    },

    "name-desc": {
        sortBy: "name",
        order: "desc",
    },

    "model-asc": {
        sortBy: "model",
        order: "asc",
    },

    "model-desc": {
        sortBy: "model",
        order: "desc",
    },

    "created-desc": {
        sortBy: "created_at",
        order: "desc",
    },

    "created-asc": {
        sortBy: "created_at",
        order: "asc",
    },
}

const { sortBy: backendSortBy, order } =
  sortMap[sortBy];

  // wire this up to a router (e.g. navigate(`/vehicles/${id}`)) to open
  // VehicleDetail.jsx for the clicked vehicle
  const handleOpenVehicle = (id) => {
    navigate(`/vehicles/details/${id}`)
    console.log('open vehicle', id)
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
            <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">fleet</span>
            <h1 className="mt-1 FleetOps-display text-2xl font-semibold tracking-tight text-[#0B0E14] sm:text-3xl">
              Vehicles
            </h1>
          </div>
          <p className="FleetOps-mono text-[12px] text-[#5B6472]">{data?.total || 0} vehicles found</p>
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
                { value: "all", label: "All Categories" },
                { value: "Normal", label: "Normal" },
                { value: "Moderate", label: "Moderate" },
                { value: "Premium", label: "Premium" },
              ]}
            />
            <SelectField
              icon="van"
              label="Filter by type"
              value={typeFilter}
              onChange={withPageReset(setTypeFilter)}
              options={[
                { value: "all", label: "All Types" },
                { value: "Car", label: "Car" },
                { value: "Bike", label: "Bike" },
                { value: "Truck", label: "Truck" },
                { value: "SUV", label: "SUV" },
                { value: "Van", label: "Van" },
                { value: "Bus", label: "Bus" },
                { value: "Other", label: "Other" },
              ]}
            />
            <SelectField
              icon="sort"
              label="Sort by"
              value={sortBy}
              onChange={setSortBy}
              options={[
                { value: "name-asc", label: "Name (A-Z)" },
                { value: "name-desc", label: "Name (Z-A)" },
                { value: "model-asc", label: "Model (A-Z)" },
                { value: "model-desc", label: "Model (Z-A)" },
                { value: "created-desc", label: "Newest" },
                { value: "created-asc", label: "Oldest" },
              ]}
            />
          </div>
        </div>

        {vehicles.length > 0 ? (
          <>
            <p className="mt-6 FleetOps-mono text-[11px] text-[#5B6472]">
              Showing {(page - 1) * PER_PAGE + 1} - {Math.min(page * PER_PAGE, data?.total || 0)} of {data?.total || 0}
            </p>
            <div className="mt-3 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
              {vehicles.map((vehicle) => (
                <VehicleCard
                  key={vehicle.id}
                  vehicle={vehicle}
                  onClick={() => handleOpenVehicle(vehicle.id)}
                />
              ))}
            </div>
          </>
        ) : (
          <div className="mt-16 flex flex-col items-center gap-3 rounded-2xl border border-dashed border-slate-300 py-16 text-center">
            {NodeIcon('inbox', 'h-8 w-8 text-[#5B6472]')}
            <p className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">No vehicles match your filters</p>
            <p className="FleetOps-body text-[13.5px] text-[#5B6472]">Try a different search term or reset the filters.</p>
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

        <Pagination
          page={page}
          totalPages={data?.total_pages || 1}
          onChange={setPage}
        />
      </main>
    </div>
  )
}