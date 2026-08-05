import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const authApi = createApi({
    reducerPath: "AuthApi",
    baseQuery: fetchBaseQuery({
        baseUrl: "http://localhost:8000/api/v1/auth",

        prepareHeaders: (headers, {getState}) => {
            const { accessToken } = getState().authReducer;
            if (accessToken) {
                headers.set("Authorization", `Bearer ${accessToken}`)
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
                query: (refreshToken) => ({
                    url: "/logout",
                    method: "POST",
                    body: {
                        refresh_token: refreshToken,
                    }
                }),
            }),

            logoutAll: builder.mutation({
                query: (refreshToken) => ({
                    url: "/logout-all",
                    method: "POST",
                    body: {
                        refresh_token: refreshToken,
                    }
                }),
            }),

            refreshToken: builder.mutation({
                query: (refreshToken) => ({
                    url: "/refresh",
                    method: "POST",
                    body: {
                        refresh_token: refreshToken,
                    }
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