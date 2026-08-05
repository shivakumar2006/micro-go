import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const vehicleApi = createApi({
    reducerPath: "vehicleApi",
    baseQuery: fetchBaseQuery({
        baseUrl: "http://localhost:8000/api/v1",

        prepareHeaders: (headers, { getState }) => {

            const auth = getState().authReducer;

            if (auth?.accessToken) {
                headers.set("Authorization", `Bearer ${auth.accessToken}`);
            }
        
            if (auth?.refreshToken) {
                headers.set("Refresh-Token", auth.refreshToken);
            }
        
            return headers;
        },
    }),

    endpoints: (builder) => ({
        getAllVehicles: builder.query({
            query: (params) => ({
                url: "/vehicles",
                methods: "GET",
                params,
            }),
        }),

        getVehicleById: builder.query({
            query: (id) => ({
                url: `/vehicles/${id}`,
                methods: "GET"
            })
        }),

        createVehicle: builder.mutation({
            query: (data) => ({
                url: "/vehicles",
                methods: "POST",
                body: data,
            }),
        }),

        updateVehicle: builder.mutation({
            query: ({id, data}) => ({
                url: `/vehicles/${id}`,
                methods: "PATCH",
                body: data,
            }),
        }),

        deleteVehicle: builder.mutation({
            query: (id) => ({
                url: `/vehicles/${id}`,
                methods: "DELETE"
            }),
        }),
    })
})

export const { useGetAllVehiclesQuery, useGetVehicleByIdQuery, useCreateVehicleMutation, useUpdateVehicleMutation, useDeleteVehicleMutation } = vehicleApi;
