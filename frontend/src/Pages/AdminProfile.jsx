import { useState } from 'react';
import { useMeQuery } from '../Redux/features/auth/authApi';

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
  mail: (
    <>
      <path d="M3.5 6.5h17a1 1 0 0 1 1 1v9a1 1 0 0 1-1 1h-17a1 1 0 0 1-1-1v-9a1 1 0 0 1 1-1z" />
      <path d="M3 7l9 6.5L21 7" />
    </>
  ),
  shield: (
    <>
      <path d="M12 3.2l7 2.8v5c0 4.6-3 7.6-7 9-4-1.4-7-4.4-7-9v-5l7-2.8z" />
      <path d="M9 12l2 2 4-4" />
    </>
  ),
  edit: (
    <>
      <path d="M4 20.5l.9-3.6 11-11a2 2 0 0 1 2.8 0l1.8 1.8a2 2 0 0 1 0 2.8l-11 11-3.6.9z" />
      <path d="M14.5 7l2.5 2.5" />
    </>
  ),
  lock: (
    <>
      <rect x="5" y="10.5" width="14" height="9" rx="1.5" />
      <path d="M8 10.5V8a4 4 0 0 1 8 0v2.5" />
    </>
  ),
}

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

const sampleUser = {
  id: 1,
  name: 'Admin',
  email: 'admin123@gmail.com',
  role_name: 'admin',
  created_at: '0001-01-01T00:00:00Z',
  updated_at: '0001-01-01T00:00:00Z',
}

function getInitials(name = '') {
  return name
    .split(' ')
    .filter(Boolean)
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()
}

function DetailRow({ icon, label, value }) {
  return (
    <div className="flex items-center gap-3 border-b border-slate-100 py-4 last:border-0">
      <div className="flex h-9 w-9 flex-none items-center justify-center rounded-lg bg-[#F3F5F8] text-[#5B6472]">
        {NodeIcon(icon, 'h-4 w-4')}
      </div>
      <div>
        <p className="FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">{label}</p>
        <p className="FleetOps-body text-[14px] font-medium text-[#0B0E14]">{value}</p>
      </div>
    </div>
  )
}

export default function Profile({ user = sampleUser }) {
  const { name, email, role_name } = user
  const [editing, setEditing] = useState(false);

  const {data, isLoading} = useMeQuery();

  const User = data;

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <div className="min-h-screen bg-[#F3F5F8] FleetOps-body text-[#0B0E14]">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@500;600;700&family=Manrope:wght@400;500;600;700&family=IBM+Plex+Mono:wght@400;500&display=swap');
        .FleetOps-display { font-family: 'Space Grotesk', ui-sans-serif, sans-serif; }
        .FleetOps-body { font-family: 'Manrope', ui-sans-serif, sans-serif; }
        .FleetOps-mono { font-family: 'IBM Plex Mono', ui-monospace, monospace; }
      `}</style>

      <main className="mx-auto max-w-md px-6 py-14">
        <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">profile</span>
        <h1 className="mt-1 FleetOps-display text-2xl font-semibold tracking-tight text-[#0B0E14] sm:text-3xl">
          Your account
        </h1>

        <div className="mt-8 rounded-[28px] border border-slate-200 bg-white p-8 text-center shadow-[0_1px_2px_rgba(15,23,42,0.04)]">
          <span className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-[#0B0E14] FleetOps-display text-2xl font-semibold text-white">
            {getInitials(name)}
          </span>

          <h2 className="mt-4 FleetOps-display text-xl font-semibold text-[#0B0E14]">{name}</h2>

          <span className="mt-2 inline-flex items-center gap-1.5 rounded-full border border-[#35455C]/25 bg-[#35455C]/[0.06] px-3 py-1 FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#35455C]">
            {NodeIcon('shield', 'h-3.5 w-3.5')}
            {User?.role}
          </span>

          <div className="mt-6 text-left">
            <DetailRow icon="mail" label="Email" value={User?.email} />
            <DetailRow icon="shield" label="Role" value={User?.role} />
          </div>

          <div className="mt-7 flex items-center gap-3">
            <button
              type="button"
              onClick={() => setEditing(true)}
              className="flex flex-1 items-center justify-center gap-1.5 rounded-full bg-[#0B0E14] px-5 py-2.5 FleetOps-body text-[13.5px] font-medium text-white transition-colors hover:bg-[#1a2030]"
            >
              {NodeIcon('edit', 'h-4 w-4')}
              Edit profile
            </button>
            <button
              type="button"
              className="flex flex-1 items-center justify-center gap-1.5 rounded-full border border-slate-200 px-5 py-2.5 FleetOps-body text-[13.5px] font-medium text-[#0B0E14] transition-colors hover:bg-[#F3F5F8]"
            >
              {NodeIcon('lock', 'h-4 w-4')}
              Change password
            </button>
          </div>
        </div>
      </main>

      {editing && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-[#0B0E14]/50 px-4 backdrop-blur-sm"
          onClick={() => setEditing(false)}
        >
          <div
            className="w-full max-w-sm rounded-[28px] border border-slate-200 bg-white p-8 text-center shadow-[0_30px_60px_-30px_rgba(15,23,42,0.35)]"
            onClick={(e) => e.stopPropagation()}
          >
            <p className="FleetOps-body text-[13.5px] text-[#5B6472]">
              Hook this up to your update-profile mutation whenever it&rsquo;s ready — this is just a placeholder for now.
            </p>
            <button
              type="button"
              onClick={() => setEditing(false)}
              className="mt-5 w-full rounded-full bg-[#0B0E14] px-5 py-2.5 FleetOps-body text-[13.5px] font-medium text-white transition-colors hover:bg-[#1a2030]"
            >
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  )
}