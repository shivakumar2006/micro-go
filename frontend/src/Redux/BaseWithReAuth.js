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

    console.log("REQUEST:", args);

    let result = await baseQuery(args, api, extraOptions);

    console.log("FIRST RESULT:", result);
    console.log("ERROR", result.error);
    console.log("STATUS", result.error?.status);

    if (result.error && (result.error.status === 401 || result.error.originalStatus === 401)) {

        console.log("ACCESS TOKEN EXPIRED");

        const refreshToken = api.getState().authReducer.refreshToken;

        console.log("REFRESH TOKEN:", refreshToken);

        const refreshResult = await baseQuery(
            {
                url: "/auth/refresh",
                method: "POST",
                body: {
                    refresh_token: refreshToken,
                },
            },
            api,
            extraOptions
        );

        console.log("REFRESH RESULT:", refreshResult);

        if (refreshResult.data) {

            console.log("NEW TOKEN RECEIVED");

            api.dispatch(
                setTokens({
                    accessToken: refreshResult.data.access_token,
                    refreshToken: refreshResult.data.refresh_token,
                    user: refreshResult.data.user,
                })
            );

            result = await baseQuery(args, api, extraOptions);

            console.log("RETRY RESULT:", result);

        } else {

            console.log("REFRESH FAILED");

            api.dispatch(clearAuth());

        }
    }

    return result;
}

// export const BaseQueryWithReAuth = async (args, api, extraOptions) => {
//     let result = await baseQuery(args, api, extraOptions)

//     // access token expired
//     if (result.error && result.error.status === 401) {
//         const refreshToken = api.getState().authReducer.refreshToken;

//         if (!refreshToken) {
//             api.dispatch(clearAuth());
//             return result;
//         }

//         const refreshResult = await baseQuery({
//             url: "/auth/refresh",
//             method: "POST",
//             body: {
//                 refresh_token: refreshToken,
//             },
//         }, api, extraOptions);

//         if (refreshResult.data) {
//             api.dispatch(
//                 setTokens({
//                     accessToken: refreshResult.data.access_token,
//                     refreshToken: refreshResult.data.refresh_token,
//                     user: refreshResult.data.user,
//                 })
//             );

//             // retry the original request
//             result = await baseQuery(args, api, extraOptions);
//         } else {
//             api.dispatch(clearAuth());
//         }
//     }

//     return result;
// }