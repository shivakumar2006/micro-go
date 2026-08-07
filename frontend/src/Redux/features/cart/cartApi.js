import { createApi } from "@reduxjs/toolkit/query/react";
import { BaseQueryWithReAuth } from "../../BaseWithReAuth";

export const cartApi = createApi({
    reducer: "cartApi",
    baseQuery: BaseQueryWithReAuth,

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
        }),

        getUserCart: builder.query({
            query: () => ({
                url: "/cart",
                method: "GET",
            })
        }),

        updateCartItem: builder.mutation({
            query: ({itemId, quantity}) => ({
                url: `/cart/${itemId}`,
                method: "PATCH",
                body: {
                    quantity: quantity,
                },
            }),
        }),

        deleteCartItem: builder.mutation({
            query: (itemId) => ({
                url: `/cart/${itemId}`,
                method: "DELETE",
            }),
        }),

        clearCart: builder.mutation({
            query: () => ({
                url: "/cart",
                method: "DELETE",
            }),
        }),

        getCartTotal: builder.query({
            query: () => ({
                url: "/cart/total",
                method: "GET",
            }),
        }),

        countItems: builder.query({
            query: () => ({
                url: "/cart/count",
                method: "GET",
            }),
        }),
    })
})

export const { useAddToCartMutation, useGetUserCartQuery, useUpdateCartItemMutation, useDeleteCartItemMutation, useClearCartMutation, useGetCartTotalQuery, useCountItemsQuery } = cartApi;