import {BrowserRouter as Router,Route,Routes}  from "react-router-dom"
import './App.css';

import React from "react";
import { NotFound,Index,Register, Main, Log } from "./components";
//const defaultBackendURL = "http://112.124.53.234:8080"
const defaultBackendURL = "http://127.0.0.1:8080"
export const Backend = React.createContext(defaultBackendURL);
function App(){
  return(
    <Router>
     <Backend.Provider value={defaultBackendURL}>
      <Routes>
      <Route path="/" Component={Index}/>
      <Route path="/login" Component={Index}/>
      <Route path="/register" Component={Register}/>
      <Route path="/main" Component={Main}/>
      <Route path="/log" Component={Log}/>
       <Route path="*" Component={NotFound}/>
       </Routes>
       </Backend.Provider>
    </Router>
  )
}

export default App;
