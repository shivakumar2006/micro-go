import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./features/auth/authSlice";
import { authApi } from "./features/auth/authApi";
import { vehicleApi } from "./features/vehicles/vehicleApi";

export const store = configureStore({
    reducer: {
        authReducer: authReducer,
        [authApi.reducerPath]: authApi.reducer,
        [vehicleApi.reducerPath]: vehicleApi.reducer,
    },

    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(authApi.middleware, vehicleApi.middleware)
})