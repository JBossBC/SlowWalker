import {BrowserRouter as Router,Route,Routes,useNavigate}  from "react-router-dom"
import './App.css';
import axios from "axios";
import React, { useEffect, useState } from "react";
<<<<<<< HEAD
import { NotFound,Index,Register, Main, Log } from "./components";
import { message } from "antd";
=======
import { NotFound,Index,Register, Main } from "./components";
import { initReactI18next } from 'react-i18next';
import i18n from 'i18next';
import messages from "./utils/translations.json";
>>>>>>> c378b46238addcc9e928f3c98dba2148025b2edd
// init the axios interceptors from error handle
const defaultToken = sessionStorage.getItem("repliteweb")!=undefined?sessionStorage.getItem("repliteweb"):"";
// axios.defaults.headers.common["Authorization"] = `Bearer ${defaultToken}`;

// const defaultBackendURL = "http://localhost:8080";
// export const Backend = React.createContext(defaultBackendURL);
<<<<<<< HEAD
=======

i18n.use(initReactI18next).init(
   {
    resources:{
        
    },
    lng: 'cn',
    interpolation: {
        escapeValue: false // react already safes from xss => https://www.i18next.com/translation-function/interpolation#unescape
      }
   }
)
>>>>>>> c378b46238addcc9e928f3c98dba2148025b2edd
function App(){

    return (
        <Router>
<<<<<<< HEAD
=======
           {/* <I18nextProvider i18n={i18n}> */}
>>>>>>> c378b46238addcc9e928f3c98dba2148025b2edd
            {/* <Backend.Provider value={defaultBackendURL}> */}
                <Routes>
                    <Route path="/" element={<Index />} /> {/* 使用element属性 */}
                    <Route path="/login" element={<Index/>} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/main" element={<Main />} />
<<<<<<< HEAD
                    <Route path="/log" element={<Log />} />
                    <Route path="*" element={<NotFound />} />
                </Routes>
=======
                    <Route path="*" element={<NotFound />} />
                </Routes>
                {/* </I18nextProvider> */}
>>>>>>> c378b46238addcc9e928f3c98dba2148025b2edd
            {/* </Backend.Provider> */}
        </Router>
    );
}

export default App;
