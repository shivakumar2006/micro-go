import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useGetVehicleByIdQuery } from '../Redux/features/vehicles/vehicleApi';
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
  arrowLeft: (
    <>
      <path d="M19 12H5" />
      <path d="M11 18l-6-6 6-6" />
    </>
  ),
  star: <path d="M12 3.5l2.6 5.3 5.9.9-4.3 4.1 1 5.8L12 16.9 6.8 19.6l1-5.8L3.5 9.7l5.9-.9L12 3.5z" />,
  box: (
    <>
      <path d="M3 8l9-4 9 4-9 4-9-4z" />
      <path d="M3 8v8l9 4 9-4V8" />
      <path d="M12 12v8" />
    </>
  ),
  tag: (
    <>
      <path d="M11.5 3.5h6a1 1 0 0 1 1 1v6a1 1 0 0 1-.3.7l-8.5 8.5a1 1 0 0 1-1.4 0l-6.7-6.7a1 1 0 0 1 0-1.4l8.5-8.5a1 1 0 0 1 .7-.3z" />
      <circle cx="16" cy="8" r="1.3" fill="currentColor" stroke="none" />
    </>
  ),
  image: (
    <>
      <rect x="3" y="4.5" width="18" height="15" rx="1.5" />
      <circle cx="8.5" cy="9.5" r="1.5" />
      <path d="M3 16l5-5 4 4 3-3 6 6" />
    </>
  ),
  building: (
    <>
      <rect x="4" y="3.5" width="10" height="17" rx="1" />
      <path d="M14 9.5h6v11h-6" />
      <path d="M7.5 7.5h1M11 7.5h1M7.5 11h1M11 11h1M7.5 14.5h1M11 14.5h1" />
    </>
  ),
  layers: (
    <>
      <path d="M12 3.2l9 5-9 5-9-5 9-5z" />
      <path d="M3 13.2l9 5 9-5" />
    </>
  ),
}

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

// sample payload — replace with the real fetch result for /vehicles/:id

function formatPrice(amount) {
  return new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
    maximumFractionDigits: 0,
  }).format(amount)
}

function stockMeta(stock) {
  if (stock <= 0) return { label: 'Out of stock', dot: '#94A3B8', bg: '#94A3B814', text: '#5B6472' }
  if (stock <= 5) return { label: `Low stock — ${stock} left`, dot: '#D97706', bg: '#D9770614', text: '#B45309' }
  return { label: `In stock — ${stock} units`, dot: '#16A34A', bg: '#16A34A14', text: '#15803D' }
}

function SpecRow({ label, value }) {
  return (
    <div className="flex items-center justify-between border-b border-slate-100 py-3 last:border-0">
      <span className="fleetos-mono text-[11px] uppercase tracking-wide text-[#5B6472]">{label}</span>
      <span className="fleetos-body text-[13.5px] font-medium text-[#0B0E14]">{value}</span>
    </div>
  )
}

