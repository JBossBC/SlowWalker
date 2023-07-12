import {BrowserRouter as Router,Route,Routes}  from "react-router-dom"
import './App.css';

import React from "react";
import { NotFound,Index,Register } from "./components";
function App(){
  return(
    <Router>
      <Routes>
      <Route path="/" Component={Index}/>
      <Route path="/login" Component={Index}/>
      <Route path="/register" Component={Register}></Route>
       <Route path="*" Component={NotFound}/>
       </Routes>
    </Router>
  )
}

export default App;
