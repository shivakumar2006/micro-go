import React from 'react';
import "./App.css";
import {Routes, Route} from "react-router-dom";
import Home from "./Pages/Home";
import Vehicles from './Pages/Vehicles';
import SignIn from './Pages/Signin';
import Signup from "./Pages/Signup";
import Navbar from './components/Navbar';

const App = () => {
  return (
    <>
    <Navbar />
    <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/vehicles" element={<Vehicles />}/>
        <Route path="/login" element={<SignIn />} />
        <Route path="/register" element={<Signup />} />
    </Routes>
    </>
  )
}

export default App