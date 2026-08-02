import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./features/auth/authSlice";
import { authApi } from "./features/auth/authApi";

export const store = configureStore({
    reducer: {
        authReducer: authReducer,
        [authApi.reducerPath]: authApi.reducer,
    },

    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(authApi.middleware)
})