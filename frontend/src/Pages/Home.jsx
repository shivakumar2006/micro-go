import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

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

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

// ---------------------------------------------------------------------------
// Small shared building blocks
// ---------------------------------------------------------------------------

function GoBadge() {
  return (
    <span className="inline-flex items-center gap-1 rounded-full border border-[#00ADD8]/30 bg-[#00ADD8]/[0.07] px-2 py-0.5 FleetOps-mono text-[10px] tracking-wide text-[#0a7f99]">
      <span className="h-1.5 w-1.5 rounded-full bg-[#00ADD8]" />
      GO
    </span>
  )
}

function ServiceCard({ icon, name, blurb }) {
  return (
    <div className="group relative flex flex-col gap-3 rounded-2xl border border-slate-200 bg-white p-4 shadow-[0_1px_2px_rgba(15,23,42,0.04)] transition-all hover:-translate-y-0.5 hover:border-slate-300 hover:shadow-[0_12px_24px_-12px_rgba(15,23,42,0.18)]">
      <div className="flex items-center justify-between">
        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-[#0B0E14] text-white">
          {NodeIcon(icon, 'w-[18px] h-[18px]')}
        </div>
        <GoBadge />
      </div>
      <div>
        <p className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">{name}</p>
        <p className="mt-0.5 FleetOps-mono text-[11px] leading-relaxed text-[#5B6472]">{blurb}</p>
      </div>
    </div>
  )
}

function InfraCard({ icon, name, sub }) {
  return (
    <div className="flex items-center gap-3 rounded-2xl border border-slate-200 bg-white px-5 py-3.5 shadow-[0_1px_2px_rgba(15,23,42,0.04)]">
      <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-[#35455C] text-white">
        {NodeIcon(icon, 'w-[18px] h-[18px]')}
      </div>
      <div>
        <p className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">{name}</p>
        <p className="FleetOps-mono text-[11px] text-[#5B6472]">{sub}</p>
      </div>
    </div>
  )
}

function RedisCard() {
  return (
    <div className="flex items-center gap-3 rounded-2xl border border-[#DC382D]/25 bg-white px-4 py-3 shadow-[0_1px_2px_rgba(15,23,42,0.04)]">
      <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-[#DC382D] text-white">
        {NodeIcon('bolt', 'w-[16px] h-[16px]')}
      </div>
      <div>
        <p className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">Redis</p>
        <p className="FleetOps-mono text-[11px] text-[#5B6472]">shared cache — vehicle &amp; cart</p>
      </div>
    </div>
  )
}

// vertical connector rail with a traveling "signal" dot
function VRail({ height = 'h-10', delay = 0 }) {
  return (
    <div className={`relative w-px ${height} bg-slate-300`}>
      <span
        className="FleetOps-flow-v absolute left-1/2 h-1.5 w-1.5 -translate-x-1/2 rounded-full bg-[#FF5A1F]"
        style={{ animationDelay: `${delay}s` }}
      />
    </div>
  )
}

// horizontal fan-out rail spanning between two column centers, with drop
// ticks into each of `count` evenly-spaced columns below it
function FanOut({ count, insetPct }) {
  return (
    <div className="relative h-10 w-full">
      <div
        className="absolute top-0 h-px bg-slate-300"
        style={{ left: `${insetPct}%`, right: `${insetPct}%` }}
      >
        <span
          className="FleetOps-flow-h absolute top-1/2 h-1.5 w-1.5 -translate-y-1/2 rounded-full bg-[#FF5A1F]"
        />
      </div>
      <div className="grid h-full" style={{ gridTemplateColumns: `repeat(${count}, minmax(0,1fr))` }}>
        {Array.from({ length: count }).map((_, i) => (
          <div key={i} className="flex justify-center">
            <div className="h-full w-px bg-slate-300" />
          </div>
        ))}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Nav
// ---------------------------------------------------------------------------

function Nav() {
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

// ---------------------------------------------------------------------------
// Hero — headline + a small "dispatch board" visual grounded in fleet telemetry
// ---------------------------------------------------------------------------

function DispatchBoard() {
  return (
    <div className="relative aspect-[4/5] w-full max-w-sm overflow-hidden rounded-[28px] border border-slate-200 bg-white shadow-[0_30px_60px_-30px_rgba(15,23,42,0.25)] sm:aspect-square lg:aspect-[4/5]">
      {/* faint dot grid */}
      <div
        className="absolute inset-0 opacity-[0.35]"
        style={{
          backgroundImage: 'radial-gradient(#CBD3E0 1px, transparent 1px)',
          backgroundSize: '18px 18px',
        }}
      />

      <div className="relative flex h-full flex-col justify-between p-6">
        <div className="flex items-center justify-between">
          <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">
            live_dispatch.view
          </span>
          <span className="flex items-center gap-1.5 FleetOps-mono text-[11px] text-[#5B6472]">
            <span className="h-1.5 w-1.5 rounded-full bg-[#FF5A1F] FleetOps-blink" />
            streaming
          </span>
        </div>

        <svg viewBox="0 0 280 220" className="w-full flex-1">
          <path
            d="M20 190 C 70 150, 90 130, 120 110 S 190 60, 250 30"
            fill="none"
            stroke="#CBD3E0"
            strokeWidth="2"
            strokeDasharray="1 9"
            strokeLinecap="round"
          />
          <circle cx="20" cy="190" r="4" fill="#0B0E14" />
          <circle cx="140" cy="97" r="4" fill="#0B0E14" />
          <circle cx="250" cy="30" r="4" fill="#0B0E14" />
          <circle r="5" fill="#FF5A1F" className="FleetOps-route-dot">
            <animateMotion
              dur="4.5s"
              repeatCount="indefinite"
              path="M20 190 C 70 150, 90 130, 120 110 S 190 60, 250 30"
            />
          </circle>
        </svg>

        <div className="flex gap-3">
          <div className="flex-1 rounded-xl border border-slate-200 bg-white/90 p-3 FleetOps-float">
            <p className="FleetOps-mono text-[10px] text-[#5B6472]">vehicles online</p>
            <p className="FleetOps-display text-lg font-semibold text-[#0B0E14]">128</p>
          </div>
          <div className="flex-1 rounded-xl border border-slate-200 bg-white/90 p-3 FleetOps-float" style={{ animationDelay: '0.6s' }}>
            <p className="FleetOps-mono text-[10px] text-[#5B6472]">orders in flight</p>
            <p className="FleetOps-display text-lg font-semibold text-[#0B0E14]">37</p>
          </div>
        </div>
      </div>
    </div>
  )
}

function Hero() {
  return (
    <section id="top" className="relative overflow-hidden bg-white pt-16 sm:pt-20">
      <div className="mx-auto grid max-w-6xl items-center gap-14 px-6 pb-20 lg:grid-cols-[1.05fr_0.95fr] lg:pb-28">
        <div>
          <span className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-[#F3F5F8] px-3 py-1 FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">
            <span className="h-1.5 w-1.5 rounded-full bg-[#FF5A1F]" />
            private beta — built on go &amp; kafka
          </span>

          <h1 className="mt-6 FleetOps-display text-[2.6rem] font-semibold leading-[1.08] tracking-tight text-[#0B0E14] sm:text-6xl">
            Every vehicle.
            <br />
            Every order.
            <br />
            <span className="relative inline-block">
              One control plane.
              <span className="absolute -bottom-1 left-0 h-[6px] w-full bg-[#FF5A1F]/25" />
            </span>
          </h1>

          <p className="mt-6 max-w-md FleetOps-body text-[16px] leading-relaxed text-[#5B6472]">
            FleetOps is the backend for fleet-based commerce — live vehicle state, cart and
            checkout, and payments, wired together over a Go microservice mesh and a
            shared event bus.
          </p>

          <div className="mt-9 flex flex-wrap items-center gap-4">
            <a
              href="#cta"
              className="group inline-flex items-center gap-2 rounded-full bg-[#FF5A1F] px-6 py-3.5 FleetOps-body text-[14.5px] font-semibold text-white shadow-[0_12px_24px_-8px_rgba(255,90,31,0.55)] transition-transform hover:-translate-y-0.5"
            >
              Request access
              <Icon path={icons.arrowRight} className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
            </a>
            <a
              href="#architecture"
              className="inline-flex items-center gap-2 FleetOps-body text-[14.5px] font-semibold text-[#0B0E14] underline decoration-slate-300 underline-offset-4 hover:decoration-[#0B0E14]"
            >
              See how it&rsquo;s built
            </a>
          </div>

          <div className="mt-14 grid max-w-md grid-cols-3 gap-6 border-t border-slate-200 pt-6">
            {[
              ['5', 'core services'],
              ['3', 'async consumers'],
              ['1', 'event bus'],
            ].map(([n, l]) => (
              <div key={l}>
                <p className="FleetOps-display text-2xl font-semibold text-[#0B0E14]">{n}</p>
                <p className="FleetOps-mono text-[11px] text-[#5B6472]">{l}</p>
              </div>
            ))}
          </div>
        </div>

        <div className="flex justify-center lg:justify-end">
          <DispatchBoard />
        </div>
      </div>
    </section>
  )
}

// ---------------------------------------------------------------------------
// Features
// ---------------------------------------------------------------------------

function Features() {
  const items = [
    {
      icon: 'truck',
      title: 'Live vehicle state',
      body: 'The vehicle service streams position and status straight into cart and order decisions, so nothing checks out against a vehicle that isn\u2019t actually available.',
    },
    {
      icon: 'cart',
      title: 'One checkout, three services',
      body: 'Cart, orders and payments call each other directly over internal service calls, not through the gateway, so checkout stays fast and consistent.',
    },
    {
      icon: 'layers',
      title: 'Kafka behind every order',
      body: 'Once an order is placed, Kafka fans it out to inventory, notifications and analytics \u2014 independently, without the order service waiting on any of them.',
    },
    {
      icon: 'gateway',
      title: 'A Go mesh underneath',
      body: 'Every core service is written in Go and sits behind a single API gateway and load balancer, so the client only ever talks to one front door.',
    },
  ]

  return (
    <section id="product" className="bg-[#F3F5F8] py-24">
      <div className="mx-auto max-w-6xl px-6">
        <div className="max-w-lg">
          <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">what FleetOps does</span>
          <h2 className="mt-3 FleetOps-display text-3xl font-semibold tracking-tight text-[#0B0E14] sm:text-4xl">
            Built around one order, end to end.
          </h2>
        </div>

        <div className="mt-12 grid gap-5 sm:grid-cols-2">
          {items.map((it) => (
            <div key={it.title} className="rounded-2xl border border-slate-200 bg-white p-6">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#0B0E14] text-white">
                {NodeIcon(it.icon, 'h-5 w-5')}
              </div>
              <h3 className="mt-4 FleetOps-display text-[17px] font-semibold text-[#0B0E14]">{it.title}</h3>
              <p className="mt-2 FleetOps-body text-[14px] leading-relaxed text-[#5B6472]">{it.body}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

// ---------------------------------------------------------------------------
// Architecture — the request-flow diagram
// ---------------------------------------------------------------------------

function Architecture() {
  const coreServices = [
    { icon: 'shield', name: 'Auth', blurb: 'sessions & tokens' },
    { icon: 'truck', name: 'Vehicle', blurb: 'position & status' },
    { icon: 'cart', name: 'Cart', blurb: 'line items & pricing' },
    { icon: 'package', name: 'Orders', blurb: 'order lifecycle' },
    { icon: 'card', name: 'Payment', blurb: 'charges & refunds' },
  ]
  const consumers = [
    { icon: 'boxes', name: 'Inventory', blurb: 'stock reconciliation' },
    { icon: 'bell', name: 'Notification', blurb: 'push, sms, email' },
    { icon: 'chart', name: 'Analytics', blurb: 'metrics & events' },
  ]

  return (
    <section id="architecture" className="bg-white py-24">
      <div className="mx-auto max-w-6xl px-6">
        <div className="max-w-xl">
          <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">request lifecycle</span>
          <h2 className="mt-3 FleetOps-display text-3xl font-semibold tracking-tight text-[#0B0E14] sm:text-4xl">
            How a request actually moves.
          </h2>
          <p className="mt-3 FleetOps-body text-[15px] leading-relaxed text-[#5B6472]">
            From the browser to the event bus — this is the exact path every order takes
            through FleetOps.
          </p>
        </div>

        <div className="mt-14 rounded-[28px] border border-slate-200 bg-[#F9FAFB] p-6 sm:p-10">
          {/* client */}
          <div className="flex justify-center">
            <InfraCard icon="monitor" name="Client" sub="react app" />
          </div>
          <div className="flex justify-center">
            <VRail delay={0} />
          </div>

          {/* load balancer */}
          <div className="flex justify-center">
            <InfraCard icon="balance" name="Load balancer" sub="traffic distribution" />
          </div>
          <div className="flex justify-center">
            <VRail delay={0.3} />
          </div>

          {/* gateway */}
          <div className="flex justify-center">
            <InfraCard icon="gateway" name="API gateway" sub="single front door" />
          </div>

          {/* fan-out to 5 core services */}
          <div className="hidden sm:block">
            <FanOut count={5} insetPct={10} />
          </div>
          <div className="flex justify-center sm:hidden">
            <VRail height="h-8" delay={0.2} />
          </div>

          <div className="grid grid-cols-2 gap-3 sm:grid-cols-5 sm:gap-4">
            {coreServices.map((s) => (
              <ServiceCard key={s.name} icon={s.icon} name={s.name} blurb={s.blurb} />
            ))}
          </div>

          {/* internal service mesh: vehicle <-> cart <-> orders <-> payment */}
          <div className="relative mt-3 hidden h-8 sm:block">
            <div
              className="absolute top-0 h-px bg-[#35455C]/40"
              style={{ left: '30%', right: '10%' }}
            />
            <div className="grid h-full grid-cols-5">
              {[false, true, true, true, true].map((connected, i) => (
                <div key={i} className="flex justify-center">
                  {connected && <div className="h-full w-px bg-[#35455C]/40" />}
                </div>
              ))}
            </div>
            <span className="absolute -top-2 left-1/2 -translate-x-1/2 rounded-full border border-[#35455C]/25 bg-white px-2 py-0.5 FleetOps-mono text-[9.5px] tracking-wide text-[#35455C]">
              internal calls
            </span>
          </div>
          <p className="mt-2 text-center FleetOps-mono text-[11px] text-[#5B6472] sm:hidden">
            vehicle ↔ cart ↔ orders ↔ payment — internal calls
          </p>

          {/* redis: shared cache for vehicle + cart only */}
          <div className="relative mt-2 hidden sm:block">
            <svg
              className="absolute left-0 top-0 h-12 w-full"
              preserveAspectRatio="none"
              viewBox="0 0 100 48"
            >
              <path
                d="M30 0 L40 40"
                fill="none"
                stroke="#DC382D"
                strokeOpacity="0.4"
                strokeWidth="1.4"
                vectorEffect="non-scaling-stroke"
              />
              <path
                d="M50 0 L40 40"
                fill="none"
                stroke="#DC382D"
                strokeOpacity="0.4"
                strokeWidth="1.4"
                vectorEffect="non-scaling-stroke"
              />
            </svg>
            <div className="flex justify-center pt-12">
              <div style={{ marginLeft: '-10%' }}>
                <RedisCard />
              </div>
            </div>
          </div>
          <div className="mt-2 flex flex-col items-center gap-1.5 sm:hidden">
            <RedisCard />
            <p className="text-center FleetOps-mono text-[11px] text-[#5B6472]">
              shared by vehicle &amp; cart
            </p>
          </div>

          {/* orders -> kafka (only under the 3rd of 5 columns) */}
          <div className="mt-2 grid grid-cols-2 sm:grid-cols-5">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="flex justify-center">
                {i === 3 && <VRail height="h-10" delay={0.6} />}
              </div>
            ))}
          </div>
          <div className="flex justify-center sm:hidden">
            <VRail height="h-8" delay={0.6} />
          </div>

          {/* kafka band */}
          <div className="flex justify-center">
            <div className="flex items-center gap-3 rounded-2xl border border-[#35455C]/30 bg-[#0B0E14] px-6 py-4">
              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-white/10 text-white">
                {NodeIcon('layers', 'w-[18px] h-[18px]')}
              </div>
              <div>
                <p className="FleetOps-display text-[15px] font-semibold text-white">Kafka</p>
                <p className="FleetOps-mono text-[11px] text-white/60">event streaming bus</p>
              </div>
            </div>
          </div>

          {/* fan-out to 3 consumers */}
          <div className="hidden sm:block">
            <FanOut count={3} insetPct={16.6667} />
          </div>
          <div className="flex justify-center sm:hidden">
            <VRail height="h-8" delay={0.8} />
          </div>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
            {consumers.map((c) => (
              <ServiceCard key={c.name} icon={c.icon} name={c.name} blurb={c.blurb} />
            ))}
          </div>
        </div>

        <p className="mx-auto mt-6 max-w-2xl text-center FleetOps-body text-[13.5px] leading-relaxed text-[#5B6472]">
          Vehicle, cart, orders and payment call each other directly for checkout. Vehicle
          and cart also share a Redis cache for fast reads. Everything downstream of an
          order — inventory, notifications, analytics — reacts through Kafka instead of a
          direct call.
        </p>
      </div>
    </section>
  )
}

// ---------------------------------------------------------------------------
// Engineering highlights — the production stack, shown as real vendor marks
// ------------------------------------------------------------------
// Two items here (RBAC, Azure) aren't third-party products with a brand
// mark to license — RBAC is an access-control pattern, and Microsoft
// withdrew the Azure logo from open icon sets over trademark terms. Both
// stay visually honest by using the site's own steel-colored line icons
// instead of a counterfeited logo. Everything else below is the real,
// openly-licensed brand mark (simple-icons, CC0) in its official color.
// ---------------------------------------------------------------------------

const JWT_PATH =
  'M10.2 0v6.456L12 8.928l1.8-2.472V0zm3.6 6.456v3.072l2.904-.96L20.52 3.36l-2.928-2.136zm2.904 2.112l-1.8 2.496 2.928.936 6.144-1.992-1.128-3.432zM17.832 12l-2.928.936 1.8 2.496 6.144 1.992 1.128-3.432zm-1.128 3.432l-2.904-.96v3.072l3.792 5.232 2.928-2.136zM13.8 17.544L12 15.072l-1.8 2.472V24h3.6zm-3.6 0v-3.072l-2.904.96L3.48 20.64l2.928 2.136zm-2.904-2.112l1.8-2.496L6.168 12 .024 13.992l1.128 3.432zM6.168 12l2.928-.936-1.8-2.496-6.144-1.992-1.128 3.432zm1.128-3.432l2.904.96V6.456L6.408 1.224 3.48 3.36Z'

const REDIS_PATH =
  'M22.71 13.145c-1.66 2.092-3.452 4.483-7.038 4.483-3.203 0-4.397-2.825-4.48-5.12.701 1.484 2.073 2.685 4.214 2.63 4.117-.133 6.94-3.852 6.94-7.239 0-4.05-3.022-6.972-8.268-6.972-3.752 0-8.4 1.428-11.455 3.685C2.59 6.937 3.885 9.958 4.35 9.626c2.648-1.904 4.748-3.13 6.784-3.744C8.12 9.244.886 17.05 0 18.425c.1 1.261 1.66 4.648 2.424 4.648.232 0 .431-.133.664-.365a100.49 100.49 0 0 0 5.54-6.765c.222 3.104 1.748 6.898 6.014 6.898 3.819 0 7.604-2.756 9.33-8.965.2-.764-.73-1.361-1.261-.73zm-4.349-5.013c0 1.959-1.926 2.922-3.685 2.922-.941 0-1.664-.247-2.235-.568 1.051-1.592 2.092-3.225 3.21-4.973 1.972.334 2.71 1.43 2.71 2.619z'

const KAFKA_PATH =
  'M9.71 2.136a1.43 1.43 0 0 0-2.047 0h-.007a1.48 1.48 0 0 0-.421 1.042c0 .41.161.777.422 1.039l.007.007c.257.264.616.426 1.019.426.404 0 .766-.162 1.027-.426l.003-.007c.261-.262.421-.629.421-1.039 0-.408-.159-.777-.421-1.042H9.71zM8.683 22.295c.404 0 .766-.167 1.027-.429l.003-.008c.261-.261.421-.631.421-1.036 0-.41-.159-.778-.421-1.044H9.71a1.42 1.42 0 0 0-1.027-.432 1.4 1.4 0 0 0-1.02.432h-.007c-.26.266-.422.634-.422 1.044 0 .406.161.775.422 1.036l.007.008c.258.262.617.429 1.02.429zm7.89-4.462c.359-.096.683-.33.882-.684l.027-.052a1.47 1.47 0 0 0 .114-1.067 1.454 1.454 0 0 0-.675-.896l-.021-.014a1.425 1.425 0 0 0-1.078-.132c-.36.091-.684.335-.881.686-.2.349-.241.75-.146 1.119.099.363.33.691.675.896h.002c.346.203.737.239 1.101.144zm-6.405-7.342a2.083 2.083 0 0 0-1.485-.627c-.58 0-1.103.242-1.482.627-.378.385-.612.916-.612 1.507s.233 1.124.612 1.514a2.08 2.08 0 0 0 2.967 0c.379-.39.612-.923.612-1.514s-.233-1.122-.612-1.507zm-.835-2.51c.843.141 1.6.552 2.178 1.144h.004c.092.093.182.196.265.299l1.446-.851a3.176 3.176 0 0 1-.047-1.808 3.149 3.149 0 0 1 1.456-1.926l.025-.016a3.062 3.062 0 0 1 2.345-.306c.77.21 1.465.721 1.898 1.482v.002c.431.757.518 1.626.313 2.408a3.145 3.145 0 0 1-1.456 1.928l-.198.118h-.02a3.095 3.095 0 0 1-2.154.201 3.127 3.127 0 0 1-1.514-.944l-1.444.848a4.162 4.162 0 0 1 0 2.879l1.444.846c.413-.47.939-.789 1.514-.944a3.041 3.041 0 0 1 2.371.319l.048.023v.002a3.17 3.17 0 0 1 1.408 1.906 3.215 3.215 0 0 1-.313 2.405l-.026.053-.003-.005a3.147 3.147 0 0 1-1.867 1.436 3.096 3.096 0 0 1-2.371-.318v-.006a3.156 3.156 0 0 1-1.456-1.927 3.175 3.175 0 0 1 .047-1.805l-1.446-.848a3.905 3.905 0 0 1-.265.294l-.004.005a3.938 3.938 0 0 1-2.178 1.138v1.699a3.09 3.09 0 0 1 1.56.862l.002.004c.565.572.914 1.368.914 2.243 0 .873-.35 1.664-.914 2.239l-.002.009a3.1 3.1 0 0 1-2.21.931 3.1 3.1 0 0 1-2.206-.93h-.002v-.009a3.186 3.186 0 0 1-.916-2.239c0-.875.35-1.672.916-2.243v-.004h.002a3.1 3.1 0 0 1 1.558-.862v-1.699a3.926 3.926 0 0 1-2.176-1.138l-.006-.005a4.098 4.098 0 0 1-1.173-2.874c0-1.122.452-2.136 1.173-2.872h.006a3.947 3.947 0 0 1 2.176-1.144V6.289a3.137 3.137 0 0 1-1.558-.864h-.002v-.004a3.192 3.192 0 0 1-.916-2.243c0-.871.35-1.669.916-2.243l.002-.002A3.084 3.084 0 0 1 8.683 0c.861 0 1.641.355 2.21.932v.002h.002c.565.574.914 1.372.914 2.243 0 .876-.35 1.667-.914 2.243l-.002.005a3.142 3.142 0 0 1-1.56.864v1.692zm8.121-1.129l-.012-.019a1.452 1.452 0 0 0-.87-.668 1.43 1.43 0 0 0-1.103.146h.002c-.347.2-.58.529-.677.896-.095.365-.054.768.146 1.119l.007.009c.2.347.519.579.874.673.357.103.755.059 1.098-.144l.019-.009a1.47 1.47 0 0 0 .657-.885 1.493 1.493 0 0 0-.141-1.118'

const STRIPE_PATH =
  'M13.976 9.15c-2.172-.806-3.356-1.426-3.356-2.409 0-.831.683-1.305 1.901-1.305 2.227 0 4.515.858 6.09 1.631l.89-5.494C18.252.975 15.697 0 12.165 0 9.667 0 7.589.654 6.104 1.872 4.56 3.147 3.757 4.992 3.757 7.218c0 4.039 2.467 5.76 6.476 7.219 2.585.92 3.445 1.574 3.445 2.583 0 .98-.84 1.545-2.354 1.545-1.875 0-4.965-.921-6.99-2.109l-.9 5.555C5.175 22.99 8.385 24 11.714 24c2.641 0 4.843-.624 6.328-1.813 1.664-1.305 2.525-3.236 2.525-5.732 0-4.128-2.524-5.851-6.594-7.305h.003z'

const PROMETHEUS_PATH =
  'M12 0C5.373 0 0 5.372 0 12c0 6.627 5.373 12 12 12s12-5.373 12-12c0-6.628-5.373-12-12-12zm0 22.46c-1.885 0-3.414-1.26-3.414-2.814h6.828c0 1.553-1.528 2.813-3.414 2.813zm5.64-3.745H6.36v-2.046h11.28v2.046zm-.04-3.098H6.391c-.037-.043-.075-.086-.111-.13-1.155-1.401-1.427-2.133-1.69-2.879-.005-.025 1.4.287 2.395.511 0 0 .513.119 1.262.255-.72-.843-1.147-1.915-1.147-3.01 0-2.406 1.845-4.508 1.18-6.207.648.053 1.34 1.367 1.387 3.422.689-.951.977-2.69.977-3.755 0-1.103.727-2.385 1.454-2.429-.648 1.069.168 1.984.894 4.256.272.854.237 2.29.447 3.201.07-1.892.395-4.652 1.595-5.605-.529 1.2.079 2.702.494 3.424.671 1.164 1.078 2.047 1.078 3.716a4.642 4.642 0 01-1.11 2.996c.792-.149 1.34-.283 1.34-.283l2.573-.502s-.374 1.538-1.81 3.019z'

const GRAFANA_PATH =
  'M23.02 10.59a8.578 8.578 0 0 0-.862-3.034 8.911 8.911 0 0 0-1.789-2.445c.337-1.342-.413-2.505-.413-2.505-1.292-.08-2.113.4-2.416.62-.052-.02-.102-.044-.154-.064-.22-.089-.446-.172-.677-.247-.231-.073-.47-.14-.711-.197a9.867 9.867 0 0 0-.875-.161C14.557.753 12.94 0 12.94 0c-1.804 1.145-2.147 2.744-2.147 2.744l-.018.093c-.098.029-.2.057-.298.088-.138.042-.275.094-.413.143-.138.055-.275.107-.41.166a8.869 8.869 0 0 0-1.557.87l-.063-.029c-2.497-.955-4.716.195-4.716.195-.203 2.658.996 4.33 1.235 4.636a11.608 11.608 0 0 0-.607 2.635C1.636 12.677.953 15.014.953 15.014c1.926 2.214 4.171 2.351 4.171 2.351.003-.002.006-.002.006-.005.285.509.615.994.986 1.446.156.19.32.371.488.548-.704 2.009.099 3.68.099 3.68 2.144.08 3.553-.937 3.849-1.173a9.784 9.784 0 0 0 3.164.501h.08l.055-.003.107-.002.103-.005.003.002c1.01 1.44 2.788 1.646 2.788 1.646 1.264-1.332 1.337-2.653 1.337-2.94v-.058c0-.02-.003-.039-.003-.06.265-.187.52-.387.758-.6a7.875 7.875 0 0 0 1.415-1.7c1.43.083 2.437-.885 2.437-.885-.236-1.49-1.085-2.216-1.264-2.354l-.018-.013-.016-.013a.217.217 0 0 1-.031-.02c.008-.092.016-.18.02-.27.011-.162.016-.323.016-.48v-.253l-.005-.098-.008-.135a1.891 1.891 0 0 0-.01-.13c-.003-.042-.008-.083-.013-.125l-.016-.124-.018-.122a6.215 6.215 0 0 0-2.032-3.73 6.015 6.015 0 0 0-3.222-1.46 6.292 6.292 0 0 0-.85-.048l-.107.002h-.063l-.044.003-.104.008a4.777 4.777 0 0 0-3.335 1.695c-.332.4-.592.84-.768 1.297a4.594 4.594 0 0 0-.312 1.817l.003.091c.005.055.007.11.013.164a3.615 3.615 0 0 0 .698 1.82 3.53 3.53 0 0 0 1.827 1.282c.33.098.66.14.971.137.039 0 .078 0 .114-.002l.063-.003c.02 0 .041-.003.062-.003.034-.002.065-.007.099-.01.007 0 .018-.003.028-.003l.031-.005.06-.008a1.18 1.18 0 0 0 .112-.02c.036-.008.072-.013.109-.024a2.634 2.634 0 0 0 .914-.415c.028-.02.056-.041.085-.065a.248.248 0 0 0 .039-.35.244.244 0 0 0-.309-.06l-.078.042c-.09.044-.184.083-.283.116a2.476 2.476 0 0 1-.475.096c-.028.003-.054.006-.083.006l-.083.002c-.026 0-.054 0-.08-.002l-.102-.006h-.012l-.024.006c-.016-.003-.031-.003-.044-.006-.031-.002-.06-.007-.091-.01a2.59 2.59 0 0 1-.724-.213 2.557 2.557 0 0 1-.667-.438 2.52 2.52 0 0 1-.805-1.475 2.306 2.306 0 0 1-.029-.444l.006-.122v-.023l.002-.031c.003-.021.003-.04.005-.06a3.163 3.163 0 0 1 1.352-2.29 3.12 3.12 0 0 1 .937-.43 2.946 2.946 0 0 1 .776-.101h.06l.07.002.045.003h.026l.07.005a4.041 4.041 0 0 1 1.635.49 3.94 3.94 0 0 1 1.602 1.662 3.77 3.77 0 0 1 .397 1.414l.005.076.003.075c.002.026.002.05.002.075 0 .024.003.052 0 .07v.065l-.002.073-.008.174a6.195 6.195 0 0 1-.08.639 5.1 5.1 0 0 1-.267.927 5.31 5.31 0 0 1-.624 1.13 5.052 5.052 0 0 1-3.237 2.014 4.82 4.82 0 0 1-.649.066l-.039.003h-.287a6.607 6.607 0 0 1-1.716-.265 6.776 6.776 0 0 1-3.4-2.274 6.75 6.75 0 0 1-.746-1.15 6.616 6.616 0 0 1-.714-2.596l-.005-.083-.002-.02v-.056l-.003-.073v-.096l-.003-.104v-.07l.003-.163c.008-.22.026-.45.054-.678a8.707 8.707 0 0 1 .28-1.355c.128-.444.286-.872.473-1.277a7.04 7.04 0 0 1 1.456-2.1 5.925 5.925 0 0 1 .953-.763c.169-.111.343-.213.524-.306.089-.05.182-.091.273-.135.047-.02.093-.042.138-.062a7.177 7.177 0 0 1 .714-.267l.145-.045c.049-.015.098-.026.148-.041.098-.029.197-.052.296-.076.049-.013.1-.02.15-.033l.15-.032.151-.028.076-.013.075-.01.153-.024c.057-.01.114-.013.171-.023l.169-.021c.036-.003.073-.008.106-.01l.073-.008.036-.003.042-.002c.057-.003.114-.008.171-.01l.086-.006h.023l.037-.003.145-.007a7.999 7.999 0 0 1 1.708.125 7.917 7.917 0 0 1 2.048.68 8.253 8.253 0 0 1 1.672 1.09l.09.077.089.078c.06.052.114.107.171.159.057.052.112.106.166.16.052.055.107.107.159.164a8.671 8.671 0 0 1 1.41 1.978c.012.026.028.052.04.078l.04.078.075.156c.023.051.05.1.07.153l.065.15a8.848 8.848 0 0 1 .45 1.34.19.19 0 0 0 .201.142.186.186 0 0 0 .172-.184c.01-.246.002-.532-.024-.856z'

const JAEGER_PATH =
  'M14.816 22.3774c0 .0724-.1283.1311-.2865.1311-.1581 0-.2864-.0587-.2864-.131 0-.0725.1283-.1311.2864-.1311.1582 0 .2865.0586.2865.131Zm-.738.1554c0 .075-.126.136-.2815.136-.1555 0-.2816-.061-.2816-.136 0-.075.126-.136.2816-.136.1555 0 .2816.061.2816.136Zm-.7427.1408c0 .0697-.1739.1262-.3884.1262-.2145 0-.3884-.0565-.3884-.1262 0-.0697.174-.1263.3884-.1263.2145 0 .3884.0566.3884.1263Zm-1.136.1116c0 .0885-.1587.1602-.3544.1602-.1957 0-.3544-.0717-.3544-.1602 0-.0885.1587-.1602.3544-.1602.1957 0 .3544.0717.3544.1602Zm-1.1943.0389c0 .0965-.2369.1747-.5291.1747-.2923 0-.5292-.0782-.5292-.1747 0-.0966.237-.1748.5292-.1748s.5291.0782.5291.1748Zm1.107-.8787c.6116-.0292.7718-.1214.7912-.3253 0 0 .1214-.1505-.2913-.233-.4126-.0826-.9952-.0826-.8932-.4127.102-.335 1.1554-.3932 1.7865-.3932s1.3447-.0388 1.2525.301c-.0874.3398-.5486.3447-.5097.5728.0388.2282.4175.1602.3495.5195-.068.3592-1.7428.4903-2.9274.5923-1.1797.102-1.9321-.0195-1.7428-.3399.1893-.3204 1.573-.2524 2.1846-.2815Zm-9.9693 1.0755c.0066.045-.128.102-.3004.127-.1724.0251-.3176.0089-.3241-.0362-.0066-.0452.128-.102.3004-.1271.1725-.025.3176-.0088.3241.0363Zm.0918.226c.0077.053-.1474.1195-.3464.1484-.199.029-.3665.0094-.3742-.0437-.0078-.053.1473-.1195.3463-.1485.199-.0289.3666-.0093.3743.0438Zm.1648.204c.0077.053-.1473.1194-.3463.1484-.199.0289-.3666.0093-.3743-.0437-.0077-.0531.1473-.1196.3463-.1485.199-.029.3666-.0094.3743.0437Zm.7387.114c.0086.0638-.1501.138-.3547.166-.2046.0278-.3774-.0013-.3861-.065-.0087-.0638.15-.138.3546-.166.2046-.0278.3775.0013.3862.065Zm1.1803.1309c.0066.0828-.2394.17-.5494.1947s-.5667-.0225-.5733-.1053c-.0066-.0829.2394-.17.5494-.1947s.5667.0224.5733.1053Zm1.1926.0011c.0182.0838-.2006.2026-.4888.2653-.2882.0627-.5366.0456-.5548-.0382-.0182-.0839.2006-.2027.4888-.2654.2882-.0627.5366-.0455.5548.0383Zm1.2463-2.0331c.3544-.0777 1.9856-.0388 2.0292.2476.0534.2913-2.67.5437-3.2866.631-.6165.0875-.6408.2622-.6408.2622 0 .0146-.0534.1554.3495.1505.5632-.0097.9564.1068.7234.3787-.233.2719-2.8109.199-3.34 0-.5292-.199-.3204-.398-.0486-.5.2719-.0971 1.0438-.1894 1.937-.3593.8933-.1699 1.9274-.7282 2.2769-.8107Zm6.7917-9.7337a.6457.6457 0 1 1-1.2914 0 .6457.6457 0 0 1 1.2914 0ZM7.3155 13.44a.6457.6457 0 1 1-1.2914 0 .6457.6457 0 0 1 1.2914 0ZM4.325 20.023c-.3544 0-.602-.1409-.7573-.4176-.0631-.1068-.102-.2379-.1311-.3738-.0631-.335-.0437-.7282.0291-1.0972.0825-.4126.2282-.7961.3835-1.0389.1408-.2233.2719-.369.3884-.466a50.61 50.61 0 0 0 .5195 1.3205c.2378.5874.4563 1.1068.6165 1.4758-.0097.0194-.0243.034-.0388.0534-.233.3107-.5583.5437-1.0098.5437Zm15.0155-2.1944c0 .0146.0049.0388.0195.0728.0534.2913.233 1.0778.3543 1.5292-1.8642.2088-4.5051.3399-7.4325.3399-2.4759 0-4.7479-.0923-6.515-.2476a86.46 86.46 0 0 1-.6748-1.607c-.597-1.4563-1.3302-3.3933-1.5583-4.5779-.2088-1.0923-.4612-2.4176.1165-3.6119A28.5682 28.5682 0 0 1 3.32 9.697c-1.369-.1456-1.5535-.4127-1.6263-.5243-.0777-.1117-.0923-.2427-.0389-.3738.1214-.3059.534-.5826 1.0146-.8108a5.028 5.028 0 0 0 .136.2816c-.437.2088-.7768.432-.8593.6457-.0194.0437-.0097.0631.0049.0874.068.097.4611.398 2.6846.4709a.6316.6316 0 0 1-.5-.505.634.634 0 0 1 .0582-.3931c-.2039.1068-.4078.233-.602.364 0 0-1.6554-2.1748-.4709-2.8933.136-.0826.267-.1214.3981-.1214.9904 0 1.7234 2.1603 1.7234 2.1603s-.1456.0534-.3592.1505c.2621.0243.7816.1457 1.8302.568 1.068.4321 1.4952.4612 1.6118.4612A.625.625 0 0 1 9 9.3426l.034.034a22.885 22.885 0 0 0 .4563-.0389c.0485-.0048.1116-.0194.1942-.0388a.6277.6277 0 0 1 .2864-.5826.6193.6193 0 0 1 .733.0437.434.434 0 0 0 .0874-.034c.2719-.1359 1.2331-1.1602 1.9565-2.034a.616.616 0 0 1 .6117-.2137c.1941-.6553.9709-2.9759 2.1554-2.9759 1.3788 0 .5777 2.7283.4418 2.9225 0 0-.3204-.0242-.7573-.0242-.4855 0-1.1166.029-1.6312.165.0146.0097.034.0243.0534.0389.267.2184.301.6117.0825.8787-.1165.1359-.267.3204-.4417.5194.3592-.1456.733-.3058 1.1068-.466.6263-.2768 1.2089-.539 1.7234-.7914h.0098l.0048-.0097c.4418-.2136.8399-.4175 1.1748-.6068 1.17-.6554 1.3157-.9418 1.3302-1.0146.0097-.0486 0-.0631-.0048-.068-.0097-.0194-.0923-.097-.4612-.097-.403 0-.9661.0922-1.5535.2184a3.849 3.849 0 0 0 .0291-.3253c.5729-.1165 1.1214-.2039 1.5292-.2039.369 0 .5923.068.7088.2136.0534.0631.1068.1748.068.335-.0437.1893-.2428.5194-1.4807 1.2137a8.6623 8.6623 0 0 1-.3738.2039c-.0922.0534-.1942.102-.2961.1553.5582.471.9758 1.1748 1.2379 2.1943.5292 2.039.937 4.4275 1.2185 6.3208-.4612.0583-.7427.2767-.801.6311-.0631.3836.136.5923.3447.8156.2233.2379.4806.505.5631 1.1069Zm-4.84-12.331c0 .0244.0048.0195.0048.0195s.0242-.0437.063-.1165c.0438-.0728.102-.1748.17-.2961.0728-.1214.1553-.2573.2573-.3787.0485-.0583.1068-.1165.165-.1602.0632-.0437.1214-.0728.1845-.0825.0194 0 .034-.0049.0486-.0049.0145.0049.034 0 .0485.0049.0146.0048.034.0048.0486.0097l.0388.0097.0097.0048c0 .0049.0049.0049.0049.0097.0097.0146.0145.0437.0194.0777 0 .068-.0049.1505-.0194.2233-.0097.0729-.0292.1457-.0437.2136-.034.136-.0631.2476-.0826.3302-.0242.0825-.034.131-.034.131s.034-.034.0874-.102c.0534-.063.1214-.165.1942-.2912.034-.0631.0729-.136.102-.2185.0291-.0825.0582-.1747.068-.2815.0048-.0534.0048-.1165-.0195-.1845-.0048-.0194-.0145-.034-.0242-.0534l-.0146-.0243-.0194-.0243c-.0194-.029-.0971-.0776-.102-.0776l-.0097-.0049h-.0048l-.0146-.0048h-.0049l-.0194-.0097-.0437-.0146c-.0242-.0097-.0582-.0146-.0874-.0194-.029-.0049-.063-.0049-.0922-.0049-.0291.0049-.0631.0049-.0922.0097a.6844.6844 0 0 0-.3204.1554c-.0874.0728-.1505.1553-.204.2379-.1019.165-.1601.33-.199.4709-.034.1456-.0534.267-.0582.3495-.0049.0388-.0049.0728-.0049.097Zm-4.2868 3.578c-.0874.0534-.1116.165-.063.2524.0145.0195.0922.1311.2572.1942a.707.707 0 0 0 .2525.0437c.1456 0 .3058-.0388.4709-.1116.3592-.1602 1.2379-1.0487 2.2428-2.2575.0631-.0776.0534-.1942-.0242-.2573-.0437-.034-.0874-.0437-.1166-.0437a.1794.1794 0 0 0-.1407.068c-.4321.5195-1.6555 1.9225-2.0924 2.1458l-.0194.0097c-.0534.0194-.102.0388-.1505.0534l-.2233.0631-.1797-.1456c-.0291-.0292-.0728-.0437-.1165-.0437a.175.175 0 0 0-.097.029ZM9.17 14.304c.0582.335.5874.5243 1.1845.4272.5971-.097 1.039-.4515.9807-.7865-.0583-.335-.5875-.5243-1.1846-.4272-.5971.102-1.034.4564-.9806.7865ZM4.4172 7.9493a2.001 2.001 0 0 0 .0729.0922s-.0049-.0145-.0049-.034c-.0049-.0194-.0049-.0485-.0146-.0873-.0097-.0777-.0388-.1845-.0776-.3107a2.4088 2.4088 0 0 0-.1748-.4175 1.8868 1.8868 0 0 0-.1408-.2233c-.0534-.0729-.1165-.1457-.2039-.2136-.0485-.0292-.097-.0632-.165-.0826-.0146-.0048-.034-.0048-.0535-.0097h-.0728c-.0097 0-.0291.0049-.034.0049l-.029.0097-.0195.0048-.0049.0049-.0097.0048a.3371.3371 0 0 0-.1602.136.4473.4473 0 0 0-.0534.1747c-.0097.1068.0146.1942.034.272.0243.0776.0534.1504.0825.2087a1.493 1.493 0 0 0 .17.2718c.0242.0291.0436.0534.0582.068l.0243.0243s-.0097-.0437-.0243-.1166c-.0146-.0728-.0388-.1796-.0582-.301a2.2904 2.2904 0 0 1-.0243-.1941c-.0049-.068-.0049-.136.0097-.1894.0048-.0242.0146-.0437.0243-.0534.0097-.0048.0145-.0097.029-.0145h.0098l.0146-.0049c.0048 0 0 .0049 0 .0049h.0048a.2045.2045 0 0 1 .0631.034c.0486.0388.102.0922.1505.1504.0534.0632.102.1214.1457.1845.0873.1262.1747.2476.2476.3496a5.199 5.199 0 0 0 .1844.2524Zm.3544 1.0777.034.0049c.131.0146.5534.0922 1.5923.5097.8545.3447 1.5001.5195 1.9128.5195.3107 0 .3932-.0971.4223-.131.0437-.0535.0486-.1069.0437-.136-.0048-.0291-.0145-.0825-.068-.1262-.0436-.034-.0873-.0437-.1164-.0437-.0292 0-.0583.0048-.0826.0194l-.0874.0437h-.1165c-.2476 0-.7622-.0825-1.767-.4903-1.204-.4855-1.6458-.539-1.7769-.539-.0291 0-.0388 0-.0582.005-.0632.0145-.102.0533-.1166.0776-.0145.0194-.0388.0728-.029.136.0145.0873.0922.1504.1795.1504h.034Zm.0437 3.6945c0 1.1262.9127 2.0438 2.039 2.039 1.1263 0 2.039-.9128 2.039-2.039 0-1.1263-.9127-2.039-2.039-2.039s-2.039.9127-2.039 2.039Zm5.5246 3.3643c-.267.102-.5.2475-.7039.3738.0486.2184.199.7088.6311.5825.3544-.1068.1505-.733.0728-.9563Zm.6214.7816c.1166.0437.5778.0825.534-.4758-.0145-.1942-.029-.3301-.0388-.4224-.0437-.0048-.0825-.0048-.1213-.0048-.136 0-.2622-.0049-.369-.0049-.1117 0-.17.0049-.1942.0097a9.3046 9.3046 0 0 0-.1553.0292c.0388.2524.1456.7961.3446.869Zm1.7866-1.3788c.432-.1942.432-.5243.4078-.6942-.0292-.1796-.2962-.5583-1.3788-.568 0 0-.335.7816-1.2428.9273-.9078.1505-1.2962-.0243-1.4515-.2185 0 0-.8739.2573-.9467.8107-.0728.5535.6603.7477 1.0098.5535.3495-.1942.8738-.6166 1.573-.7088.3737-.0486 1.5971.0922 2.0292-.102Zm2.471-4.2333c0-1.1262-.9127-2.039-2.039-2.039s-2.039.9128-2.039 2.039c0 1.1263.9127 2.039 2.039 2.039s2.039-.9175 2.039-2.039ZM3.1162 18.1975c-.0437.267-.0728.6117-.0243.9564-.6602-.1408-1.034-.3058-1.034-.4757 0-.17.3835-.335 1.0583-.4807Zm18.4866.0389c.5729.136.8884.2864.8836.4418 0 .2476-.767.4757-2.0633.6553.233-.063.4952-.233.7962-.5146.1505-.1407.267-.3107.3495-.5a.5922.5922 0 0 0 .034-.0825Zm-.3495-.0728a1.0453 1.0453 0 0 1-.2767.3883c-.3884.369-.6214.4467-.7525.4467-.2573 0-.4466-.3593-.5583-1.0632-.0049-.0486-.0146-.0971-.0194-.1456-.0923-.6991-.4127-1.039-.6457-1.2817-.2039-.2136-.301-.3301-.2621-.5486.0534-.3155.4369-.3835.7476-.3835.1456 0 .2573.0146.2573.0146h.0194c.5728.0097 1.17.5291 1.4515 1.2525.1554.398.2185.903.0389 1.3205ZM4.7085 6.3618c-.1068-.1456-.2136-.267-.3253-.3592-.0097-.3981-.0097-.7574.0098-.8981.0388-.2622.6893-1.0244 1.2962-1.607.5097-.4903 1.4855-1.335 2.301-1.5.1263-.0243.272-.0389.4419-.0389 1.5535 0 4.6993 1.1651 4.8546 1.3205l.0097.0097c.0389-.3787.0777-.9078.0098-1.0583 0 0 .296.0922.3495.2767 0 0 .733-.738.6602-1.6554 0 0 .165 0 .1845.296 0 0 .5292-1.2524.9175-1.1408.3884.1117-.1262 1.8594-.4029 2.1895 0 0 .4418-.1068.4952-.403 0 0 .0534.2574-.0194.403 0 0 .9952-.8107.9563-1.1214 0 0 .1311.3155.0923.5534 0 0 .8253-.6602.8253-.937 0 0 .1262 1.6021-1.471 2.1895 0 0 .4224-.0534.6263-.0922 0 0-.1748.2281-.4418.5534-.1602-.102-.3496-.1505-.5632-.1505-.6408 0-1.2282.4855-1.7476 1.4419.0776-.4175-.0292-.7574-.1845-.7914-.0049-.0048-.0146-.0048-.0243-.0048-.1602 0-.3884.2767-.4806.6894-.0923.4369.0146.8058.1748.8398.0582.0146.131-.0194.2039-.0874-.1505.3593-.267.6942-.3447.9321-3.976 1.7914-6.5053 1.7527-7.656 1.573-.0825-.2137-.233-.5778-.4417-.9419-.0583-.097-.1166-.1941-.1748-.2815 2.7769.7476 6.4761-.6506 7.9568-1.2962-.0777-.2185-.0874-.5-.0194-.806.0631-.296.1942-.5485.3544-.7184-.0291-.1165-.0534-.199-.0728-.2427-.3399-.1845-3.1556-1.2331-4.6265-1.2331-.1457 0-.272.0145-.3787.034-1.301.267-3.277 2.5292-3.3497 2.8497-.0243.1553-.0146.6796.0048 1.2136Z'

const DOCKER_PATH =
  'M13.983 11.078h2.119a.186.186 0 00.186-.185V9.006a.186.186 0 00-.186-.186h-2.119a.185.185 0 00-.185.185v1.888c0 .102.083.185.185.185m-2.954-5.43h2.118a.186.186 0 00.186-.186V3.574a.186.186 0 00-.186-.185h-2.118a.185.185 0 00-.185.185v1.888c0 .102.082.185.185.185m0 2.716h2.118a.187.187 0 00.186-.186V6.29a.186.186 0 00-.186-.185h-2.118a.185.185 0 00-.185.185v1.887c0 .102.082.185.185.186m-2.93 0h2.12a.186.186 0 00.184-.186V6.29a.185.185 0 00-.185-.185H8.1a.185.185 0 00-.185.185v1.887c0 .102.083.185.185.186m-2.964 0h2.119a.186.186 0 00.185-.186V6.29a.185.185 0 00-.185-.185H5.136a.186.186 0 00-.186.185v1.887c0 .102.084.185.186.186m5.893 2.715h2.118a.186.186 0 00.186-.185V9.006a.186.186 0 00-.186-.186h-2.118a.185.185 0 00-.185.185v1.888c0 .102.082.185.185.185m-2.93 0h2.12a.185.185 0 00.184-.185V9.006a.185.185 0 00-.184-.186h-2.12a.185.185 0 00-.184.185v1.888c0 .102.083.185.185.185m-2.964 0h2.119a.185.185 0 00.185-.185V9.006a.185.185 0 00-.184-.186h-2.12a.186.186 0 00-.186.186v1.887c0 .102.084.185.186.185m-2.92 0h2.12a.185.185 0 00.184-.185V9.006a.185.185 0 00-.184-.186h-2.12a.185.185 0 00-.184.185v1.888c0 .102.082.185.185.185M23.763 9.89c-.065-.051-.672-.51-1.954-.51-.338.001-.676.03-1.01.087-.248-1.7-1.653-2.53-1.716-2.566l-.344-.199-.226.327c-.284.438-.49.922-.612 1.43-.23.97-.09 1.882.403 2.661-.595.332-1.55.413-1.744.42H.751a.751.751 0 00-.75.748 11.376 11.376 0 00.692 4.062c.545 1.428 1.355 2.48 2.41 3.124 1.18.723 3.1 1.137 5.275 1.137.983.003 1.963-.086 2.93-.266a12.248 12.248 0 003.823-1.389c.98-.567 1.86-1.288 2.61-2.136 1.252-1.418 1.998-2.997 2.553-4.4h.221c1.372 0 2.215-.549 2.68-1.009.309-.293.55-.65.707-1.046l.098-.288Z'

const KUBERNETES_PATH =
  'M10.204 14.35l.007.01-.999 2.413a5.171 5.171 0 0 1-2.075-2.597l2.578-.437.004.005a.44.44 0 0 1 .484.606zm-.833-2.129a.44.44 0 0 0 .173-.756l.002-.011L7.585 9.7a5.143 5.143 0 0 0-.73 3.255l2.514-.725.002-.009zm1.145-1.98a.44.44 0 0 0 .699-.337l.01-.005.15-2.62a5.144 5.144 0 0 0-3.01 1.442l2.147 1.523.004-.002zm.76 2.75l.723.349.722-.347.18-.78-.5-.623h-.804l-.5.623.179.779zm1.5-3.095a.44.44 0 0 0 .7.336l.008.003 2.134-1.513a5.188 5.188 0 0 0-2.992-1.442l.148 2.615.002.001zm10.876 5.97l-5.773 7.181a1.6 1.6 0 0 1-1.248.594l-9.261.003a1.6 1.6 0 0 1-1.247-.596l-5.776-7.18a1.583 1.583 0 0 1-.307-1.34L2.1 5.573c.108-.47.425-.864.863-1.073L11.305.513a1.606 1.606 0 0 1 1.385 0l8.345 3.985c.438.209.755.604.863 1.073l2.062 8.955c.108.47-.005.963-.308 1.34zm-3.289-2.057c-.042-.01-.103-.026-.145-.034-.174-.033-.315-.025-.479-.038-.35-.037-.638-.067-.895-.148-.105-.04-.18-.165-.216-.216l-.201-.059a6.45 6.45 0 0 0-.105-2.332 6.465 6.465 0 0 0-.936-2.163c.052-.047.15-.133.177-.159.008-.09.001-.183.094-.282.197-.185.444-.338.743-.522.142-.084.273-.137.415-.242.032-.024.076-.062.11-.089.24-.191.295-.52.123-.736-.172-.216-.506-.236-.745-.045-.034.027-.08.062-.111.088-.134.116-.217.23-.33.35-.246.25-.45.458-.673.609-.097.056-.239.037-.303.033l-.19.135a6.545 6.545 0 0 0-4.146-2.003l-.012-.223c-.065-.062-.143-.115-.163-.25-.022-.268.015-.557.057-.905.023-.163.061-.298.068-.475.001-.04-.001-.099-.001-.142 0-.306-.224-.555-.5-.555-.275 0-.499.249-.499.555l.001.014c0 .041-.002.092 0 .128.006.177.044.312.067.475.042.348.078.637.056.906a.545.545 0 0 1-.162.258l-.012.211a6.424 6.424 0 0 0-4.166 2.003 8.373 8.373 0 0 1-.18-.128c-.09.012-.18.04-.297-.029-.223-.15-.427-.358-.673-.608-.113-.12-.195-.234-.329-.349-.03-.026-.077-.062-.111-.088a.594.594 0 0 0-.348-.132.481.481 0 0 0-.398.176c-.172.216-.117.546.123.737l.007.005.104.083c.142.105.272.159.414.242.299.185.546.338.743.522.076.082.09.226.1.288l.16.143a6.462 6.462 0 0 0-1.02 4.506l-.208.06c-.055.072-.133.184-.215.217-.257.081-.546.11-.895.147-.164.014-.305.006-.48.039-.037.007-.09.02-.133.03l-.004.002-.007.002c-.295.071-.484.342-.423.608.061.267.349.429.645.365l.007-.001.01-.003.129-.029c.17-.046.294-.113.448-.172.33-.118.604-.217.87-.256.112-.009.23.069.288.101l.217-.037a6.5 6.5 0 0 0 2.88 3.596l-.09.218c.033.084.069.199.044.282-.097.252-.263.517-.452.813-.091.136-.185.242-.268.399-.02.037-.045.095-.064.134-.128.275-.034.591.213.71.248.12.556-.007.69-.282v-.002c.02-.039.046-.09.062-.127.07-.162.094-.301.144-.458.132-.332.205-.68.387-.897.05-.06.13-.082.215-.105l.113-.205a6.453 6.453 0 0 0 4.609.012l.106.192c.086.028.18.042.256.155.136.232.229.507.342.84.05.156.074.295.145.457.016.037.043.09.062.129.133.276.442.402.69.282.247-.118.341-.435.213-.71-.02-.039-.045-.096-.065-.134-.083-.156-.177-.261-.268-.398-.19-.296-.346-.541-.443-.793-.04-.13.007-.21.038-.294-.018-.022-.059-.144-.083-.202a6.499 6.499 0 0 0 2.88-3.622c.064.01.176.03.213.038.075-.05.144-.114.28-.104.266.039.54.138.87.256.154.06.277.128.448.173.036.01.088.019.13.028l.009.003.007.001c.297.064.584-.098.645-.365.06-.266-.128-.537-.423-.608zM16.4 9.701l-1.95 1.746v.005a.44.44 0 0 0 .173.757l.003.01 2.526.728a5.199 5.199 0 0 0-.108-1.674A5.208 5.208 0 0 0 16.4 9.7zm-4.013 5.325a.437.437 0 0 0-.404-.232.44.44 0 0 0-.372.233h-.002l-1.268 2.292a5.164 5.164 0 0 0 3.326.003l-1.27-2.296h-.01zm1.888-1.293a.44.44 0 0 0-.27.036.44.44 0 0 0-.214.572l-.003.004 1.01 2.438a5.15 5.15 0 0 0 2.081-2.615l-2.6-.44-.004.005z'

const GITHUB_ACTIONS_PATH =
  'M10.984 13.836a.5.5 0 0 1-.353-.146l-.745-.743a.5.5 0 1 1 .706-.708l.392.391 1.181-1.18a.5.5 0 0 1 .708.707l-1.535 1.533a.504.504 0 0 1-.354.146zm9.353-.147l1.534-1.532a.5.5 0 0 0-.707-.707l-1.181 1.18-.392-.391a.5.5 0 1 0-.706.708l.746.743a.497.497 0 0 0 .706-.001zM4.527 7.452l2.557-1.585A1 1 0 0 0 7.09 4.17L4.533 2.56A1 1 0 0 0 3 3.406v3.196a1.001 1.001 0 0 0 1.527.85zm2.03-2.436L4 6.602V3.406l2.557 1.61zM24 12.5c0 1.93-1.57 3.5-3.5 3.5a3.503 3.503 0 0 1-3.46-3h-2.08a3.503 3.503 0 0 1-3.46 3 3.502 3.502 0 0 1-3.46-3h-.558c-.972 0-1.85-.399-2.482-1.042V17c0 1.654 1.346 3 3 3h.04c.244-1.693 1.7-3 3.46-3 1.93 0 3.5 1.57 3.5 3.5S13.43 24 11.5 24a3.502 3.502 0 0 1-3.46-3H8c-2.206 0-4-1.794-4-4V9.899A5.008 5.008 0 0 1 0 5c0-2.757 2.243-5 5-5s5 2.243 5 5a5.005 5.005 0 0 1-4.952 4.998A2.482 2.482 0 0 0 7.482 12h.558c.244-1.693 1.7-3 3.46-3a3.502 3.502 0 0 1 3.46 3h2.08a3.503 3.503 0 0 1 3.46-3c1.93 0 3.5 1.57 3.5 3.5zm-15 8c0 1.378 1.122 2.5 2.5 2.5s2.5-1.122 2.5-2.5-1.122-2.5-2.5-2.5S9 19.122 9 20.5zM5 9c2.206 0 4-1.794 4-4S7.206 1 5 1 1 2.794 1 5s1.794 4 4 4zm9 3.5c0-1.378-1.122-2.5-2.5-2.5S9 11.122 9 12.5s1.122 2.5 2.5 2.5 2.5-1.122 2.5-2.5zm9 0c0-1.378-1.122-2.5-2.5-2.5S18 11.122 18 12.5s1.122 2.5 2.5 2.5 2.5-1.122 2.5-2.5zm-13 8a.5.5 0 1 0 1 0 .5.5 0 0 0-1 0zm2 0a.5.5 0 1 0 1 0 .5.5 0 0 0-1 0zm12 0c0 1.93-1.57 3.5-3.5 3.5a3.503 3.503 0 0 1-3.46-3.002c-.007.001-.013.005-.021.005l-.506.017h-.017a.5.5 0 0 1-.016-.999l.506-.017c.018-.002.035.006.052.007A3.503 3.503 0 0 1 20.5 17c1.93 0 3.5 1.57 3.5 3.5zm-1 0c0-1.378-1.122-2.5-2.5-2.5S18 19.122 18 20.5s1.122 2.5 2.5 2.5 2.5-1.122 2.5-2.5z'

// each entry is either a real brand mark (brandPath + its official hex) or,
// for the two non-vendor items, one of the site's own line icons tinted
// with the structural "steel" color so it reads as pattern, not product
const stackGroups = [
  {
    label: 'security',
    items: [
      { name: 'JWT', blurb: 'stateless session tokens', hex: '#0B0E14', brandPath: JWT_PATH },
      { name: 'RBAC', blurb: 'route-level permission checks', hex: '#35455C', iconName: 'shield' },
    ],
  },
  {
    label: 'data & transactions',
    items: [
      { name: 'Redis', blurb: 'shared cache layer', hex: '#FF4438', brandPath: REDIS_PATH },
      { name: 'Kafka', blurb: 'order fan-out & event bus', hex: '#231F20', brandPath: KAFKA_PATH },
      { name: 'Stripe', blurb: 'checkout charges & refunds', hex: '#635BFF', brandPath: STRIPE_PATH },
    ],
  },
  {
    label: 'observability',
    items: [
      { name: 'Prometheus', blurb: 'service-level metrics', hex: '#E6522C', brandPath: PROMETHEUS_PATH },
      { name: 'Grafana', blurb: 'live ops dashboards', hex: '#F46800', brandPath: GRAFANA_PATH },
      { name: 'Jaeger', blurb: 'distributed request tracing', hex: '#66CFE3', brandPath: JAEGER_PATH },
    ],
  },
  {
    label: 'platform & delivery',
    items: [
      { name: 'Docker', blurb: 'every service, containerized', hex: '#2496ED', brandPath: DOCKER_PATH },
      { name: 'Kubernetes', blurb: 'orchestration & scaling', hex: '#326CE5', brandPath: KUBERNETES_PATH },
      { name: 'GitHub Actions', blurb: 'CI on every push', hex: '#2088FF', brandPath: GITHUB_ACTIONS_PATH },
      { name: 'Azure', blurb: 'cloud hosting & scaling', hex: '#35455C', iconName: 'cloud' },
    ],
  },
]

function StackBadge({ name, blurb, hex, brandPath, iconName }) {
  return (
    <div className="flex w-[168px] flex-none flex-col gap-3 rounded-2xl border border-slate-200 bg-white p-4 shadow-[0_1px_2px_rgba(15,23,42,0.04)] transition-all hover:-translate-y-0.5 hover:border-slate-300 hover:shadow-[0_12px_24px_-12px_rgba(15,23,42,0.18)] sm:w-[184px]">
      <div
        className="flex h-10 w-10 items-center justify-center rounded-xl"
        style={{ backgroundColor: `${hex}14`, color: hex }}
      >
        {brandPath ? (
          <svg viewBox="0 0 24 24" className="h-[18px] w-[18px]" fill="currentColor" aria-hidden="true">
            <path d={brandPath} />
          </svg>
        ) : (
          NodeIcon(iconName, 'h-[18px] w-[18px]')
        )}
      </div>
      <div>
        <p className="FleetOps-display text-[14px] font-semibold text-[#0B0E14]">{name}</p>
        <p className="mt-0.5 FleetOps-mono text-[11px] leading-relaxed text-[#5B6472]">{blurb}</p>
      </div>
    </div>
  )
}

function EngineeringStack() {
  return (
    <section id="stack" className="bg-[#F3F5F8] py-24">
      <div className="mx-auto max-w-6xl px-6">
        <div className="max-w-lg">
          <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">engineering highlights</span>
          <h2 className="mt-3 FleetOps-display text-3xl font-semibold tracking-tight text-[#0B0E14] sm:text-4xl">
            Not a prototype — this is what&rsquo;s running.
          </h2>
          <p className="mt-3 FleetOps-body text-[15px] leading-relaxed text-[#5B6472]">
            Every mark below is wired into the request flow above. Nothing here is on a
            roadmap, it&rsquo;s in the diagram.
          </p>
        </div>

        <div className="mt-12 space-y-10">
          {stackGroups.map((group) => (
            <div key={group.label}>
              <p className="FleetOps-mono text-[11px] uppercase tracking-wide text-[#5B6472]">{group.label}</p>
              <div className="mt-4 flex flex-wrap gap-3">
                {group.items.map((item) => (
                  <StackBadge key={item.name} {...item} />
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

// ---------------------------------------------------------------------------
// CTA + Footer
// ---------------------------------------------------------------------------

function CTA() {
  return (
    <section id="cta" className="bg-[#0B0E14] py-24">
      <div className="mx-auto flex max-w-6xl flex-col items-start justify-between gap-8 px-6 sm:flex-row sm:items-center">
        <div className="max-w-md">
          <h2 className="FleetOps-display text-3xl font-semibold tracking-tight text-white sm:text-4xl">
            Building the next fleet platform?
          </h2>
          <p className="mt-3 FleetOps-body text-[15px] leading-relaxed text-white/60">
            We&rsquo;re onboarding a small number of teams during the private beta.
          </p>
        </div>
        <a
          href="#"
          className="group inline-flex shrink-0 items-center gap-2 rounded-full bg-[#FF5A1F] px-7 py-4 FleetOps-body text-[14.5px] font-semibold text-white shadow-[0_16px_30px_-10px_rgba(255,90,31,0.5)] transition-transform hover:-translate-y-0.5"
        >
          Request access
          <Icon path={icons.arrowRight} className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
        </a>
      </div>
    </section>
  )
}

function Footer() {
  const cols = [
    { title: 'Product', links: ['Overview', 'Architecture', 'Changelog'] },
    { title: 'Company', links: ['About', 'Careers', 'Contact'] },
    { title: 'Resources', links: ['Docs', 'Status', 'Support'] },
  ]
  return (
    <footer className="bg-[#0B0E14] pb-10 pt-16">
      <div className="mx-auto max-w-6xl px-6">
        <div className="grid gap-10 border-t border-white/10 pt-12 sm:grid-cols-[1.2fr_repeat(3,1fr)]">
          <div className="flex items-center gap-2.5">
            <span className="relative flex h-7 w-7 items-center justify-center">
              <span className="absolute inset-0 rotate-45 rounded-[7px] bg-white/10" />
              <span className="relative h-1.5 w-1.5 rounded-full bg-[#FF5A1F]" />
            </span>
            <span className="FleetOps-display text-[16px] font-semibold text-white">FleetOps</span>
          </div>
          {cols.map((c) => (
            <div key={c.title}>
              <p className="FleetOps-mono text-[11px] tracking-wide text-white/40">{c.title}</p>
              <ul className="mt-3 space-y-2">
                {c.links.map((l) => (
                  <li key={l}>
                    <a href="#" className="FleetOps-body text-[13.5px] text-white/70 hover:text-white">
                      {l}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
        <p className="mt-12 FleetOps-mono text-[11px] text-white/30">
          &copy; {new Date().getFullYear()} FleetOps. Built on Go, React, and Kafka.
        </p>
      </div>
    </footer>
  )
}

// ---------------------------------------------------------------------------
// Page
// ---------------------------------------------------------------------------

export default function Home() {
  return (
    <div className="min-h-screen bg-white FleetOps-body text-[#0B0E14]">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@500;600;700&family=Manrope:wght@400;500;600;700&family=IBM+Plex+Mono:wght@400;500&display=swap');
        .FleetOps-display { font-family: 'Space Grotesk', ui-sans-serif, sans-serif; }
        .FleetOps-body { font-family: 'Manrope', ui-sans-serif, sans-serif; }
        .FleetOps-mono { font-family: 'IBM Plex Mono', ui-monospace, monospace; }

        @keyframes FleetOps-flow-v {
          0%   { top: 0%;   opacity: 0; }
          12%  { opacity: 1; }
          88%  { opacity: 1; }
          100% { top: 100%; opacity: 0; }
        }
        .FleetOps-flow-v { animation: FleetOps-flow-v 2.6s ease-in-out infinite; }

        @keyframes FleetOps-flow-h {
          0%   { left: 0%;   opacity: 0; }
          12%  { opacity: 1; }
          88%  { opacity: 1; }
          100% { left: 100%; opacity: 0; }
        }
        .FleetOps-flow-h { animation: FleetOps-flow-h 3.2s ease-in-out infinite; }

        @keyframes FleetOps-float {
          0%, 100% { transform: translateY(0); }
          50% { transform: translateY(-5px); }
        }
        .FleetOps-float { animation: FleetOps-float 3.4s ease-in-out infinite; }

        @keyframes FleetOps-blink {
          0%, 100% { opacity: 1; }
          50% { opacity: 0.25; }
        }
        .FleetOps-blink { animation: FleetOps-blink 1.6s ease-in-out infinite; }

        @media (prefers-reduced-motion: reduce) {
          .FleetOps-flow-v, .FleetOps-flow-h, .FleetOps-float, .FleetOps-blink, .FleetOps-route-dot {
            animation: none !important;
          }
        }
      `}</style>

      <Nav />
      <main>
        <Hero />
        <Features />
        <Architecture />
        <EngineeringStack />
      </main>
      <CTA />
      <Footer />
    </div>
  )
}