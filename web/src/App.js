import {BrowserRouter as Router,Route,Routes,useNavigate}  from "react-router-dom"
import './App.css';
import axios from "axios";
import React, { useEffect, useState } from "react";

import { NotFound,Index,Register, Main } from "./components";
import { message } from "antd";
import { IntlProvider, useIntl,FormattedMessage } from "react-intl";
import messages from "./utils/language";
import { Select } from "antd";
// init the axios interceptors from error handle
const defaultToken = sessionStorage.getItem("repliteweb")!=undefined?sessionStorage.getItem("repliteweb"):"";
// axios.defaults.headers.common["Authorization"] = `Bearer ${defaultToken}`;

// const defaultBackendURL = "http://localhost:8080";
// export const Backend = React.createContext(defaultBackendURL);
const { Option } = Select;

function App() {
    const [locale, setLocale] = useState('zh');

    const handleLocaleChange = (lang) => {
        setLocale(lang);
    };

    useEffect(() => {
        setLocale(window.navigator.language.toLowerCase().slice(0, 2));
    }, []);

    const currentMessages = messages[locale];

    return (
        <IntlProvider locale={locale} messages={currentMessages}>
            <Router>
                <div className="container">
                    <div style={{ position: "fixed", top: "20px", right: "20px", zIndex: 1 }}>
                        <Select defaultValue={locale} onChange={handleLocaleChange}>
                            <Option value="en">
                                <FormattedMessage id="English" />
                            </Option>
                            <Option value="zh">
                                <FormattedMessage id="中文" />
                            </Option>
                        </Select>
                    </div>
                    <Routes>
                        <Route path="/" element={<Index />} />
                        <Route path="/login" element={<Index />} />
                        <Route path="/register" element={<Register />} />
                        <Route path="/main" element={<Main />} />
                        <Route path="*" element={<NotFound />} />
                    </Routes>
                </div>
            </Router>
        </IntlProvider>

    );
}
export default App;
