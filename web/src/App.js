import {BrowserRouter as Router,Route,Routes,useNavigate}  from "react-router-dom"
import './App.css';
import axios from "axios";
import React, { useEffect, useState } from "react";

import { NotFound,Index,Register, Main } from "./components";
import { initReactI18next } from 'react-i18next';
import i18n from 'i18next';
import messages from "./utils/translations.json";

// init the axios interceptors from error handle
const defaultToken = sessionStorage.getItem("repliteweb")!=undefined?sessionStorage.getItem("repliteweb"):"";
// axios.defaults.headers.common["Authorization"] = `Bearer ${defaultToken}`;

// const defaultBackendURL = "http://localhost:8080";
// export const Backend = React.createContext(defaultBackendURL);

// i18n.use(initReactI18next).init(
//    {
//     resources:{
        
//     },
//     lng: 'cn',
//     interpolation: {
//         escapeValue: false // react already safes from xss => https://www.i18next.com/translation-function/interpolation#unescape
//       }
//    }
// )
function App(){
    return (
        <Router>
           {/* <I18nextProvider i18n={i18n}> */}
            {/* <Backend.Provider value={defaultBackendURL}> */}
                <Routes>
                    <Route path="/" element={<Index />} /> {/* 使用element属性 */}
                    <Route path="/login" element={<Index/>} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/main" element={<Main />} />
                    <Route path="*" element={<NotFound />} />
                </Routes>
                {/* </I18nextProvider> */}

            {/* </Backend.Provider> */}
        </Router>
    );
}
export default App;
