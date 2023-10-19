import {BrowserRouter as Router,Route,Routes,useNavigate, BrowserRouter}  from "react-router-dom"
import './App.css';
import axios from "axios";
import React, { useEffect, useState } from "react";
import { NotFound,Index,Register, Main } from "./components";
import { message } from "antd";
// init the axios interceptors from error handle
const defaultToken = sessionStorage.getItem("repliteweb")!=undefined?sessionStorage.getItem("repliteweb"):"";
axios.defaults.headers.common["Authorization"] = `Bearer ${defaultToken}`;
const defaultBackendURL = "http://localhost:8080";
export const Backend = React.createContext(defaultBackendURL);



function App(){
    // console.log(navigate);
    const {Token,setToken} = useState(defaultToken);
    // const navigate =useNavigate();
    useEffect(()=>{
        axios.interceptors.response.use(null,(error)=>{
            // if (error.response.state ==304){
            setToken("");
            sessionStorage.removeItem("repliteweb");
            // navigate("/login");
            // }
            message.error("系统出错啦.....");
            console.log(error);
            return Promise.reject(error);
        })
    },[])
    useEffect(()=>{
        axios.defaults.headers.common["Authorization"] = `Bearer ${Token}`
    },[Token]);
    return(
        // <RouterProvider router={routers}/>
        <BrowserRouter>
            <Backend.Provider value={defaultBackendURL}>
                <Routes>
                    <Route path="/" element={<Index Token={Token} setToken={setToken} />}/>
                    <Route path="/login"  index element={<Index Token={Token} setToken={setToken}/>}/>
                    <Route path="/register" element={<Register/>}/>
                    <Route path="/main" element={<Main/>}/>
                    <Route path="*" element={<NotFound/>}/>
                </Routes>
            </Backend.Provider>
        </BrowserRouter>
    )

}

export default App;