function VehicleImage({ src, alt }) {
  const [failed, setFailed] = useState(false)

  if (!src || failed) {
    return (
      <div className="flex aspect-[4/3] w-full flex-col items-center justify-center gap-2 rounded-[28px] border border-slate-200 bg-[#F3F5F8] text-[#5B6472]">
        {NodeIcon('image', 'h-8 w-8')}
        <span className="fleetos-mono text-[11px]">No image available</span>
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-[28px] border border-slate-200 bg-[#F3F5F8]">
      <img
        src={src}
        alt={alt}
        onError={() => setFailed(true)}
        className="aspect-[4/3] w-full object-cover"
      />
    </div>
  )
}

export default function VehicleDetail() {
  const {id} = useParams();
  const navigate = useNavigate();

  const {data, isLoading, error} = useGetVehicleByIdQuery(id);

  const vehicle = data;

  if (isLoading) {
    return <div>Loading...</div>;
  }
  
  if (error) {
    return <div>Error loading vehicle</div>;
  }
  
  if (!vehicle) {
      return <div>No vehicle found</div>;
  }

  const {
    name,
    brand,
    model,
    stock,
    price,
    description,
    image_url,
    category,
    type,
    createdAt,
} = vehicle;


  const stockInfo = stockMeta(stock)

  console.log("vehicle by id data : ", data);

  return (
    <div className="min-h-screen bg-[#F3F5F8] fleetos-body text-[#0B0E14]">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@500;600;700&family=Manrope:wght@400;500;600;700&family=IBM+Plex+Mono:wght@400;500&display=swap');
        .fleetos-display { font-family: 'Space Grotesk', ui-sans-serif, sans-serif; }
        .fleetos-body { font-family: 'Manrope', ui-sans-serif, sans-serif; }
        .fleetos-mono { font-family: 'IBM Plex Mono', ui-monospace, monospace; }
      `}</style>

      <main className="mx-auto max-w-6xl px-6 py-10">
        <a
          onClick={() => navigate("/vehicles")}
          className="inline-flex items-center gap-1.5 fleetos-body text-[13.5px] font-medium text-[#5B6472] transition-colors hover:text-[#0B0E14]"
        >
          {NodeIcon('arrowLeft', 'h-4 w-4')}
          Back to vehicles
        </a>

        <div className="mt-6 grid gap-10 lg:grid-cols-2">
          {/* image */}
          <div>
            <VehicleImage src={image_url} alt={`${brand} ${name} ${model}`} />
          </div>

          {/* info */}
          <div>
            <div className="flex flex-wrap items-center gap-2">
              <span className="inline-flex items-center gap-1.5 rounded-full border border-slate-200 bg-white px-3 py-1 fleetos-mono text-[10.5px] tracking-wide text-[#5B6472]">
                {NodeIcon('building', 'h-3.5 w-3.5')}
                {brand}
              </span>
              <span className="inline-flex items-center gap-1.5 rounded-full border border-slate-200 bg-white px-3 py-1 fleetos-mono text-[10.5px] tracking-wide text-[#5B6472]">
                {NodeIcon('layers', 'h-3.5 w-3.5')}
                {type}
              </span>
              {category && (
                <span className="inline-flex items-center gap-1.5 rounded-full border border-[#35455C]/25 bg-[#35455C]/[0.06] px-3 py-1 fleetos-mono text-[10.5px] tracking-wide text-[#35455C]">
                  {NodeIcon('star', 'h-3.5 w-3.5')}
                  {category}
                </span>
              )}
            </div>

            <h1 className="mt-4 fleetos-display text-3xl font-semibold tracking-tight text-[#0B0E14] sm:text-4xl">
              {brand} {name}
            </h1>
            <p className="mt-1 fleetos-mono text-[13px] tracking-wide text-[#5B6472]">{model}</p>

            <p className="mt-5 fleetos-display text-2xl font-semibold text-[#0B0E14]">{formatPrice(price)}</p>

            <span
              className="mt-3 inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 fleetos-mono text-[10.5px] font-medium"
              style={{ backgroundColor: stockInfo.bg, color: stockInfo.text }}
            >
              <span className="h-1.5 w-1.5 rounded-full" style={{ backgroundColor: stockInfo.dot }} />
              {stockInfo.label}
            </span>

            <p className="mt-6 fleetos-body text-[14.5px] leading-relaxed text-[#5B6472]">{description || "no description available"}</p>

            <div className="mt-6 flex flex-wrap items-center gap-3">
              <button
                type="button"
                disabled={stock <= 0}
                className="rounded-full bg-[#FF5A1F] px-6 py-3 fleetos-body text-[14px] font-semibold text-white shadow-[0_12px_24px_-8px_rgba(255,90,31,0.55)] transition-all hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:translate-y-0"
              >
                Add to fleet
              </button>
              <a
                onClick={() => navigate("/vehicles")}
                className="rounded-full border border-slate-200 bg-white px-6 py-3 fleetos-body text-[14px] font-semibold text-[#0B0E14] transition-colors hover:bg-[#F3F5F8]"
              >
                Back to vehicles
              </a>
            </div>

            {/* spec sheet */}
            <div className="mt-8 rounded-2xl border border-slate-200 bg-white px-5">
              <SpecRow label="Brand" value={brand} />
              <SpecRow label="Model" value={model} />
              <SpecRow label="Type" value={type} />
              <SpecRow label="Category" value={category} />
              <SpecRow label="Stock" value={`${stock} units`} />
              <SpecRow label="Price" value={formatPrice(price)} />
              <SpecRow
                  label="Created"
                  value={new Date(createdAt).toLocaleDateString()}
              />
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}