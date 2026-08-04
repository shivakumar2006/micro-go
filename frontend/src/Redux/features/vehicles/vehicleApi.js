import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const vehicleApi = createApi({
    reducerPath: "vehicleApi",
    baseQuery: fetchBaseQuery({
        baseUrl: "http://localhost:8000/api/v1",

        prepareHeaders: (headers, {getState}) => {
            const { accessToken, refreshToken } = getState.authReducer;

            if (accessToken) {
                headers.set("Authorization", `Bearer ${accessToken}`)
            }
            if (refreshToken) {
                headers.set("Refresh-Token", `${refreshToken}`)
            }
            return headers;
        }
    }),

    endpoints: (builder) => ({
        getAllVehicles: builder.mutation({
            query: () => ({
                url: "/vehicles",
                methods: "GET",

            }),
        }),

        getVehicleById: builder.mutation({
            query: (id) => ({
                url: `/vehicles/${id}`,
                methods: "GET"
            })
        })
    })
})

export const { useGetAllVehiclesMutation, useGetVehicleByIdMutation } = vehicleApi;
