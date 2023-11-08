import { Form, Input, Checkbox, Button, message,Card } from "antd";
import { LockOutlined, UserOutlined } from "@ant-design/icons";
import {useNavigate}  from 'react-router-dom';
import logo from "../public/logo.png";
import React, { useState,useContext,useEffect} from "react";
import axios from "../utils/axios";
// import axios from "../utils/axios";
//import { useHistory } from "react-router-dom";
const Login = (props) => {
    const {setToken,token} = props;
    const [loading, setLoading] = useState(false);
    const [disableAll, setDisableAll] = useState(false);
    const navigate =useNavigate();
    useEffect(()=>{
        if (sessionStorage.getItem("repliteweb")!=undefined&&sessionStorage.getItem("repliteweb") !=""){
            // init the token for system
            // login
            navigate('/main');
        }
    },[]);
    const onFinish = async (value) => {
        const { username, password } = value; // Obtain the value of the form input
        try {
            setLoading(true);
            setDisableAll(true); // 禁用其他链接和按钮
            const response = await axios.get("/user/login",{params:{
                "username":username,
                "password":password,
            }});
            // if (response.status!=200){
            //        message.error("系统出错啦");
            //        return;
            // }
            const {state, message: resMessage} = response.data;

            //先检查所有错误并处理
            if (!state ) {
                //登录失败
                let Message = resMessage;
                if (Message == undefined || message == "") {
                    Message = "系统错误";
                }
                message.error(Message);
                return
            }
            // 登录成功，保存JWT Token到浏览器
            const jwt = String(response.headers.get('Authorization')).replace("Bearer ","");
            sessionStorage.setItem('repliteweb', jwt);
            setToken(jwt);
            axios.defaults.headers['Authorization']= 'Bearer ' + sessionStorage.getItem('repliteweb');
            message.success(resMessage).then(()=>navigate("/main"));

        } finally {
            setLoading(false);
            setDisableAll(false); // 解除禁用其他链接和按钮
        }
    };

    return (
        <div style={{paddingRight:"32px",paddingLeft:"32px",maxWidth:"100%",maxHeight:"100%"}} >
            <div style={{width:"100%",height:"20%"}}>
                <div style={{display:"flex",padding:"32px"}}>
                    <img style={{ display:"block",cursor:"pointer"}} width="72px" height="72px" src={logo} alt="logo" />
                </div>
            </div>
            <div style={{width:"100%",height:"70%",display:"flex",justifyContent:"space-between",alignItems:"center",marginTop:"64px"}}>
                <div style={{width:"50%",height:"100%",paddingTop:"12px",paddingBottom:"12px",paddingLeft:"32px",paddingRight:"32px",display:"flex",flexDirection:"column",alignContent:"space-between"}}>
                    <div style={{fontFamily:" Montserrat, sans-serif",fontSize:"72px"}}>Welcome to</div>
                    <div style={{display:"flex",justifyContent:"end"}}>
                        <div style={{fontFamily:"Pacifico cursive",fontSize:"64px"}}>
                            Replite Utils
                        </div>
                    </div>
                </div>
                <Card title="登录" hoverable style={{border:"2px",boxShadow: '0 0 5px rgba(0, 0, 0, 0.3)',padding:"16px",width:"25%",marginRight:"108px",height:"100%"}}>
                    <Form
                        name="normal_login"
                        className="login-form"
                        initialValues={{ remember: true }}
                        onFinish={onFinish}
                    >
                        <Form.Item
                            name="username"
                            rules={[{ required: true, message: 'Please input your Username!' }]}
                        >
                            <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Username" />
                        </Form.Item>
                        <Form.Item
                            name="password"
                            rules={[{ required: true, message: 'Please input your Password!' }]}
                        >
                            <Input
                                prefix={<LockOutlined className="site-form-item-icon" />}
                                type="password"
                                placeholder="Password"
                            />
                        </Form.Item>
                        <Form.Item>
                            <Form.Item name="remember" valuePropName="checked" noStyle>
                                <Checkbox>请记住</Checkbox>
                            </Form.Item>

                            <a className="login-form-forgot" href="/forget" disabled={disableAll}>
                                忘记密码
                            </a>
                        </Form.Item>

                        <Form.Item>
                            <Button type="primary" htmlType="submit" className="login-form-button"loading={loading} disabled={loading || disableAll}>
                                登录
                            </Button>
                            Or <a href="/register"disabled={disableAll}>去注册!</a>
                        </Form.Item>
                    </Form>
                </Card>
            </div>
            <div style={{height:"10%"}}></div>
        </div>

    )
}

export default Login