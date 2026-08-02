import { createSlice } from "@reduxjs/toolkit"
import { authApi } from "./authApi"

const initialState = {
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
}

const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        clearAuth: (state) => {
            state.user = null
            state.accessToken = null
            state.refreshToken = null
            state.isAuthenticated = false
        },

        setTokens: (state, action) => {
            state.accessToken = action.payload.accessToken
            state.refreshToken = action.payload.refreshToken
        },
    },

    extraReducers: (builder) => {
        builder.addMatcher(
            authApi.endpoints.register.matchFulfilled,
            (state, action) => {
                state.user = action.payload.user
                state.accessToken = action.payload.access_token
                state.refreshToken = action.payload.refresh_token
                state.isAuthenticated = true
            }
        )

        builder.addMatcher(
            authApi.endpoints.login.matchFulfilled,
            (state, action) => {
                state.user = action.payload.user
                state.accessToken = action.payload.access_token
                state.refreshToken = action.payload.refresh_token
                state.isAuthenticated = true 
            }
        )

        builder.addMatcher(
            authApi.endpoints.refreshToken.matchFulfilled,
            (state, action) => {
                state.user = action.payload.user
                state.accessToken = action.payload.access_token
                state.refreshToken = action.payload.refresh_token
                state.isAuthenticated = true
            }
        )

        builder.addMatcher(
            authApi.endpoints.logout.matchFulfilled,
            (state) => {
                authSlice.caseReducers.clearAuth(state)
            }
        )

        builder.addMatcher(
            authApi.endpoints.logoutAll.matchFulfilled,
            (state) => {
                authSlice.caseReducers.clearAuth(state)
            }
        )

        builder.addMatcher(
            authApi.endpoints.me.matchFulfilled,
            (state, action) => {
                state.user = action.payload
                state.isAuthenticated = true
            }
        )
    }
})

export default authSlice.reducer;
export const { clearAuth, setTokens } = authSlice.actions;