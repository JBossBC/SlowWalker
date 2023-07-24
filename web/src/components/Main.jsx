import React, { useContext,useEffect } from 'react';
import { DesktopOutlined,ExclamationCircleOutlined} from '@ant-design/icons';
import { Breadcrumb, Layout, Menu, theme,message } from 'antd';
import {Backend} from "../App"
import axios from 'axios'


// const items1 = ['1', '2', '3'].map((key) => ({
//   key,
//   label: `nav ${key}`,
// }));
// const items2 = [UserOutlined, LaptopOutlined, NotificationOutlined].map((icon, index) => {
//   const key = String(index + 1);
//   return {
//     key: `sub${key}`,
//     icon: React.createElement(icon),
//     label: `subnav ${key}`,
//     children: new Array(4).fill(null).map((_, j) => {
//       const subKey = index * 4 + j + 1;
//       return {
//         key: subKey,
//         label: `option${subKey}`,
//       };
//     }),
//   };
// });

const { Header, Content, Sider } = Layout;
const itemsForIcon = {"function":DesktopOutlined,"system":ExclamationCircleOutlined}


//TODO 对每一个功能添加对应的类别
// const menuItems = [DesktopOutlined,ExclamationCircleOutlined].map((icon,index)=>{
//    const key=String(index+1);
//    return {
       
//    }
// })
const Main = () => {
  const BackendURL = useContext(Backend)
  const {
    token: { colorBgContainer },
  } = theme.useToken();
  useEffect(()=>{
    async function initMain(){
      // get the authorization data
      await axios.get(BackendURL+"/rule/query").then((response)=>{
        let data =response.data;
        if (response.status != '200' ||data.state!=true){
          message.error("加载首页出错");
          // 返回首页

          return
        }
      })
    }

    initMain(); //1

  },[])
  return (
    <Layout>
      <Header
        style={{
          display: 'flex',
          alignItems: 'center',
        }}
      >
      </Header>
      <Layout>
        <Sider
          width={200}
          style={{
            background: colorBgContainer,
          }}
        >
          <Menu
            mode="inline"
            defaultSelectedKeys={['1']}
            defaultOpenKeys={['sub1']}
            style={{
              height: '100%',
              borderRight: 0,
            }}
          />
        </Sider>
        <Layout
          style={{
            padding: '0 24px 24px',
          }}
        >
          
          <Breadcrumb
            style={{
              margin: '16px 0',
            }}
          >
            <Breadcrumb.Item>首页</Breadcrumb.Item>
            <Breadcrumb.Item>List</Breadcrumb.Item>
            <Breadcrumb.Item>App</Breadcrumb.Item>
          </Breadcrumb>
          <Content
            style={{
              padding: 24,
              margin: 0,
              minHeight: 280,
              background: colorBgContainer,
            }}
          >
            Content
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};
export default Main;