import { useState } from 'react'
import {
  useAddToCartMutation,
  useGetUserCartQuery,
  useUpdateCartItemMutation,
  useDeleteCartItemMutation,
  useClearCartMutation,
  useGetCartTotalQuery,
  useCountItemsQuery,
} from '../Redux/features/cart/cartApi'

/*
  FleetOps — cart
  ------------------------------------------------------------------
  Wired to cartApi (RTK Query) — no local mock state. No navbar here
  on purpose, per your call — drop this inside whatever layout you're
  adding the navbar to.

  Two shape assumptions, both handled defensively since the exact
  response payloads weren't given — adjust the two spots marked below
  if your backend differs:

  1. Cart items — GET /cart is assumed to return either
     { data: [...] } or a bare array, and each item is assumed to
     either nest vehicle details under `item.vehicle` or have them
     flattened directly on the item. getVehicleInfo() below tries
     both so the UI degrades gracefully either way.
  2. Total / count — GET /cart/total and /cart/count are assumed to
     return { total } / { count } (or nested under `.data`). If
     neither is present, the page falls back to computing them
     client-side from the loaded cart items.
*/

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
  bag: (
    <>
      <path d="M3 4h2l2.2 10a2 2 0 0 0 2 1.6h6.6a2 2 0 0 0 2-1.6L19.5 8H6" />
      <circle cx="9" cy="20" r="1.3" />
      <circle cx="17" cy="20" r="1.3" />
    </>
  ),
  plus: (
    <>
      <path d="M12 5v14" />
      <path d="M5 12h14" />
    </>
  ),
  minus: <path d="M5 12h14" />,
  trash: (
    <>
      <path d="M4 7h16" />
      <path d="M9 7V5a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
      <path d="M6 7l1 13a1.5 1.5 0 0 0 1.5 1.4h7a1.5 1.5 0 0 0 1.5-1.4L18 7" />
      <path d="M10 11v6" />
      <path d="M14 11v6" />
    </>
  ),
  image: (
    <>
      <rect x="3" y="4.5" width="18" height="15" rx="1.5" />
      <circle cx="8.5" cy="9.5" r="1.5" />
      <path d="M3 16l5-5 4 4 3-3 6 6" />
    </>
  ),
  arrowRight: (
    <>
      <path d="M4 12h16" />
      <path d="M14 6l6 6-6 6" />
    </>
  ),
  close: <path d="M6 6l12 12M18 6L6 18" />,
}

const NodeIcon = (name, className) => <Icon path={icons[name]} className={className} />

function formatPrice(amount) {
  const n = Number(amount)
  if (Number.isNaN(n)) return '—'
  return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(n)
}

// pulls vehicle details off a cart item whether they're nested under
// `item.vehicle` or flattened directly on the item — see note at top
function getVehicleInfo(item) {
  const v = item.vehicle || item
  return {
    name: v.name ?? 'Vehicle',
    brand: v.brand ?? '',
    model: v.model ?? '',
    price: Number(v.price) || 0,
    image_url: v.image_url ?? null,
    stock: v.stock ?? null,
  }
}

function Thumb({ src, alt }) {
  const [failed, setFailed] = useState(false)
  if (!src || failed) {
    return (
      <div className="flex h-20 w-20 flex-none items-center justify-center rounded-xl bg-[#F3F5F8] text-[#5B6472]">
        {NodeIcon('image', 'h-6 w-6')}
      </div>
    )
  }
  return (
    <img
      src={src}
      alt={alt}
      onError={() => setFailed(true)}
      className="h-20 w-20 flex-none rounded-xl border border-slate-100 object-cover"
    />
  )
}

function ClearCartConfirm({ onCancel, onConfirm, isClearing }) {
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
        <h2 className="mt-4 FleetOps-display text-lg font-semibold text-[#0B0E14]">Clear your cart?</h2>
        <p className="mt-2 FleetOps-body text-[13.5px] text-[#5B6472]">
          This removes every item from your cart. This can&rsquo;t be undone.
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
            disabled={isClearing}
            className="flex-1 rounded-full bg-[#DC2626] px-5 py-2.5 FleetOps-body text-[13.5px] font-semibold text-white transition-colors hover:bg-[#B91C1C] disabled:cursor-not-allowed disabled:opacity-50"
          >
            {isClearing ? 'Clearing…' : 'Clear cart'}
          </button>
        </div>
      </div>
    </div>
  )
}

