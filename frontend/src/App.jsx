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

const App = () => {

  const location = useLocation();

  const hideNavbar2 = [
    "/login", "/register",
  ].includes(location.pathname)

  const hideNavbar = [
  "/",
  "/login",
  "/register",
  "/vehicles",
].includes(location.pathname) || location.pathname.startsWith("/vehicles/details");

  return (
    <>
    {!hideNavbar && <Navbar />}
    {!hideNavbar2 && <Navbar2 />}
    <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/vehicles" element={<Vehicles />}/>
        <Route path="/login" element={<SignIn />} />
        <Route path="/register" element={<Signup />} />
        <Route path="/vehicles/details/:id" element={<VehicleDetail />} />
    </Routes>
    </>
  )
}

export default App