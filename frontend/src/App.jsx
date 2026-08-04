import React from 'react';
import "./App.css";
import {Routes, Route, useLocation} from "react-router-dom";
import Home from "./Pages/Home";
import Vehicles from './Pages/Vehicles';
import SignIn from './Pages/Signin';
import Signup from "./Pages/Signup";
import Navbar from './components/Navbar';
import Navbar2 from "./components/Navbar2";
import VehicleDetail from './Pages/VechileDetail';
import AdminVehicles from './Pages/AdminVehicles';
import AdminNavbar from './components/AdminNavbar';

const App = () => {

  const location = useLocation();

  const hideNavbar2 = [
    "/login", "/register",
  ].includes(location.pathname) || location.pathname.startsWith("/vehicles/admin");

  const hideNavbar = [
  "/",
  "/login",
  "/register",
  "/vehicles",
].includes(location.pathname) || location.pathname.startsWith("/vehicles/details") || location.pathname.startsWith("/vehicles/admin");

  const adminNavbar = [
    "/", "/vehicles", "/login", "/register",
  ].includes(location.pathname) || location.pathname.startsWith("/vehicles/details/");

  return (
    <>
    {!hideNavbar && <Navbar />}
    {!hideNavbar2 && <Navbar2 />}
    {!adminNavbar && <AdminNavbar />}
    <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/vehicles" element={<Vehicles />}/>
        <Route path="/login" element={<SignIn />} />
        <Route path="/register" element={<Signup />} />
        <Route path="/vehicles/details/:id" element={<VehicleDetail />} />
        <Route path="/vehicles/admin" element={<AdminVehicles />} />
    </Routes>
    </>
  )
}

export default App