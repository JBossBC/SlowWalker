import { Form, Input, Checkbox, Button, message, Modal } from "antd";
import { LockOutlined, UserOutlined } from "@ant-design/icons";
import logo from "../public/logo.png";
import React, { useState,useContext} from "react";
import axios from "axios";
import {Backend} from "../App";
//import { useHistory } from "react-router-dom";

//const localURL="http://localhost:8080/"
//const testURL="http://112.124.53.234:8080/"
const Login = () => {
    const [loading, setLoading] = useState(false);
    const [disableAll, setDisableAll] = useState(false);
    const backendURL = useContext(Backend);
    console.log(backendURL);
    //const history=useHistory();
    const onFinish = async (value) => {
        const { username, password } = value; // Obtain the value of the form input
        try {
            setLoading(true);
            setDisableAll(true); // 禁用其他链接和按钮
            const response = await axios.get( backendURL+ "user/login?username=" + username + "&password=" + password)
            const {state, message: resMessage} = response.data;
            //先检查所有错误并处理
            if (response.status != "200" || !state) {
                //登录失败
                let Message = resMessage;
                if (Message == undefined || message == "") {
                    Message = "系统错误";
                }
                message.error(Message);
                return
            }
            // 登录成功，保存JWT Token到浏览器
            const jwt = response.headers.get('Authorization');
            console.log("Authorization", jwt);
            //localStorage.setItem("jwtToken", jwt);
            /*使用浏览器的会话存储（session Storage）来存储 JWT Token。
            与 Local Storage 不同，会话存储只在当前会话期间有效，
            当浏览器标签页或窗口关闭时会被清除。同样存在安全风险和 XSS 攻击的风险。*/
            sessionStorage.setItem('jwtToken', jwt);
            message.success(resMessage);
            // TODO: 跳转到首页或其他页面
            //history.push("/main");
            window.location.href = "/main";
            console.log(jwt);
        } finally {
            //TODO:自旋
            setLoading(false);
            setDisableAll(false); // 解除禁用其他链接和按钮
        }
    };

    /* const showConfirm = () => {
         Modal.confirm({
             title: "确认要退出吗？",
             content: "退出后将不再保留登录状态。",
             okText: "确认",
             cancelText: "取消",
             onOk: () => {
                 // 清除浏览器中的JWT Token
                 localStorage.removeItem("jwtToken");

                 // TODO: 跳转到登录页面或其他操作
             },
         });
     };*/
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
                <div style={{border:"2px",boxShadow: '0 0 5px rgba(0, 0, 0, 0.3)',padding:"16px",width:"25%",marginRight:"108px",height:"100%"}}>
                    <div  style={{margin:"12px",marginBottom:"32px",fontFamily:" Montserrat, sans-serif",fontSize:"25px"}}>
                        登录
                    </div>
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
                </div>
            </div>
            <div style={{height:"10%"}}></div>
        </div>

    )
}

export default Login