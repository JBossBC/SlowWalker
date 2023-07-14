

import {Form,Input,Checkbox,Button} from "antd"

import { LockOutlined, UserOutlined } from '@ant-design/icons';

import logo from '../public/logo.png'

import React from "react";
const Login = ()=>{
  function onFinish(){

  }
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

        <a className="login-form-forgot" href="/forget">
          忘记密码
        </a>
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" className="login-form-button">
          登录
        </Button>
        Or <a href="/register">去注册!</a>
      </Form.Item>
    </Form>
    </div>
        </div>
        <div style={{height:"10%"}}></div> 
        </div>
   )
}

export default Login