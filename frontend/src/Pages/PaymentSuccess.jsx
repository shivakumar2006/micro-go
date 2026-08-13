import { CheckCircle2, ArrowRight } from "lucide-react";
import { useNavigate, useSearchParams } from "react-router-dom";

export default function PaymentSuccess() {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();

    const orderId = searchParams.get("order_id");

    return (
        <div className="min-h-screen bg-[#F7F8FA] flex items-center justify-center px-5">
            <div className="w-full max-w-[500px]">

                <div className="rounded-[28px] border border-slate-200 bg-white px-8 py-10 shadow-[0_20px_60px_-25px_rgba(15,23,42,0.25)]">

                    {/* Success Icon */}
                    <div className="flex justify-center">
                        <div className="flex h-20 w-20 items-center justify-center rounded-full bg-emerald-50">
                            <CheckCircle2
                                className="h-12 w-12 text-emerald-500"
                                strokeWidth={1.8}
                            />
                        </div>
                    </div>

                    {/* Content */}
                    <div className="mt-7 text-center">

                        <p className="FleetOps-body text-xs font-semibold uppercase tracking-[0.18em] text-emerald-600">
                            Payment Successful
                        </p>

                        <h1 className="FleetOps-heading mt-2 text-[32px] font-semibold tracking-tight text-[#0B0E14]">
                            Thank you for your purchase!
                        </h1>

                        <p className="FleetOps-body mx-auto mt-3 max-w-[380px] text-sm leading-6 text-[#687386]">
                            Your payment has been successfully processed and
                            your order has been placed.
                        </p>
                    </div>

                    {/* Order ID */}
                    {orderId && (
                        <div className="mt-7 rounded-2xl bg-[#F8F9FB] border border-slate-100 px-5 py-4 text-center">
                            <p className="FleetOps-body text-xs text-[#8A93A3]">
                                Order ID
                            </p>

                            <p className="FleetOps-body mt-1 text-lg font-semibold text-[#0B0E14]">
                                #{orderId}
                            </p>
                        </div>
                    )}

                    {/* Button */}
                    <button
                        onClick={() =>
                            orderId
                                ? navigate(`/orders/${orderId}`)
                                : navigate("/orders")
                        }
                        className="group mt-7 flex h-12 w-full items-center justify-center gap-2 rounded-xl bg-[#0B0E14] FleetOps-body text-sm font-semibold text-white transition-all hover:bg-[#171B24] active:scale-[0.99]"
                    >
                        View My Order

                        <ArrowRight
                            className="h-4 w-4 transition-transform group-hover:translate-x-1"
                        />
                    </button>

                    <button
                        onClick={() => navigate("/vehicles")}
                        className="mt-3 h-12 w-full rounded-xl border border-slate-200 bg-white FleetOps-body text-sm font-semibold text-[#0B0E14] transition-colors hover:bg-[#F5F6F8]"
                    >
                        Continue Shopping
                    </button>

                    <p className="FleetOps-body mt-7 text-center text-xs text-[#9AA2AF]">
                        Your payment confirmation has been recorded securely.
                    </p>

                </div>

                <p className="FleetOps-body mt-5 text-center text-[11px] text-[#A0A7B3]">
                    © {new Date().getFullYear()} FleetOps
                </p>

            </div>
        </div>
    );
}
