import { useState } from "react";
import { useNavigate } from "react-router-dom";

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
  monitor: (
    <>
      <rect x="2.5" y="4.5" width="19" height="13" rx="1.5" />
      <path d="M2.5 8h19" />
      <circle cx="5.4" cy="6.25" r="0.55" fill="currentColor" stroke="none" />
    </>
  ),
  balance: (
    <>
      <path d="M12 4v4" />
      <path d="M5 15v-3.5A2 2 0 0 1 7 9.5h10A2 2 0 0 1 19 11.5V15" />
      <circle cx="5" cy="18" r="2" />
      <circle cx="12" cy="18" r="2" />
      <circle cx="19" cy="18" r="2" />
    </>
  ),
  gateway: (
    <>
      <path d="M8 3.5h8l4 6.5-4 6.5H8l-4-6.5 4-6.5z" />
      <path d="M9 10h6" />
    </>
  ),
  shield: (
    <>
      <path d="M12 3.2l7 2.8v5c0 4.6-3 7.6-7 9-4-1.4-7-4.4-7-9v-5l7-2.8z" />
      <path d="M9 12l2 2 4-4" />
    </>
  ),
  truck: (
    <>
      <rect x="2" y="9" width="11" height="7" rx="1" />
      <path d="M13 11.5h3.5L19 14v2h-6z" />
      <circle cx="6.2" cy="18" r="1.5" />
      <circle cx="16.5" cy="18" r="1.5" />
    </>
  ),
  cart: (
    <>
      <path d="M3 4h2l2.2 10a2 2 0 0 0 2 1.6h6.6a2 2 0 0 0 2-1.6L19.5 8H6" />
      <circle cx="9" cy="20" r="1.3" />
      <circle cx="17" cy="20" r="1.3" />
    </>
  ),
  package: (
    <>
      <path d="M3 8l9-4 9 4-9 4-9-4z" />
      <path d="M3 8v8l9 4 9-4V8" />
      <path d="M12 12v8" />
    </>
  ),
  card: (
    <>
      <rect x="2.5" y="5.5" width="19" height="13" rx="2" />
      <path d="M2.5 9.5h19" />
      <path d="M6 14.5h4" />
    </>
  ),
  boxes: (
    <>
      <path d="M3 7.5l4-1.8 4 1.8-4 1.8-4-1.8z" />
      <path d="M3 7.5v4.8l4 1.8 4-1.8V7.5" />
      <path d="M13 7.5l4-1.8 4 1.8-4 1.8-4-1.8z" />
      <path d="M13 7.5v4.8l4 1.8 4-1.8V7.5" />
    </>
  ),
  bell: (
    <>
      <path d="M6 10.2a6 6 0 0 1 12 0c0 3.8 1.3 5.3 1.3 5.3H4.7S6 14 6 10.2z" />
      <path d="M10 18a2 2 0 0 0 4 0" />
    </>
  ),
  chart: (
    <>
      <path d="M4 20V11" />
      <path d="M11 20V4" />
      <path d="M18 20v-6.5" />
    </>
  ),
  cloud: (
    <>
      <path d="M7 17.5a4 4 0 0 1-.5-7.97 5 5 0 0 1 9.65-1.62A4.5 4.5 0 0 1 17.5 17.5H7z" />
    </>
  ),
  layers: (
    <>
      <path d="M12 3.2l9 5-9 5-9-5 9-5z" />
      <path d="M3 13.2l9 5 9-5" />
    </>
  ),
  arrowRight: (
    <>
      <path d="M4 12h16" />
      <path d="M14 6l6 6-6 6" />
    </>
  ),
  menu: <path d="M4 7h16M4 12h16M4 17h16" />,
  close: <path d="M6 6l12 12M18 6L6 18" />,
  bolt: <path d="M13 2 L4 14 H10 L9 22 L20 10 H13 Z" fill="currentColor" stroke="none" />,
}


export default function Nav() {
  const [open, setOpen] = useState(false)
  const links = [
    { label: 'Product', href: '/vehicles' },
    { label: 'Architecture', href: '#architecture' },
    { label: 'Docs', href: '#' },
  ]
  const navigate = useNavigate();

  return (
    <header className="sticky top-0 z-40 border-b border-slate-200/70 bg-white/80 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
        <a href="#top" className="flex items-center gap-2.5">
          <span className="relative flex h-7 w-7 items-center justify-center">
            <span className="absolute inset-0 rotate-45 rounded-[7px] bg-[#0B0E14]" />
            <span className="relative h-1.5 w-1.5 rounded-full bg-[#FF5A1F]" />
          </span>
          <span className="FleetOps-display text-[17px] font-semibold tracking-tight text-[#0B0E14]">
            FleetOps
          </span>
        </a>

        <nav className="hidden items-center gap-8 md:flex">
          {links.map((l) => (
            <a
              key={l.label}
              href={l.href}
              className="FleetOps-body text-[13.5px] font-medium text-[#5B6472] transition-colors hover:text-[#0B0E14]"
            >
              {l.label}
            </a>
          ))}
        </nav>

        <div className="hidden items-center gap-3 md:flex cursor-pointer">
          <a
            onClick={() => navigate("/login")}
            className="FleetOps-body text-[13.5px] font-medium text-[#5B6472] transition-colors hover:text-[#0B0E14]"
          >
            Sign in
          </a>
          <a
            onClick={() => navigate("/login")}
            className="rounded-full bg-[#0B0E14] px-4 py-2 FleetOps-body text-[13.5px] font-medium text-white transition-colors hover:bg-[#1a2030]"
          >
            Request access
          </a>
        </div>

        <button
          className="flex h-9 w-9 items-center justify-center rounded-lg text-[#0B0E14] md:hidden"
          onClick={() => setOpen((v) => !v)}
          aria-label="Toggle menu"
        >
          <Icon path={open ? icons.close : icons.menu} className="h-5 w-5" />
        </button>
      </div>

      {open && (
        <div className="flex flex-col gap-1 border-t border-slate-200 bg-white px-6 py-4 md:hidden">
          {links.map((l) => (
            <a key={l.label} href={l.href} className="py-2 FleetOps-body text-sm text-[#0B0E14]">
              {l.label}
            </a>
          ))}
          <a href="#cta" className="mt-2 rounded-full bg-[#0B0E14] px-4 py-2.5 text-center FleetOps-body text-sm font-medium text-white">
            Request access
          </a>
        </div>
      )}
    </header>
  )
}
