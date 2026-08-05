import { createApi } from "@reduxjs/toolkit/query/react";
import { BaseQueryWithReAuth } from "../../BaseWithReAuth";

export const vehicleApi = createApi({
    reducerPath: "vehicleApi",
    baseQuery: BaseQueryWithReAuth,

    tagTypes: ["Vehicle"],

    endpoints: (builder) => ({
        getAllVehicles: builder.query({
            query: (params) => ({
                url: "/vehicles",
                method: "GET",
                params,
            }),

            providesTags: ["Vehicle"],
        }),

        getVehicleById: builder.query({
            query: (id) => ({
                url: `/vehicles/${id}`,
                method: "GET"
            })
        }),

        createVehicle: builder.mutation({
            query: (data) => ({
                url: "/vehicles",
                method: "POST",
                body: data,
            }),
            invalidatesTags: ["Vehicle"],
        }),

        updateVehicle: builder.mutation({
            query: ({id, data}) => ({
                url: `/vehicles/${id}`,
                method: "PUT",
                body: data,
            }),
            invalidatesTags: ["Vehicle"]
        }),

        deleteVehicle: builder.mutation({
            query: (id) => ({
                url: `/vehicles/${id}`,
                method: "DELETE"
            }),
            invalidatesTags: ["Vehicle"]
        }),
    })
})

export const { useGetAllVehiclesQuery, useGetVehicleByIdQuery, useCreateVehicleMutation, useUpdateVehicleMutation, useDeleteVehicleMutation } = vehicleApi;
