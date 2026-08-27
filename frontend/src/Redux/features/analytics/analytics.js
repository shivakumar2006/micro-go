import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const analyticsApi = createApi({
    reducerPath: "analyticsApi",
    baseQuery: fetchBaseQuery({
        baseUrl: "http://localhost:8087/api/v1", 
    }),

    endpoints: (builder) => ({
        getPaymentAnalytic: builder.query({
            query: () => ({
                url: "/analytics/payment"
            })
        }),
    })
})