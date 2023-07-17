import React from "react"
import { useState } from 'react';
import axios from 'axios';
import {
    Button,
    Checkbox,
    Col,
    Form,
    Input,
    Row,
    Select,
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
    //const [autoCompleteResult, setAutoCompleteResult] = useState([]);
    const [registerResult, setRegisterResult] = useState(null);
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
          const response = axios.post('http://localhost:8080/user/register', values);   //response就是一个promise对象，里面的data就是后端返回的值
          const data = response.data;
            if (data.success) {
                setRegisterResult("注册成功！"); // 注册成功
            } else {
                setRegisterResult(data.message); // 注册失败，显示后端返回的错误信息
            }
            } catch (error) {
                setRegisterResult("注册失败，请稍后再试！"); 
            }
      };

      const sendVerificationCode = async () => {
        try {
          const response = await axios.post('http://localhost:8080/phone/send');
          const data = response.data; 
          if (data.success) {
            console.log("验证码发送成功");
          } else {
            console.log("验证码发送失败");
          }
          } catch (error) {
            console.log(error);
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
          <Button onClick={sendVerificationCode}>获取验证码</Button>   
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