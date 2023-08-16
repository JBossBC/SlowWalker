import React, { useEffect } from "react"
import { useState,useContext } from 'react';
import { useNavigate } from "react-router-dom";
import axios from '../utils/axios';
import {
    Button,
    Checkbox,
    Col,
    Form,
    Input,
    Row,
    Select,
    message
  } from 'antd';
  const formItemLayout = {
    labelCol: {
      xs: {
        span: 24,
      },
      sm: {
        span: 8,
      },
    },
    wrapperCol: {
      xs: {
        span: 24,
      },
      sm: {
        span: 16,
      },
    },
  };
  const tailFormItemLayout = {
    wrapperCol: {
      xs: {
        span: 24,
        offset: 0,
      },
      sm: {
        span: 16,
        offset: 8,
      },
    },
  };
  const { Option } = Select;
const Register=()=>{
    const [form] = Form.useForm();
    // const backendURL = useContext(Backend);
    const navigate = useNavigate();
    //allow to send code
    const [sendCode,setSendCode]=useState(false);
    const [timer,setTimer]=useState(60);
    //const [autoCompleteResult, setAutoCompleteResult] = useState([]);

    const prefixSelector = (
        <Form.Item name="prefix" noStyle>
          <Select
            style={{
              width: 70,
            }}
          >
            <Option value="86">+86</Option>
          </Select>
        </Form.Item>
      );
      const onFinish = async(values) => {
        try {
          let form=new FormData()
          form.append("username",values.username);
          form.append("password",values.password);
          form.append("phone",values.phone);
          form.append("code",values.captcha);
          const response = await axios.post('/user/register', form,{
            headers:{
              'Content-Type':'multipart/form-data'
            }
          });   //response就是一个promise对象，里面的data就是后端返回的值
          const data = response.data;
            if (response.status!='200' || data.state!=true){
             let msg = data.message;
             if (msg == undefined ||msg == ""){
              msg = "系统正在开小差";
             }
             message.error(msg);
             return
          }
          message.success("注册成功").then(()=>navigate('/login'));
            } catch (error) {
              console.log(error);
                message.error("系统出错");  
            }
      };
      useEffect(()=>{
         if(sendCode){
          setTimer(60);
          const task = setInterval(() => {
            if (timer>0){
            setTimer((prevCount) => prevCount - 1);
            }
          }, 1000);
          setTimeout(()=>{clearInterval(task);setSendCode(false)},60000);
         }
      },[sendCode])
      const sendVerificationCode = async (event) => {
        if (sendCode){
          return
        }
        try {
          const response = await axios.get('/phone/send', {
          params: {
              phoneNumber: form.getFieldValue('phone')  //获取phone这个字段的值，携带参数的axios的get请求就等于http://localhost:8080/phone/send?phone=1234567890
          }
          });
          const data = response.data;
          if (response.status!='200' || data.state!=true){
             let msg = data.message;
             if (msg == undefined ||msg == ""){
              msg = "系统正在开小差";
             }
             message.error(msg);
             return
          }
          setSendCode(true);
        }catch(e){
               console.log(e)
               message.error("系统出错");
          }
        
      };


    return(

        <div style={{height:"100%",width:"100%",display:"flex",justifyContent:"center",alignItems:"center"}}>
            <div style={{maxWidth:"600px",height:"80%",textAlign:"center"}}>
            {/* <div style={{margin:"12px",marginBottom:"32px",maxWidth:"600px",fontFamily:" Montserrat, sans-serif",fontSize:"25px"}}>注册</div> */}
            <Form

      {...formItemLayout}
      form={form}
      name="register"
      onFinish={onFinish} //这里才是核心关键所在
      initialValues={{
        prefix: '86',
      }}
      style={{
        maxWidth: 600,
      }}
      scrollToFirstError
    >
      <div style={{margin:"12px",marginBottom:"32px",maxWidth:"600px",fontFamily:" Montserrat, sans-serif",fontSize:"25px"}}>注册</div>
  
      <Form.Item
        name="username"
        label="账号"
        rules={[
          {
            required: true,
            message: '请输入你的账号!',
          },
        ]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        name="password"
        label="密码"
        rules={[
          {
            required: true,
            message: '请输入你的密码!',
          },
        ]}
        tooltip="密码必须包含特殊字符,大写字母和数字,长度必须大于12位"
        hasFeedback
      >
        <Input.Password />
      </Form.Item>

      <Form.Item
        name="confirm"
        label="再次输入密码"
        dependencies={['password']}
        hasFeedback
        tooltip="密码必须包含特殊字符,大写字母和数字,长度必须大于12位"
        rules={[
          {
            required: true,
            message: 'Please confirm your password!',
          },
          ({ getFieldValue }) => ({
            validator(_, value) {
              if (!value || getFieldValue('password') === value) {
                return Promise.resolve();
              }
              return Promise.reject(new Error('The new password that you entered do not match!'));
            },
          }),
        ]}
      >
        <Input.Password />
      </Form.Item>

      <Form.Item
        name="phone"
        label="手机号码"
        rules={[
          {
            required: true,
            message: '请输入你的手机号码!',
          },
        ]}
        tooltip="用来找回密码"
      >
        <Input
          addonBefore={prefixSelector}
          style={{
            width: '100%',
          }}
        />
      </Form.Item>

      <Form.Item label="验证码" extra="We must make sure that your are a human.">
        <Row gutter={8}>
          <Col span={12}>
            <Form.Item
              name="captcha"
              noStyle
              rules={[
                {
                  required: true,
                  message: '验证码不能为空!',
                },
              ]}
            >
              <Input />
            </Form.Item>
          </Col>
          <Col span={12}>
          <Button style={{width:"100px"}} disabled={sendCode} onClick={sendVerificationCode}>{sendCode?timer+"秒":"获取验证码"}</Button>   
          </Col>
        </Row>
      </Form.Item>

      <Form.Item
        name="agreement"
        valuePropName="checked"
        rules={[
          {
            validator: (_, value) =>
              value ? Promise.resolve() : Promise.reject(new Error('Should accept agreement')),
          },
        ]}
        {...tailFormItemLayout}
      >
        <Checkbox>
         我已经阅读 <a href="">协议</a>
        </Checkbox>
      </Form.Item>
      <Form.Item {...tailFormItemLayout}>
        <Button type="primary" htmlType="submit">
          注册
        </Button>
      </Form.Item>
    </Form>
        </div>
    </div>
    )
}


export default Register