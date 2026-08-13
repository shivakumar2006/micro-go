import { useState } from 'react'
import {
    useGetOrdersByUserIdQuery,
    useCancelOrderMutation
} from '../Redux/features/order/order';
import { useSelector } from "react-redux";

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
  box: (
    <>
      <path d="M3 8l9-4 9 4-9 4-9-4z" />
      <path d="M3 8v8l9 4 9-4V8" />
      <path d="M12 12v8" />
    </>
  ),
  chevronDown: <path d="M6 9l6 6 6-6" />,
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
  inbox: (
    <>
      <path d="M3.5 12h5l1.5 3h4l1.5-3h5" />
      <path d="M5 6.5h14L21 12v6a1.5 1.5 0 0 1-1.5 1.5h-15A1.5 1.5 0 0 1 3 18v-6l2-5.5z" />
    </>
  ),
  cancel: (
    <>
      <circle cx="12" cy="12" r="9" />
      <path d="M9 9l6 6M15 9l-6 6" />
    </>
  ),
  user: (
    <>
      <circle cx="12" cy="8.2" r="3.2" />
      <path d="M5 19.5c0-3.3 3.13-5.5 7-5.5s7 2.2 7 5.5" />
    </>
  ),
}

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

const statusMeta = {
  pending: { label: 'Pending', bg: '#FEF3C7', text: '#B45309', dot: '#B45309' },
  confirmed: { label: 'Confirmed', bg: '#DBEAFE', text: '#1D4ED8', dot: '#1D4ED8' },
  processing: { label: 'Processing', bg: '#DBEAFE', text: '#1D4ED8', dot: '#1D4ED8' },
  shipped: { label: 'Shipped', bg: '#EDE9FE', text: '#6D28D9', dot: '#6D28D9' },
  delivered: { label: 'Delivered', bg: '#DCFCE7', text: '#15803D', dot: '#16A34A' },
  cancelled: { label: 'Cancelled', bg: '#F1F5F9', text: '#475569', dot: '#94A3B8' },
  paid: { label: 'Paid', bg: '#DCFCE7', text: '#15803D', dot: '#16A34A' },
}

const cancellableStatuses = ['pending', 'confirmed']

function formatPrice(amount) {
  const n = Number(amount)
  if (Number.isNaN(n)) return '—'
  return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(n)
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('en-IN', { day: '2-digit', month: 'short', year: 'numeric' })
}

// reads vehicle details off an order item whether they're nested under
// `item.vehicle` or flattened directly on the item
function getOrderItemInfo(item) {
  const v = item.vehicle || item
  return {
    name: v.name ?? 'Vehicle',
    brand: v.brand ?? '',
    model: v.model ?? '',
    image_url: v.image_url ?? null,
    price: Number(item.price ?? v.price) || 0,
    quantity: item.quantity ?? 1,
  }
}

function Thumb({ src, alt }) {
  const [failed, setFailed] = useState(false)
  if (!src || failed) {
    return (
      <div className="flex h-12 w-12 flex-none items-center justify-center rounded-lg bg-[#F3F5F8] text-[#5B6472]">
        {NodeIcon('image', 'h-4 w-4')}
      </div>
    )
  }
  return (
    <img
      src={src}
      alt={alt}
      onError={() => setFailed(true)}
      className="h-12 w-12 flex-none rounded-lg border border-slate-100 object-cover"
    />
  )
}

function CancelConfirm({ order, onCancel, onConfirm, isCancelling }) {
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
          {NodeIcon('cancel', 'h-5 w-5')}
        </div>
        <h2 className="mt-4 FleetOps-display text-lg font-semibold text-[#0B0E14]">
          Cancel order #{order.id}?
        </h2>
        <p className="mt-2 FleetOps-body text-[13.5px] text-[#5B6472]">This can&rsquo;t be undone.</p>
        <div className="mt-6 flex items-center justify-center gap-3">
          <button
            type="button"
            onClick={onCancel}
            className="flex-1 rounded-full border border-slate-200 px-5 py-2.5 FleetOps-body text-[13.5px] font-medium text-[#0B0E14] transition-colors hover:bg-[#F3F5F8]"
          >
            Keep order
          </button>
          <button
            type="button"
            onClick={onConfirm}
            disabled={isCancelling}
            className="flex-1 rounded-full bg-[#DC2626] px-5 py-2.5 FleetOps-body text-[13.5px] font-semibold text-white transition-colors hover:bg-[#B91C1C] disabled:cursor-not-allowed disabled:opacity-50"
          >
            {isCancelling ? 'Cancelling…' : 'Cancel order'}
          </button>
        </div>
      </div>
    </div>
  )
}

