import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const analyticsApi = createApi({
    reducerPath: "analyticsApi",
    baseQuery: fetchBaseQuery({
        baseUrl: "http://localhost:8087/api/v1", 
    }),

    tagTypes: ["Analytics"],

    endpoints: (builder) => ({
        getPaymentAnalytic: builder.query({
            query: () => ({
                url: "/analytics",
                method: "GET"
            }),
            providesTags: ["Analytics"],
        }),
    })
})

export const { useGetPaymentAnalyticQuery } = analyticsApi;