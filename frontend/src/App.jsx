import React from 'react';
import "./App.css";
import {Routes, Route} from "react-router-dom";
import Home from "./Pages/Home";
import Vehicles from './Pages/Vehicles';
import SignIn from './Pages/Signin';

const App = () => {
  return (
    <>
    <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/vehicles" element={<Vehicles />}/>
        <Route path="/login" element={<SignIn />} />
    </Routes>
    </>
  )
}

export default App