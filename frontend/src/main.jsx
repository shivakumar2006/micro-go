import React from "react"; 
import ReactDOM from "react-dom/client"
import App from "./App.jsx";
import { BrowserRouter } from "react-router-dom";
import { Provider } from "react-redux";
import { store } from "./Redux/store";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

const root = ReactDOM.createRoot(document.getElementById("root"))

root.render(
    <Provider store={store}>
        <BrowserRouter>
            <ToastContainer />
            <App />
        </BrowserRouter>
    </Provider>
)