function CartItemRow({ item, onIncrease, onDecrease, onRemove, isUpdating, isRemoving }) {
  const v = getVehicleInfo(item)
  const atMax = v.stock != null && item.quantity >= v.stock
  const subtotal = v.price * item.quantity

  return (
    <div className="flex gap-4 border-b border-slate-100 py-5 last:border-0">
      <Thumb src={v.image_url} alt={v.name} />

      <div className="flex flex-1 flex-col justify-between">
        <div className="flex items-start justify-between gap-3">
          <div>
            <p className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">{v.name}</p>
            <p className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">
              {v.brand} {v.model && `· ${v.model}`}
            </p>
            <p className="mt-1 FleetOps-body text-[13px] text-[#5B6472]">{formatPrice(v.price)} each</p>
          </div>
          <button
            type="button"
            onClick={onRemove}
            disabled={isRemoving}
            className="flex h-8 w-8 flex-none items-center justify-center rounded-lg text-[#5B6472] transition-colors hover:bg-red-50 hover:text-[#DC2626] disabled:cursor-not-allowed disabled:opacity-50"
            aria-label={`Remove ${v.name}`}
          >
            {NodeIcon('trash', 'h-4 w-4')}
          </button>
        </div>

        <div className="mt-3 flex items-center justify-between">
          <div className="flex items-center gap-1 rounded-full border border-slate-200 p-1">
            <button
              type="button"
              onClick={onDecrease}
              disabled={isUpdating || item.quantity <= 1}
              className="flex h-7 w-7 items-center justify-center rounded-full text-[#0B0E14] transition-colors hover:bg-[#F3F5F8] disabled:cursor-not-allowed disabled:opacity-30"
              aria-label="Decrease quantity"
            >
              {NodeIcon('minus', 'h-3.5 w-3.5')}
            </button>
            <span className="w-6 text-center FleetOps-mono text-[13px] font-medium text-[#0B0E14]">{item.quantity}</span>
            <button
              type="button"
              onClick={onIncrease}
              disabled={isUpdating || atMax}
              className="flex h-7 w-7 items-center justify-center rounded-full text-[#0B0E14] transition-colors hover:bg-[#F3F5F8] disabled:cursor-not-allowed disabled:opacity-30"
              aria-label="Increase quantity"
            >
              {NodeIcon('plus', 'h-3.5 w-3.5')}
            </button>
          </div>

          <p className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">{formatPrice(subtotal)}</p>
        </div>
        {atMax && <p className="mt-1 FleetOps-mono text-[10.5px] text-[#B45309]">Max available stock reached</p>}
      </div>
    </div>
  )
}