function OrderCard({ order, onCancelClick }) {
  const [expanded, setExpanded] = useState(false)
  const s = statusMeta[order.status] || { label: order.status || 'Unknown', bg: '#F1F5F9', text: '#475569', dot: '#475569' }
  const items = order.items || []
  const canCancel = cancellableStatuses.includes(order.status)

  return (
    <div className="rounded-2xl border border-slate-200 bg-white">
      <button
        type="button"
        onClick={() => setExpanded((v) => !v)}
        className="flex w-full items-center justify-between gap-4 px-5 py-4 text-left"
      >
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 flex-none items-center justify-center rounded-xl bg-[#0B0E14] text-white">
            {NodeIcon('box', 'h-4 w-4')}
          </div>
          <div>
            <p className="FleetOps-display text-[14.5px] font-semibold text-[#0B0E14]">Order #{order.id}</p>
            <p className="flex items-center gap-1.5 FleetOps-mono text-[10.5px] text-[#5B6472]">
              {NodeIcon('calendar', 'h-3 w-3')}
              {formatDate(order.created_at)}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <span
            className="hidden items-center gap-1.5 rounded-full px-2.5 py-1 FleetOps-mono text-[10px] font-medium sm:inline-flex"
            style={{ backgroundColor: s.bg, color: s.text }}
          >
            <span className="h-1.5 w-1.5 rounded-full" style={{ backgroundColor: s.dot }} />
            {s.label}
          </span>
          <span className="FleetOps-display text-[14.5px] font-semibold text-[#0B0E14]">{formatPrice(order.total)}</span>
          <span className={`text-[#5B6472] transition-transform ${expanded ? 'rotate-180' : ''}`}>
            {NodeIcon('chevronDown', 'h-4 w-4')}
          </span>
        </div>
      </button>

      <span
        className="mx-5 mb-3 inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 FleetOps-mono text-[10px] font-medium sm:hidden"
        style={{ backgroundColor: s.bg, color: s.text }}
      >
        <span className="h-1.5 w-1.5 rounded-full" style={{ backgroundColor: s.dot }} />
        {s.label}
      </span>

      {expanded && (
        <div className="border-t border-slate-100 px-5 py-4">
          <div className="space-y-3">
            {items.map((item, i) => {
              const info = getOrderItemInfo(item)
              return (
                <div key={item.id ?? i} className="flex items-center gap-3">
                  <Thumb src={info.image_url} alt={info.name} />
                  <div className="flex-1">
                    <p className="FleetOps-body text-[13.5px] font-medium text-[#0B0E14]">{info.name}</p>
                    <p className="FleetOps-mono text-[10.5px] text-[#5B6472]">
                      {info.brand} {info.model && `· ${info.model}`} · Qty {info.quantity}
                    </p>
                  </div>
                  <p className="FleetOps-body text-[13.5px] font-medium text-[#0B0E14]">
                    {formatPrice(info.price * info.quantity)}
                  </p>
                </div>
              )
            })}
          </div>

          <div className="mt-4 flex items-center justify-between border-t border-slate-100 pt-4">
            <span className="FleetOps-body text-[13.5px] font-semibold text-[#0B0E14]">Total</span>
            <span className="FleetOps-display text-[16px] font-semibold text-[#0B0E14]">{formatPrice(order.total)}</span>
          </div>

          {canCancel && (
            <button
              type="button"
              onClick={() => onCancelClick(order)}
              className="mt-4 flex items-center gap-1.5 rounded-full border border-slate-200 px-4 py-2 FleetOps-body text-[13px] font-medium text-[#5B6472] transition-colors hover:bg-red-50 hover:text-[#DC2626]"
            >
              {NodeIcon('cancel', 'h-3.5 w-3.5')}
              Cancel order
            </button>
          )}
        </div>
      )}
    </div>
  )
}

export default function Orders() {
  const { User } = useSelector((state) => state.authReducer);

  console.log("USER FROM REDUX:", User);

  const userId = User?.id;
  console.log("USER ID FROM REDUX:", userId);

  const { data, isLoading, error } = useGetOrdersByUserIdQuery(userId, { skip: !userId })
  const [cancelOrder, { isLoading: isCancelling }] = useCancelOrderMutation()
  const [cancelTarget, setCancelTarget] = useState(null)

  // --- unwrap the response — see note at top of file -------------------------
  const orders = data?.data ?? (Array.isArray(data) ? data : []) ?? []

  async function handleConfirmCancel() {
    try {
      await cancelOrder(cancelTarget.id).unwrap()
    } finally {
      setCancelTarget(null)
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

      <main className="mx-auto max-w-3xl px-6 py-10">
        <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">orders</span>
        <h1 className="mt-1 FleetOps-display text-2xl font-semibold tracking-tight text-[#0B0E14] sm:text-3xl">
          Your orders
        </h1>

        {!userId ? (
          <div className="mt-8 flex flex-col items-center gap-3 rounded-2xl border border-dashed border-slate-300 bg-white py-20 text-center">
            {NodeIcon('user', 'h-9 w-9 text-[#5B6472]')}
            <p className="FleetOps-display text-[16px] font-semibold text-[#0B0E14]">No orders found</p>
          </div>
        ) : isLoading ? (
          <div className="mt-8 rounded-2xl border border-slate-200 bg-white py-16 text-center">
            <p className="FleetOps-body text-[13.5px] text-[#5B6472]">Loading your orders…</p>
          </div>
        ) : error ? (
          <div className="mt-8 rounded-2xl border border-red-100 bg-red-50 py-16 text-center">
            <p className="FleetOps-body text-[13.5px] text-[#DC2626]">Couldn&rsquo;t load your orders. Please try again.</p>
          </div>
        ) : orders.length === 0 ? (
          <div className="mt-8 flex flex-col items-center gap-3 rounded-2xl border border-dashed border-slate-300 bg-white py-20 text-center">
            {NodeIcon('inbox', 'h-9 w-9 text-[#5B6472]')}
            <p className="FleetOps-display text-[16px] font-semibold text-[#0B0E14]">No orders yet</p>
            <p className="FleetOps-body text-[13.5px] text-[#5B6472]">Orders you place will show up here.</p>
          </div>
        ) : (
          <div className="mt-8 space-y-4">
            {orders.map((order) => (
              <OrderCard key={order.id} order={order} onCancelClick={setCancelTarget} />
            ))}
          </div>
        )}
      </main>

      {cancelTarget && (
        <CancelConfirm
          order={cancelTarget}
          isCancelling={isCancelling}
          onCancel={() => setCancelTarget(null)}
          onConfirm={handleConfirmCancel}
        />
      )}
    </div>
  )
}