import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const authApi = createApi({
    reducerPath: "AuthApi",
    baseQuery: fetchBaseQuery({
        baseURL: "http://localhost:8000/api/v1/auth",

        prepareHeaders: (headers, {getState}) => {
            const { accessToken, refreshToken } = getState().auth;
            if (accessToken) {
                headers.set("Authorization", `Bearer ${accessToken}`)
            }

            if (refreshToken) {
                headers.set("Refresh-Token", `${refreshToken}`)
            }

            return headers;
        },
    }),

    endpoints: (builder) => ({
            register: builder.mutation({
                query: (body) => ({
                    url: "/register",
                    method: "POST",
                    body,
                }),
            }),

            login: builder.mutation({
                query: (body) => ({
                    url: "/login",
                    method: "POST",
                    body,
                }),
            }),

            logout: builder.mutation({
                query: () => ({
                    url: "/logout",
                    method: "POST",
                }),
            }),

            logoutAll: builder.mutation({
                query: () => ({
                    url: "/logout-all",
                    method: "POST",
                }),
            }),

            refreshToken: builder.mutation({
                query: () => ({
                    url: "/refresh",
                    method: "POST",
                })
            }),

            me: builder.query({
                query: () => ({
                    url: "/me",
                    method: "GET",
                }),
            })
    })
})

export const {useRegisterMutation, useLoginMutation, useLogoutMutation, useLogoutAllMutation, useRefreshTokenMutation, useMeQuery} = authApi