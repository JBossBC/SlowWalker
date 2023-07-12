
import {Result,Button} from  "antd";

import React from "react";
import {useNavigate} from "react-router-dom"
const NotFound=()=>{
    const navigate =useNavigate();
    return (<Result
    status="404"
    title="404"
    subTitle="抱歉,你访问的页面不存在."
    extra={<Button type="primary" onClick={()=>navigate("/")}>返回主页</Button>}
  />)
} 


export default NotFound