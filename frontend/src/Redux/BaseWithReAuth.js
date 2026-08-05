import { fetchBaseQuery } from "@reduxjs/toolkit/query";
import { clearAuth, setTokens } from "./features/auth/authSlice";

const baseQuery = fetchBaseQuery({
    baseUrl: "http://localhost:8000/api/v1",

    prepareHeaders: (headers, { getState }) => {
        const { accessToken } = getState().authReducer;

        if (accessToken) {
            headers.set("Authorization", `Bearer ${accessToken}`);
        }
        return headers;
    },
})

export const BaseQueryWithReAuth = async (args, api, extraOptions) => {
    let result = await baseQuery(args, api, extraOptions)

    // access token expired
    if (result.error && result.error.status === 401) {
        const refreshToken = api.getState().authReducer.refreshToken;

        if (!refreshToken) {
            api.dispatch(clearAuth());
            return result;
        }

        const refreshResult = await baseQuery({
            url: "/auth/refresh",
            method: "POST",
            body: {
                refresh_token: refreshToken,
            },
        }, api, extraOptions);

        if (refreshResult.data) {
            api.dispatch(
                setTokens({
                    accessToken: refreshResult.data.access_token,
                    refreshToken: refreshResult.data.refresh_token,
                    user: refreshResult.data.user,
                })
            );

            // retry the original request
            result = await baseQuery(args, api, extraOptions);
        } else {
            api.dispatch(clearAuth());
        }
    }

    return result;
}