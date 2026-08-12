import { createApi } from "@reduxjs/toolkit/query/react";
import { BaseQueryWithReAuth } from "../../BaseWithReAuth";

export const paymentApi = createApi({
    reducerPath: "paymentApi",
    baseQuery: BaseQueryWithReAuth,

    tagTypes: ["payments"],

    endpoints: (builder) => ({
        createPayment: builder.mutation({
            query: (payment) => ({
                url: "/payments/create-checkout-session",
                method: "POST",
                body: payment,
            }),
            invalidatesTags: ["payments"],
        }),

        getPaymentById: builder.query({
            query: (id) => ({
                url: `/payments/${id}`,
                method: "GET",
            }),
            providesTags: ["payments"],
        }),

        getPaymentByOrderId: builder.query({
            query: (orderId) => ({
                url: `/payments/order/${orderId}`,
                method: "GET",
            }),
            providesTags: ["payments"],
        }),
    }),
});

export const {
    useCreatePaymentMutation,
    useGetPaymentByIdQuery,
    useGetPaymentByOrderIdQuery,
} = paymentApi;