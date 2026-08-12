import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./features/auth/authSlice";
import { authApi } from "./features/auth/authApi";
import { vehicleApi } from "./features/vehicles/vehicleApi";
import { persistReducer, persistStore } from "redux-persist"
import storageImport from "redux-persist/lib/storage";
import { cartApi } from "./features/cart/cartApi";
import { orderApi } from "./features/order/order";

const storage = storageImport.default ?? storageImport;

const persistConfig = {
    key: "auth",
    storage,
    whitelist: ["user", "accessToken", "refreshToken", "role", "isAuthenticated"],
};

const persistedAuthReducer = persistReducer(persistConfig, authReducer);

export const store = configureStore({
    reducer: {
        authReducer: persistedAuthReducer,
        [authApi.reducerPath]: authApi.reducer,
        [vehicleApi.reducerPath]: vehicleApi.reducer,
        [cartApi.reducerPath]: cartApi.reducer,
        [orderApi.reducerPath]: orderApi.reducer,
    },

    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware({
            serializableCheck: {
                ignoredActions: [
                    "persist/PERSIST",
                    "persist/REHYDRATE",
                    "persist/PAUSE",
                    "persist/FLUSH",
                    "persist/PURGE",
                    "persist/REGISTER",
                ]
            }
        }).concat(authApi.middleware, vehicleApi.middleware, cartApi.middleware, orderApi.middleware)
})

export const persistor = persistStore(store);