export default function Cart() {
  const { data: cartData, isLoading, error } = useGetUserCartQuery()
  const { data: totalData } = useGetCartTotalQuery()
  const { data: countData } = useCountItemsQuery()

  const [updateCartItem] = useUpdateCartItemMutation()
  const [deleteCartItem] = useDeleteCartItemMutation()
  const [clearCart, { isLoading: isClearing }] = useClearCartMutation()

  const [updatingId, setUpdatingId] = useState(null)
  const [removingId, setRemovingId] = useState(null)
  const [confirmClear, setConfirmClear] = useState(false)

  // --- unwrap the response — see note at top of file -------------------------
  const items = cartData?.data ?? (Array.isArray(cartData) ? cartData : []) ?? []
  const computedTotal = items.reduce((sum, it) => sum + getVehicleInfo(it).price * it.quantity, 0)
  const computedCount = items.reduce((sum, it) => sum + it.quantity, 0)

  const total = totalData?.total ?? totalData?.data?.total ?? computedTotal
  const itemCount = countData?.count ?? countData?.data?.count ?? computedCount

  async function handleIncrease(item) {
    setUpdatingId(item.id)
    try {
      await updateCartItem({ itemId: item.id, quantity: item.quantity + 1 }).unwrap()
    } finally {
      setUpdatingId(null)
    }
  }

  async function handleDecrease(item) {
    if (item.quantity <= 1) return
    setUpdatingId(item.id)
    try {
      await updateCartItem({ itemId: item.id, quantity: item.quantity - 1 }).unwrap()
    } finally {
      setUpdatingId(null)
    }
  }

  async function handleRemove(item) {
    setRemovingId(item.id)
    try {
      await deleteCartItem(item.id).unwrap()
    } finally {
      setRemovingId(null)
    }
  }

  async function handleClearCart() {
    try {
      await clearCart().unwrap()
    } finally {
      setConfirmClear(false)
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

      <main className="mx-auto max-w-5xl px-6 py-10">
        <div className="flex flex-col gap-1 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <span className="FleetOps-mono text-[11px] tracking-wide text-[#5B6472]">cart</span>
            <h1 className="mt-1 FleetOps-display text-2xl font-semibold tracking-tight text-[#0B0E14] sm:text-3xl">
              Your cart
            </h1>
          </div>
          {items.length > 0 && (
            <p className="FleetOps-mono text-[12px] text-[#5B6472]">{itemCount} item{itemCount === 1 ? '' : 's'}</p>
          )}
        </div>

        {isLoading ? (
          <div className="mt-8 rounded-2xl border border-slate-200 bg-white py-16 text-center">
            <p className="FleetOps-body text-[13.5px] text-[#5B6472]">Loading your cart…</p>
          </div>
        ) : error ? (
          <div className="mt-8 rounded-2xl border border-red-100 bg-red-50 py-16 text-center">
            <p className="FleetOps-body text-[13.5px] text-[#DC2626]">Couldn&rsquo;t load your cart. Please try again.</p>
          </div>
        ) : items.length === 0 ? (
          <div className="mt-8 flex flex-col items-center gap-3 rounded-2xl border border-dashed border-slate-300 bg-white py-20 text-center">
            {NodeIcon('bag', 'h-9 w-9 text-[#5B6472]')}
            <p className="FleetOps-display text-[16px] font-semibold text-[#0B0E14]">Your cart is empty</p>
            <p className="FleetOps-body text-[13.5px] text-[#5B6472]">Vehicles you add will show up here.</p>
          </div>
        ) : (
          <div className="mt-8 grid gap-8 lg:grid-cols-[1fr_320px]">
            <div className="rounded-2xl border border-slate-200 bg-white px-5">
              {items.map((item) => (
                <CartItemRow
                  key={item.id}
                  item={item}
                  isUpdating={updatingId === item.id}
                  isRemoving={removingId === item.id}
                  onIncrease={() => handleIncrease(item)}
                  onDecrease={() => handleDecrease(item)}
                  onRemove={() => handleRemove(item)}
                />
              ))}
            </div>

            <div>
              <div className="rounded-2xl border border-slate-200 bg-white p-6">
                <p className="FleetOps-mono text-[11px] uppercase tracking-wide text-[#5B6472]">Order summary</p>
                <div className="mt-4 flex items-center justify-between">
                  <span className="FleetOps-body text-[13.5px] text-[#5B6472]">Subtotal</span>
                  <span className="FleetOps-body text-[13.5px] font-medium text-[#0B0E14]">{formatPrice(total)}</span>
                </div>
                <div className="mt-4 flex items-center justify-between border-t border-slate-100 pt-4">
                  <span className="FleetOps-display text-[15px] font-semibold text-[#0B0E14]">Total</span>
                  <span className="FleetOps-display text-[18px] font-semibold text-[#0B0E14]">{formatPrice(total)}</span>
                </div>

                <button
                  type="button"
                  className="mt-5 flex w-full items-center justify-center gap-2 rounded-full bg-[#FF5A1F] px-6 py-3 FleetOps-body text-[14px] font-semibold text-white shadow-[0_12px_24px_-8px_rgba(255,90,31,0.55)] transition-all hover:-translate-y-0.5"
                >
                  Proceed to checkout
                  {NodeIcon('arrowRight', 'h-4 w-4')}
                </button>
              </div>

              <button
                type="button"
                onClick={() => setConfirmClear(true)}
                className="mt-3 flex w-full items-center justify-center gap-1.5 rounded-full border border-slate-200 bg-white px-6 py-2.5 FleetOps-body text-[13px] font-medium text-[#5B6472] transition-colors hover:bg-white hover:text-[#DC2626]"
              >
                {NodeIcon('close', 'h-3.5 w-3.5')}
                Clear cart
              </button>
            </div>
          </div>
        )}
      </main>

      {confirmClear && (
        <ClearCartConfirm onCancel={() => setConfirmClear(false)} onConfirm={handleClearCart} isClearing={isClearing} />
      )}
    </div>
  )
}