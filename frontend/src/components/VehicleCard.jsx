import React, { useState } from "react";
import { useParams } from "react-router-dom";
import { useGetVehicleByIdQuery } from "../Redux/features/vehicles/vehicleApi";

const Icon = ({ path, className = "w-5 h-5" }) => (
  <svg
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="1.6"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    {path}
  </svg>
);

const icons = {
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
  stock: (
    <>
      <rect x="2.5" y="8" width="17" height="8" rx="1.5" />
      <path d="M21.5 10.5v3" />
    </>
  ),
  calendar: (
    <>
      <rect x="3.5" y="5" width="17" height="15" rx="1.5" />
      <path d="M3.5 9.5h17" />
      <path d="M8 3v4" />
      <path d="M16 3v4" />
    </>
  ),
  image: (
    <>
      <rect x="3" y="4.5" width="18" height="15" rx="1.5" />
      <circle cx="8.5" cy="9.5" r="1.5" />
      <path d="M3 16l5-5 4 4 3-3 6 6" />
    </>
  ),
};

const NodeIcon = (name, className) => (
  <Icon path={icons[name]} className={className} />
);

const categoryMeta = {
  Normal: { bg: "#DBEAFE", text: "#1D4ED8", dot: "#1D4ED8", label: "Normal" },
  Moderate: { bg: "#FEF3C7", text: "#B45309", dot: "#B45309", label: "Moderate" },
  Premium: { bg: "#F3E8FF", text: "#7E22CE", dot: "#7E22CE", label: "Premium" },
};

function formatPrice(amount) {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    maximumFractionDigits: 0,
  }).format(amount);
}

function formatDate(date) {
  return new Date(date).toLocaleDateString("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

function VehicleImage({ src, alt, typeIcon, typeLabel }) {
  const [failed, setFailed] = useState(false);
  const showImage = src && !failed;

  return (
    <div className="relative aspect-[4/3] w-full overflow-hidden bg-[#F3F5F8]">
      {showImage ? (
        <img
          src={src}
          alt={alt}
          onError={() => setFailed(true)}
          className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
        />
      ) : (
        <div className="flex h-full w-full flex-col items-center justify-center gap-2 text-[#5B6472]">
          {NodeIcon("image", "h-7 w-7")}
          <span className="fleetos-mono text-[10.5px]">No image available</span>
        </div>
      )}

      {typeLabel && (
        <span className="absolute bottom-3 left-3 inline-flex items-center gap-1.5 rounded-full bg-[#0B0E14]/75 px-2.5 py-1 fleetos-mono text-[10px] uppercase tracking-wide text-white backdrop-blur">
          {NodeIcon(typeIcon, "h-3.5 w-3.5")}
          {typeLabel}
        </span>
      )}
    </div>
  );
}

export default function VehicleCard({ vehicle, onClick }) {
  const s =
    categoryMeta[vehicle.category] || {
      bg: "#F1F5F9",
      text: "#475569",
      dot: "#475569",
      label: vehicle.category || "Unknown",
    };

  const clickable = typeof onClick === "function";
  const typeKey = vehicle.type?.toLowerCase();
  const typeIcon = icons[typeKey] ? typeKey : "car";
  const stockPercentage = Math.min(vehicle.stock, 100);

  const {id} = useParams();

  const {data, isLoading, error} = useGetVehicleByIdQuery(id);

  console.log("vehicle by id data : ", data);

  return (
    <div
      onClick={onClick}
      role={clickable ? "button" : undefined}
      tabIndex={clickable ? 0 : undefined}
      onKeyDown={
        clickable
          ? (e) => {
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                onClick();
              }
            }
          : undefined
      }
      className={`group overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-[0_1px_2px_rgba(15,23,42,0.04)] transition-all hover:-translate-y-0.5 hover:border-slate-300 hover:shadow-[0_16px_32px_-16px_rgba(15,23,42,0.2)] ${
        clickable ? "cursor-pointer focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/15" : ""
      }`}
    >
      {/* image + overlaid badges */}
      <div className="relative">
        <VehicleImage
          src={vehicle.image_url}
          alt={`${vehicle.brand} ${vehicle.name}`}
          typeIcon={typeIcon}
          typeLabel={vehicle.type}
        />
        <span
          className="absolute right-3 top-3 inline-flex items-center gap-1.5 rounded-full bg-white/90 px-2.5 py-1 fleetos-mono text-[10px] font-medium backdrop-blur"
          style={{ color: s.text }}
        >
          <span className="h-1.5 w-1.5 rounded-full" style={{ background: s.dot }} />
          {s.label}
        </span>
      </div>

      <div className="p-5">
        <div>
          <h3 className="fleetos-display text-[16px] font-semibold text-[#0B0E14]">{vehicle.name}</h3>
          <p className="fleetos-mono text-[11px] tracking-wide text-[#5B6472]">
            {vehicle.brand} · {vehicle.model}
          </p>
        </div>

        <p className="mt-3 fleetos-display text-xl font-semibold text-[#0B0E14]">
          {formatPrice(vehicle.price)}
        </p>

        {vehicle.description && (
          <p className="mt-3 line-clamp-2 fleetos-body text-[13px] leading-relaxed text-[#5B6472]">
            {vehicle.description}
          </p>
        )}

        <div className="mt-4 flex items-center justify-between border-t border-slate-100 pt-4">
          <div className="flex items-center gap-2">
            {NodeIcon("stock", "h-4 w-4 text-[#5B6472]")}
            <div className="h-1.5 w-16 overflow-hidden rounded-full bg-slate-100">
              <div
                className="h-full rounded-full bg-[#0B0E14]"
                style={{ width: `${stockPercentage}%` }}
              />
            </div>
            <span className="fleetos-mono text-[11px] text-[#5B6472]">{vehicle.stock}</span>
          </div>

          {vehicle.createdAt && (
            <span className="flex items-center gap-1.5 fleetos-mono text-[10.5px] text-[#5B6472]">
              {NodeIcon("calendar", "h-3.5 w-3.5")}
              {formatDate(vehicle.createdAt)}
            </span>
          )}
        </div>
      </div>
    </div>
  );
}