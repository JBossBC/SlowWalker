import {BrowserRouter as Router,Route,Routes,useNavigate}  from "react-router-dom"
import './App.css';
import axios from "axios";
import React, { useEffect, useState } from "react";
import { NotFound,Index,Register, Main } from "./components";
import { message } from "antd";
// init the axios interceptors from error handle
const defaultToken = sessionStorage.getItem("repliteweb")!=undefined?sessionStorage.getItem("repliteweb"):"";
// axios.defaults.headers.common["Authorization"] = `Bearer ${defaultToken}`;

// const defaultBackendURL = "http://localhost:8080";
// export const Backend = React.createContext(defaultBackendURL);
function App(){

    return (
        <Router>
            {/* <Backend.Provider value={defaultBackendURL}> */}
                <Routes>
                    <Route path="/" element={<Index />} /> {/* 使用element属性 */}
                    <Route path="/login" element={<Index/>} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/main" element={<Main />} />
                    <Route path="*" element={<NotFound />} />
                </Routes>
            {/* </Backend.Provider> */}
        </Router>
    );
}

export default App;
