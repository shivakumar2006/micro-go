import { useState } from 'react'
import { FcGoogle } from "react-icons/fc";
import { useNavigate } from "react-router-dom";
import { useRegisterMutation } from '../Redux/features/auth/authApi';
import { useDispatch } from 'react-redux';
import { setTokens } from '../Redux/features/auth/authSlice';
import { Bounce, toast } from "react-toastify";

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
  lock: (
    <>
      <rect x="5" y="10.5" width="14" height="9" rx="1.5" />
      <path d="M8 10.5V8a4 4 0 0 1 8 0v2.5" />
    </>
  ),
  eye: (
    <>
      <path d="M2.5 12S5.7 5.5 12 5.5 21.5 12 21.5 12 18.3 18.5 12 18.5 2.5 12 2.5 12z" />
      <circle cx="12" cy="12" r="2.6" />
    </>
  ),
  eyeOff: (
    <>
      <path d="M3 3l18 18" />
      <path d="M10.6 5.7A10.6 10.6 0 0 1 12 5.5c6.3 0 9.5 6.5 9.5 6.5a15.6 15.6 0 0 1-3.3 4.2" />
      <path d="M6.6 6.6C4 8.4 2.5 12 2.5 12s3.2 6.5 9.5 6.5a9.9 9.9 0 0 0 3.9-.8" />
      <path d="M9.9 10a2.6 2.6 0 0 0 3.6 3.6" />
    </>
  ),
  arrowRight: (
    <>
      <path d="M4 12h16" />
      <path d="M14 6l6 6-6 6" />
    </>
  ),
  person: (
  <>
    <path d="M12 12a4 4 0 1 0-4-4 4 4 0 0 0 4 4z" />
    <path d="M4 20a8 8 0 0 1 16 0" />
  </>
),
}

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

const GOOGLE_PATH =
  'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10c8.39 0 10.31-7.8 9.54-11H12v3.5h5.3c-.45 2.28-2.4 3.9-5.3 3.9-3.48 0-6.3-2.82-6.3-6.3S8.52 5.8 12 5.8c1.57 0 2.99.56 4.09 1.49l2.48-2.48C16.86 3.18 14.57 2 12 2z'
// ---------------------------------------------------------------------------
// Shared building blocks
// ---------------------------------------------------------------------------

function Logo({ dark = false }) {
  return (
    <a href="#top" className="flex items-center gap-2.5">
      <span className="relative flex h-7 w-7 items-center justify-center">
        <span className={`absolute inset-0 rotate-45 rounded-[7px] ${dark ? 'bg-white/10' : 'bg-[#0B0E14]'}`} />
        <span className="relative h-1.5 w-1.5 rounded-full bg-[#FF5A1F]" />
      </span>
      <span className={`FleetOps-display text-[17px] font-semibold tracking-tight ${dark ? 'text-white' : 'text-[#0B0E14]'}`}>
        FleetOps
      </span>
    </a>
  )
}

