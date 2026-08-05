import { createApi } from "@reduxjs/toolkit/query/react";
import { BaseQueryWithReAuth } from "../../BaseWithReAuth";

export const authApi = createApi({
    reducerPath: "AuthApi",
    baseQuery: BaseQueryWithReAuth,

    endpoints: (builder) => ({
            register: builder.mutation({
                query: (body) => ({
                    url: "/auth/register",
                    method: "POST",
                    body,
                }),
            }),

            login: builder.mutation({
                query: (body) => ({
                    url: "/auth/login",
                    method: "POST",
                    body,
                }),
            }),

            logout: builder.mutation({
                query: (refreshToken) => ({
                    url: "/auth/logout",
                    method: "POST",
                    body: {
                        refresh_token: refreshToken,
                    }
                }),
            }),

            logoutAll: builder.mutation({
                query: (refreshToken) => ({
                    url: "/auth/logout-all",
                    method: "POST",
                    body: {
                        refresh_token: refreshToken,
                    }
                }),
            }),

            refreshToken: builder.mutation({
                query: (refreshToken) => ({
                    url: "/auth/refresh",
                    method: "POST",
                    body: {
                        refresh_token: refreshToken,
                    }
                })
            }),

            me: builder.query({
                query: () => ({
                    url: "/auth/me",
                    method: "GET",
                }),
            })
    })
})

export const {useRegisterMutation, useLoginMutation, useLogoutMutation, useLogoutAllMutation, useRefreshTokenMutation, useMeQuery} = authApi