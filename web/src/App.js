import {BrowserRouter as Router,Route,Routes}  from "react-router-dom"
import './App.css';

import React from "react";
import { NotFound,Index,Register } from "./components";
const backend = React.createContext("http://localhost:8080/");
function App(){
  return(
    <Router>
      <Routes>
        <backend.Provider>
      <Route path="/" Component={Index}/>
      <Route path="/login" Component={Index}/>
      <Route path="/register" Component={Register}></Route>
       <Route path="*" Component={NotFound}/>
       </backend.Provider>
       </Routes>
    </Router>
  )
}

export default App;