function FieldInput({ label, type = 'text', icon, value, onChange, placeholder, trailing, autoComplete }) {
  return (
    <label className="block">
      <span className="FleetOps-mono text-[10.5px] uppercase tracking-wide text-[#5B6472]">{label}</span>
      <div className="relative mt-1.5">
        <span className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-[#5B6472]">
          {NodeIcon(icon, 'h-[18px] w-[18px]')}
        </span>
        <input
          type={type}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          autoComplete={autoComplete}
          className="w-full rounded-xl border border-slate-200 bg-white py-3 pl-11 pr-11 FleetOps-body text-[14.5px] text-[#0B0E14] placeholder:text-slate-400 transition-colors focus:border-[#0B0E14] focus:outline-none focus:ring-2 focus:ring-[#0B0E14]/10"
        />
        {trailing && <span className="absolute right-3.5 top-1/2 -translate-y-1/2 flex items-center">{trailing}</span>}
      </div>
    </label>
  )
}

// vertical rail with a traveling "signal" dot, echoes the homepage diagram
function AuthFlow() {
  const stages = [
    { label: 'Client', sub: 'browser session' },
    { label: 'API gateway', sub: 'token verified' },
    { label: 'Services', sub: 'scoped access granted' },
    { label: 'Event bus', sub: 'session logged' },
  ]
  return (
    <div>
      {stages.map((s, i) => (
        <div key={s.label}>
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 flex-none items-center justify-center rounded-lg bg-white/10">
              <span className="FleetOps-mono text-[11px] text-white/70">{i + 1}</span>
            </div>
            <div>
              <p className="FleetOps-display text-[13.5px] font-semibold text-white">{s.label}</p>
              <p className="FleetOps-mono text-[10.5px] text-white/45">{s.sub}</p>
            </div>
          </div>
          {i < stages.length - 1 && (
            <div className="relative ml-4 h-6 w-px bg-white/15">
              <span
                className="FleetOps-flow-v absolute left-1/2 h-1.5 w-1.5 -translate-x-1/2 rounded-full bg-[#FF5A1F]"
                style={{ animationDelay: `${i * 0.3}s` }}
              />
            </div>
          )}
        </div>
      ))}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Page
// ---------------------------------------------------------------------------

export default function SignUp() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState("customer");

  const [showPassword, setShowPassword] = useState(false)
  const navigate = useNavigate();

  const [register, { isLoading }] = useRegisterMutation();
  
  const dispatch = useDispatch();

  const emailValid = /\S+@\S+\.\S+/.test(email)
  const passwordError =
  password.length > 0 && password.length < 8
    ? "Password must be at least 8 characters."
    : "";
  const canSubmit = name.trim() && emailValid && password.length >= 8 && role;

  const handleRegister = async (e) => {
    e.preventDefault();

    try{
        const res = await register({name, email, password, role}).unwrap();

        dispatch(setTokens({
            accessToken: res.access_token,
            refreshToken: res.refresh_token,
        }));

        navigate("/login");
        toast.success('User Registered successfully 🎉', {
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
    } catch(err) {
        console.error("Error during registration:", err);
        toast.error('Failed to register user 😞', {
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
    }
  }

  return (
    <div className="min-h-screen bg-white FleetOps-body text-[#0B0E14] lg:flex">
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

        @keyframes FleetOps-blink {
          0%, 100% { opacity: 1; }
          50% { opacity: 0.25; }
        }
        .FleetOps-blink { animation: FleetOps-blink 1.6s ease-in-out infinite; }

        @media (prefers-reduced-motion: reduce) {
          .FleetOps-flow-v, .FleetOps-blink { animation: none !important; }
        }
      `}</style>

      {/* left — brand panel, hidden below lg */}
      <div className="hidden lg:flex lg:w-[42%] lg:flex-col lg:justify-between bg-[#0B0E14] px-12 py-12 xl:w-[38%] xl:px-14">
        <div>
          <Logo dark />
          <span className="mt-10 inline-flex items-center gap-2 rounded-full border border-white/15 px-3 py-1 FleetOps-mono text-[11px] tracking-wide text-white/60">
            <span className="h-1.5 w-1.5 rounded-full bg-[#FF5A1F] FleetOps-blink" />
            private beta — built on go &amp; kafka
          </span>
          <h1 className="mt-6 FleetOps-display text-3xl font-semibold leading-[1.15] tracking-tight text-white">
            Welcome back to the control plane.
          </h1>
          <p className="mt-4 max-w-sm FleetOps-body text-[14.5px] leading-relaxed text-white/55">
            Pick up right where you left off — live vehicle state, orders and payments,
            all in one dashboard.
          </p>
        </div>

        <div className="max-w-xs">
          <AuthFlow />
        </div>

        <div className="grid max-w-xs grid-cols-2 gap-6 border-t border-white/10 pt-6">
          <div>
            <p className="FleetOps-display text-xl font-semibold text-white">128</p>
            <p className="FleetOps-mono text-[10.5px] text-white/45">vehicles online</p>
          </div>
          <div>
            <p className="FleetOps-display text-xl font-semibold text-white">37</p>
            <p className="FleetOps-mono text-[10.5px] text-white/45">orders in flight</p>
          </div>
        </div>
      </div>

      {/* right — form panel */}
      <div className="flex flex-1 flex-col">
        <header className="flex items-center justify-between px-6 py-6 sm:px-10">
          <span className="lg:hidden">
            <Logo />
          </span>
          <span className="hidden lg:block" />
          <a
            className="FleetOps-body text-[13.5px] font-medium text-[#5B6472] transition-colors hover:text-[#0B0E14]"
          >
            Already have an account? <span onClick={() => navigate("/login")} className="font-semibold text-[#0B0E14] cursor-pointer">Sign in</span>
          </a>
        </header>

        <main className="flex flex-1 items-center justify-center px-6 pb-16">
          <div className="w-full max-w-[420px]">
            <div className="rounded-[28px] border border-slate-200 bg-white p-8 shadow-[0_1px_2px_rgba(15,23,42,0.04)] sm:p-10">
              <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">sign up</span>
              <h2 className="mt-2 FleetOps-display text-2xl font-semibold tracking-tight text-[#0B0E14]">
                Sign up to FleetOps
              </h2>
              <p className="mt-2 FleetOps-body text-[14px] text-[#5B6472]">
                Enter your credentials to access the dashboard.
              </p>

              <form className="mt-8 space-y-5" onSubmit={handleRegister}>
                <FieldInput
                  label="Name"
                  type="text"
                  icon="person"
                  value={name}
                  onChange={setName}
                  placeholder="John Doe"
                  autoComplete="name"
                />

                <FieldInput
                  label="Email address"
                  type="email"
                  icon="mail"
                  value={email}
                  onChange={setEmail}
                  placeholder="you@company.com"
                  autoComplete="email"
                />

                <FieldInput
                  label="Password"
                  type={showPassword ? 'text' : 'password'}
                  icon="lock"
                  value={password}
                  onChange={setPassword}
                  placeholder="••••••••"
                  autoComplete="current-password"
                  trailing={
                    <button
                      type="button"
                      onClick={() => setShowPassword((v) => !v)}
                      className="text-[#5B6472] transition-colors hover:text-[#0B0E14]"
                      aria-label={showPassword ? 'Hide password' : 'Show password'}
                    >
                      {NodeIcon(showPassword ? 'eyeOff' : 'eye', 'h-[18px] w-[18px]')}
                    </button>
                  }
                />
                {passwordError && (
                    <p className="mt-2 text-sm text-red-500">
                        {passwordError}
                    </p>
                )}

                <div className='flex flex-col items-start justify-between'>
                  <p className='text-gray-500 text-sm py-2'>Select Role</p>
                  <select 
                    value={role}
                    onChange={(e) => setRole(e.target.value)}
                    className="w-full h-10 border border-gray-300 rounded-lg px-3 text-sm text-gray-500 focus:outline-none focus:ring-0 focus:border-[#0B0E14]">
                    <option value="customer">Customer</option>
                    <option value="admin">Admin</option>
                  </select>
                </div>

                <div className="flex items-center justify-between">
                  <a
                    href="#"
                    className="FleetOps-body text-[13px] font-medium text-[#0B0E14] underline decoration-slate-300 underline-offset-4 hover:decoration-[#0B0E14]"
                  >
                    Forgot password?
                  </a>
                </div>

                <button
                  type="submit"
                  disabled={!canSubmit}
                  className="group flex w-full items-center justify-center gap-2 rounded-full bg-[#FF5A1F] px-6 py-3.5 FleetOps-body text-[14.5px] font-semibold text-white shadow-[0_12px_24px_-8px_rgba(255,90,31,0.55)] transition-all hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-35 disabled:hover:translate-y-0"
                >
                  {isLoading ? "Creating Account..." : "Sign up"}
                  <Icon path={icons.arrowRight} className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
                </button>
              </form>

              <div className="my-7 flex items-center gap-3">
                <div className="h-px flex-1 bg-slate-200" />
                <span className="FleetOps-mono text-[10px] tracking-wide text-[#5B6472]">or continue with</span>
                <div className="h-px flex-1 bg-slate-200" />
              </div>

              <button
                type="button"
                className="flex w-full items-center justify-center gap-2.5 rounded-full border border-slate-200 bg-white px-6 py-3.5 FleetOps-body text-[14px] font-semibold text-[#0B0E14] transition-colors hover:bg-[#F3F5F8]"
              >
            
                <FcGoogle className='text-2xl'/>
                Continue with Google
              </button>
            </div>

            <p className="mt-6 text-center FleetOps-body text-[13.5px] text-[#5B6472] lg:hidden">
              New to FleetOps?{' '}
              <a
                href="#"
                className="font-semibold text-[#0B0E14] underline decoration-slate-300 underline-offset-4 hover:decoration-[#0B0E14]"
              >
                Create an account
              </a>
            </p>
          </div>
        </main>
      </div>
    </div>
  )
}