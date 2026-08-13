import { createApi } from "@reduxjs/toolkit/query/react";
import { BaseQueryWithReAuth } from "../../BaseWithReAuth";

export const orderApi = createApi({
    reducerPath: "orderApi",
    baseQuery: BaseQueryWithReAuth,

    tagTypes: ["orders"],

    endpoints: (builder) => ({
        getOrderById: builder.query({
            query: (id) => ({
                url: `/orders/${id}`,
                method: "GET",
            }),
            providesTags: ["orders"],
        }),

        getOrdersByUserId: builder.query({
            query: (userId) => ({
                url: `/orders/user/${userId}`,
                method: "GET",
            }),
            providesTags: ["orders"],
        }),


        createOrder: builder.mutation({
            query: (order) => ({
                url: "/orders",
                method: "POST",
                body: order,
            }),
            invalidatesTags: ["orders"],
        }),

        updateOrderStatus: builder.mutation({
            query: ({ id, ...body }) => ({
                url: `/orders/${id}/status`,
                method: "PATCH",
                body,
            }),
            invalidatesTags: ["orders"],
        }),

        cancelOrder: builder.mutation({
            query: (id) => ({
                url: `/orders/${id}/cancel`,
                method: "PATCH",
            }),
            invalidatesTags: ["orders"],
        }),

        markOrderPaid: builder.mutation({
            query: ({ id, ...body }) => ({
                url: `/orders/${id}/pay`,
                method: "PATCH",
                body,
            }),
            invalidatesTags: ["orders"],
        }),
    }),
});

export const {
    useGetOrderByIdQuery,
    useGetOrdersByUserIdQuery,

    useCreateOrderMutation,
    useUpdateOrderStatusMutation,
    useCancelOrderMutation,
    useMarkOrderPaidMutation,
} = orderApi;
