import React from "react";
import "./App.css";
import { Routes, Route, useLocation } from "react-router-dom";
import { useSelector } from "react-redux";

import Home from "./Pages/Home";
import Vehicles from "./Pages/Vehicles";
import VehicleDetail from "./Pages/VechileDetail";
import SignIn from "./Pages/Signin";
import Signup from "./Pages/Signup";
import AdminVehicles from "./Pages/AdminVehicles";

import Navbar from "./components/Navbar";
import Navbar2 from "./components/Navbar2";
import AdminNavbar from "./components/AdminNavbar";
import ProtectedRoutes from "./components/ProtectedRoutes";
import Cart from "./Pages/Cart";

const App = () => {
  const location = useLocation();

  const { isAuthenticated, User } = useSelector(
    (state) => state.authReducer
  );

  const role = User?.role_name;

  // Route Checks
  const isHome = location.pathname === "/";
  const isLogin = location.pathname === "/login";
  const isRegister = location.pathname === "/register";
  const isVehicles = location.pathname === "/vehicles";
  const isVehicleDetails = location.pathname.startsWith("/vehicles/details");
  const isAdminPage = location.pathname.startsWith("/vehicles/admin");
  const isCart = location.pathname === "/cart";

  return (
    <>
      {/* Home Page -> Always Navbar2 */}
      {isHome && <Navbar2 />}

      {/* Guest Navbar */}
      {!isAuthenticated &&
        !isHome &&
        !isLogin &&
        !isRegister &&
        !isAdminPage && <Navbar2 />}

      {/* Customer Navbar */}
      {isAuthenticated &&
        role === "customer" &&
        (isVehicles || isVehicleDetails || isCart) && <Navbar />}

      {/* Admin Navbar */}
      {isAuthenticated &&
        role === "admin" &&
        isAdminPage && <AdminNavbar />}

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<SignIn />} />
        <Route path="/register" element={<Signup />} />

        <Route path="/vehicles" element={<Vehicles />} />
        <Route path="/vehicles/details/:id" element={<VehicleDetail />} />
        <Route path="/cart" element={<Cart />} />

        <Route element={<ProtectedRoutes allowRoles={["admin"]} />}>
          <Route path="/vehicles/admin" element={<AdminVehicles />} />
        </Route>
      </Routes>
    </>
  );
};

export default App;