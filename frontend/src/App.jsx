import React from 'react';
import "./App.css";
import {Routes, Route} from "react-router-dom";
import Vehicles from './Pages/Vehicles';

const App = () => {
  return (
    <>
    <Routes>
        <Route path="/" element={<Vehicles />}/>
    </Routes>
    </>
  )
}

export default App