import {BrowserRouter as Router,Route,Routes}  from "react-router-dom"
import './App.css';

import React from "react";
import { NotFound,Index,Register, Main } from "./components";
const defaultBackendURL = "http://112.124.53.234:8080"
export const Backend = React.createContext(defaultBackendURL);
function App(){
  return(
    <Router>
     <Backend.Provider value={defaultBackendURL}>
      <Routes>
      <Route path="/" Component={Index}/>
      <Route path="/login" Component={Index}/>
      <Route path="/register" Component={Register}></Route>
      <Route path="/main" Component={Main}></Route>
       <Route path="*" Component={NotFound}/>
       </Routes>
       </Backend.Provider>
    </Router>
  )
}

export default App;
