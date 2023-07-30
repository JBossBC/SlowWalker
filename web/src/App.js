import {BrowserRouter as Router,Route,Routes,useNavigate}  from "react-router-dom"
import './App.css';
import axios from "axios";
import React, { useEffect, useState } from "react";
import { NotFound,Index,Register, Main } from "./components";
import { message } from "antd";
// init the axios interceptors from error handle
const defaultToken = sessionStorage.getItem("repliteweb")!=undefined?sessionStorage.getItem("repliteweb"):"";
axios.defaults.headers.common["Authorization"] = `Bearer ${defaultToken}`;

const defaultBackendURL = "http://112.124.53.234:8080";
export const Backend = React.createContext(defaultBackendURL);
function App(){
  const navigate =useNavigate();
  
  const {Token,setToken} = useState(defaultToken);
  axios.interceptors.response.use(null,(error)=>{
    if (error.response.state ==304){
      sessionStorage.removeItem("repliteweb");
      navigate("/");
    }
     message.error("系统出错");
     console.log(error);
     return Promise.reject(error);
  })
  useEffect(()=>{
    axios.defaults.headers.common["Authorization"] = `Bearer ${Token}`
  },[Token]);
  return(
    <Router>
     <Backend.Provider value={defaultBackendURL}>
      <Routes>
      <Route path="/" element={Index(Token,setToken)}/>
      <Route path="/login" element={Index(Token,setToken)}/>
      <Route path="/register" element={Register}/>
      <Route path="/main" element={Main}/>
       <Route path="*" element={NotFound}/>
       </Routes>
       </Backend.Provider>
    </Router>
  )
}

export default App;
