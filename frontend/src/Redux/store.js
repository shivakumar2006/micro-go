import { configureStore } from "@reduxjs/toolkit";
import {AuthApi} from "./features/auth/authApi";

export const store = configureStore({
    reducer: {
        auth: authReducer,
        [AuthApi.reducerPath]: AuthApi.reducer,
    },

    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(AuthApi.middleware)
})