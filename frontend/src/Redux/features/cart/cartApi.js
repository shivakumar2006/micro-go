import { createApi } from "@reduxjs/toolkit/query/react";
import { BaseQueryWithReAuth } from "../../BaseWithReAuth";

export const cartApi = createApi({
    reducer: "cartApi",
    baseQuery: BaseQueryWithReAuth,

    tagTypes: ["cart"],

    endpoints: (builder) => ({
        addToCart: builder.mutation({
            query: ({ vehicleId, quantity}) => ({
                url: "/cart",
                method: "POST",
                body: {
                    vehicle_id: vehicleId,
                    quantity: quantity,
                },
            }),

            invalidatesTags: ["cart"],
        }),

        getUserCart: builder.query({
            query: (userId) => ({
                url: "/cart",
                method: "GET",
            }),
            providesTags: (result, error, userId) => [
                { type: "cart", id: userId }
            ],
        }),

        updateCartItem: builder.mutation({
            query: ({itemId, quantity}) => ({
                url: `/cart/${itemId}`,
                method: "PUT",
                body: {
                    quantity: quantity,
                },
            }),

            invalidatesTags: ["cart"],
        }),

        deleteCartItem: builder.mutation({
            query: (itemId) => ({
                url: `/cart/${itemId}`,
                method: "DELETE",
            }),
            invalidatesTags: ["cart"],
        }),

        clearCart: builder.mutation({
            query: () => ({
                url: "/cart",
                method: "DELETE",
            }),
            invalidatesTags: ["cart"],
        }),

        getCartTotal: builder.query({
            query: () => ({
                url: "/cart/total",
                method: "GET",
            }),
            providesTags: ["cart"],
        }),

        countItems: builder.query({
            query: () => ({
                url: "/cart/count",
                method: "GET",
            }),
            providesTags: ["cart"],
        }),
    })
})

export const { useAddToCartMutation, useGetUserCartQuery, useUpdateCartItemMutation, useDeleteCartItemMutation, useClearCartMutation, useGetCartTotalQuery, useCountItemsQuery } = cartApi;