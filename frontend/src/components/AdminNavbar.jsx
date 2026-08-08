import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useLogoutMutation, useMeQuery } from '../Redux/features/auth/authApi';
import { useDispatch, useSelector } from 'react-redux';
import { clearAuth } from '../Redux/features/auth/authSlice';
import { Bounce, toast } from 'react-toastify';

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
  bell: (
    <>
      <path d="M6 10.2a6 6 0 0 1 12 0c0 3.8 1.3 5.3 1.3 5.3H4.7S6 14 6 10.2z" />
      <path d="M10 18a2 2 0 0 0 4 0" />
    </>
  ),
  chevronDown: <path d="M6 9l6 6 6-6" />,
  arrowLeft: (
    <>
      <path d="M19 12H5" />
      <path d="M11 18l-6-6 6-6" />
    </>
  ),
  menu: <path d="M4 7h16M4 12h16M4 17h16" />,
  close: <path d="M6 6l12 12M18 6L6 18" />,
  user: (
    <>
      <circle cx="12" cy="8.2" r="3.2" />
      <path d="M5 19.5c0-3.3 3.13-5.5 7-5.5s7 2.2 7 5.5" />
    </>
  ),
  settings: (
    <>
      <circle cx="12" cy="12" r="2.8" />
      <path d="M19.4 13.5a1.7 1.7 0 0 0 .34 1.87l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.7 1.7 0 0 0-1.87-.34 1.7 1.7 0 0 0-1.04 1.56V19.5a2 2 0 1 1-4 0v-.09a1.7 1.7 0 0 0-1.11-1.56 1.7 1.7 0 0 0-1.87.34l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.7 1.7 0 0 0 .34-1.87 1.7 1.7 0 0 0-1.56-1.04H4.5a2 2 0 1 1 0-4h.09A1.7 1.7 0 0 0 6.15 6.94a1.7 1.7 0 0 0-.34-1.87l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.7 1.7 0 0 0 1.87.34H10.5a1.7 1.7 0 0 0 1.04-1.56V.5a2 2 0 1 1 4 0v.09a1.7 1.7 0 0 0 1.04 1.56 1.7 1.7 0 0 0 1.87-.34l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.7 1.7 0 0 0-.34 1.87V6.5a1.7 1.7 0 0 0 1.56 1.04h.09a2 2 0 1 1 0 4h-.09a1.7 1.7 0 0 0-1.5 1z" />
    </>
  ),
  logout: (
    <>
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
      <path d="M16 17l5-5-5-5" />
      <path d="M21 12H9" />
    </>
  ),
}

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

function Avatar({ initials }) {
  return (
    <span className="flex h-8 w-8 flex-none items-center justify-center rounded-full bg-[#0B0E14] FleetOps-mono text-[12px] font-semibold text-white">
      {initials}
    </span>
  )
}

function getInitials(name) {
  return name
    .split(' ')
    .filter(Boolean)
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()
}

