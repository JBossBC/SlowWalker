import React, { useEffect } from 'react';
import axios from 'axios';

const setHeadersToken = (token) => {
    if (token) {
        // 将JWT令牌添加到请求头部
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
        // 如果没有提供令牌，则删除请求头部中的Authorization字段
        delete axios.defaults.headers.common['Authorization'];
    }
};

const JWTHeadersRequest = () => {
    useEffect(() => {
        const jwtToken = sessionStorage.getItem('repliteweb');
        setHeadersToken(jwtToken);
    }, []);

    return <></>;
};

export default JWTHeadersRequest;
