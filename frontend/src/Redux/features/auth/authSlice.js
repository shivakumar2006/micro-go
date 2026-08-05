import { createSlice } from "@reduxjs/toolkit"
import { authApi } from "./authApi"

const initialState = {
    User: null,
    accessToken: null,
    refreshToken: null,
    role: null,
    isAuthenticated: false,
}

const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        clearAuth: (state) => {
            state.User = null
            state.accessToken = null
            state.refreshToken = null
            state.role = null
            state.isAuthenticated = false
        },

        setTokens: (state, action) => {
            state.accessToken = action.payload.accessToken
            state.refreshToken = action.payload.refreshToken
            state.role = action.payload.role
        },
    },

    extraReducers: (builder) => {
        builder.addMatcher(
            authApi.endpoints.register.matchFulfilled,
            (state, action) => {
                state.User = action.payload.User
                state.accessToken = action.payload.access_token
                state.refreshToken = action.payload.refresh_token
                state.isAuthenticated = true
            }
        )

        builder.addMatcher(
            authApi.endpoints.login.matchFulfilled,
            (state, action) => {
                state.User = action.payload.User
                state.accessToken = action.payload.access_token
                state.refreshToken = action.payload.refresh_token
                state.role = action.payload.role
                state.isAuthenticated = true 
            }
        )

        builder.addMatcher(
            authApi.endpoints.refreshToken.matchFulfilled,
            (state, action) => {
                state.User = action.payload.User
                state.accessToken = action.payload.access_token
                state.refreshToken = action.payload.refresh_token
                state.role = action.payload.role
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
                state.User = action.payload
                state.role = action.payload
                state.isAuthenticated = true
            }
        )
    }
})

export default authSlice.reducer;
export const { clearAuth, setTokens } = authSlice.actions;