export default function AdminNavbar({ userName = 'Admin' }) {
  const [profileOpen, setProfileOpen] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false);

  const dispatch = useDispatch();
  const navigate = useNavigate();
  
  const [logout] = useLogoutMutation();
  const { data, isLoading } = useMeQuery();

  const { refreshToken } = useSelector((state) => state.authReducer);
  
  if (isLoading) {
    return <div>Loading...</div>
  }

  const handleLogout = async () => {
    try {
        await logout(refreshToken).unwrap();

        dispatch(clearAuth());

        toast.success('Admin Logged out successfully 🎉', {
            position: "top-right",
            autoClose: 5000,
            hideProgressBar: false,
            closeOnClick: false,
            pauseOnHover: true,
            draggable: true,
            progress: undefined,
            theme: "light",
            transition: Bounce,
        });

        navigate("/login");
    } catch(err) {
        console.error(err);

        dispatch(clearAuth())

        toast.error('Failed to Logged out the admin 😔', {
            position: "top-right",
            autoClose: 5000,
            hideProgressBar: false,
            closeOnClick: false,
            pauseOnHover: true,
            draggable: true,
            progress: undefined,
            theme: "light",
            transition: Bounce,
        });

        navigate("/login")
    }
  }

  console.log(data);

  return (
    <header className="sticky top-0 z-40 border-b border-slate-200/70 bg-white/85 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
        <div className="flex items-center gap-3">
          <a href="#top" className="flex items-center gap-2.5">
            <span className="relative flex h-7 w-7 items-center justify-center">
              <span className="absolute inset-0 rotate-45 rounded-[7px] bg-[#0B0E14]" />
              <span className="relative h-1.5 w-1.5 rounded-full bg-[#FF5A1F]" />
            </span>
            <span className="FleetOps-display text-[17px] font-semibold tracking-tight text-[#0B0E14]">
              FleetOps
            </span>
          </a>
          <span className="rounded-full border border-[#35455C]/30 bg-[#35455C]/[0.07] px-2 py-0.5 FleetOps-mono text-[10px] font-medium tracking-wide text-[#35455C]">
            ADMIN
          </span>
        </div>

        <div className="hidden items-center gap-5 md:flex">
          <a
            href="#"
            className="flex items-center gap-1.5 FleetOps-body text-[13.5px] font-medium text-[#5B6472] transition-colors hover:text-[#0B0E14]"
          >
            {NodeIcon('arrowLeft', 'h-4 w-4')}
            Back to dashboard
          </a>

          <button
            type="button"
            className="flex h-9 w-9 items-center justify-center rounded-full text-[#5B6472] transition-colors hover:bg-[#F3F5F8] hover:text-[#0B0E14]"
            aria-label="Notifications"
          >
            {NodeIcon('bell', 'h-[18px] w-[18px]')}
          </button>

          <div className="relative">
            <button
              type="button"
              onClick={() => setProfileOpen((v) => !v)}
              className="flex items-center gap-2 rounded-full border border-slate-200 py-1 pl-1 pr-2.5 transition-colors hover:bg-[#F3F5F8]"
            >
              <Avatar initials={getInitials(userName)} />
              <span className="FleetOps-body text-[13px] font-medium text-[#0B0E14]">{userName}</span>
              {NodeIcon('chevronDown', 'h-3.5 w-3.5 text-[#5B6472]')}
            </button>

            {profileOpen && (
              <>
                <div className="fixed inset-0 z-10" onClick={() => setProfileOpen(false)} />
                <div className="absolute right-0 z-20 mt-2 w-48 rounded-xl border border-slate-200 bg-white p-1.5 shadow-[0_12px_24px_-8px_rgba(15,23,42,0.18)]">
                  <a
                    onClick={() => navigate("/admin/profile")}
                    className="flex items-center gap-2.5 rounded-lg px-3 py-2 FleetOps-body text-[13.5px] text-[#0B0E14] hover:bg-[#F3F5F8]"
                  >
                    {NodeIcon('user', 'h-4 w-4 text-[#5B6472]')}
                    Profile
                  </a>
                  <a
                    href="#"
                    className="flex items-center gap-2.5 rounded-lg px-3 py-2 FleetOps-body text-[13.5px] text-[#0B0E14] hover:bg-[#F3F5F8]"
                  >
                    {NodeIcon('settings', 'h-4 w-4 text-[#5B6472]')}
                    Settings
                  </a>
                  <div className="my-1.5 h-px bg-slate-100" />
                  <a
                    onClick={handleLogout}
                    className="flex items-center gap-2.5 rounded-lg px-3 py-2 FleetOps-body text-[13.5px] text-[#0B0E14] hover:bg-[#F3F5F8] cursor-pointer"
                  >
                    {NodeIcon('logout', 'h-4 w-4 text-[#5B6472]')}
                    Sign out
                  </a>
                </div>
              </>
            )}
          </div>
        </div>

        <button
          className="flex h-9 w-9 items-center justify-center rounded-lg text-[#0B0E14] md:hidden"
          onClick={() => setMobileOpen((v) => !v)}
          aria-label="Toggle menu"
        >
          {NodeIcon(mobileOpen ? 'close' : 'menu', 'h-5 w-5')}
        </button>
      </div>

      {mobileOpen && (
        <div className="flex flex-col gap-1 border-t border-slate-200 bg-white px-6 py-4 md:hidden">
          <a href="#" className="flex items-center gap-1.5 rounded-lg px-3 py-2.5 FleetOps-body text-sm text-[#0B0E14]">
            {NodeIcon('arrowLeft', 'h-4 w-4')}
            Back to dashboard
          </a>
          <div className="mt-3 flex items-center gap-2.5 border-t border-slate-100 pt-3">
            <Avatar initials={getInitials(userName)} />
            <span className="FleetOps-body text-[13.5px] font-medium text-[#0B0E14]">{userName}</span>
          </div>
        </div>
      )}
    </header>
  )
}