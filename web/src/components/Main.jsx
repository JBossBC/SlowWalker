import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router';  //add
import { Layout, Menu, Breadcrumb, Form,theme, Button } from 'antd';
import { LaptopOutlined, UserOutlined, PlusOutlined, MonitorOutlined,LogoutOutlined } from '@ant-design/icons';
import axios from "../utils/axios";
import rbac from "../utils/rbac.json";
import IPQuery from './IPQuery';
import FileMergeCut from "./FileMergeCut";
import AddFunctionModal from './AddFunctionModal';
import UserManage from "./users";
import Log from "./Log";
import Search from './Search';  //add
const { Header, Content, Sider, Footer } = Layout;  //add Footer


const Main = () => {
    const { token: { colorBgContainer } } = theme.useToken();
    const [selectedMenuKey, setSelectedMenuKey] = useState("");
    const [breadcrumbItem, setBreadcrumbItem] = useState("");
    const [menuItems, setMenuItems] = useState([]);
    const [Showcomponent, setShowComponent] = useState(null);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [form] = Form.useForm();
    const [showLog, setShowLog] = useState(false); // add
    const navigate = useNavigate(); //add
    

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            const response = await axios.get("/rule/query");
            const { state, data } = response.data;
            //const funcResponse = await axios.get("func");
            //const { funcState, funcData } = funcResponse.data; // funcResponse
            if (state) {
                const items = [];
                const map1 = rbac.RBAC; //map1=rbac+funcState

                map1.forEach((item) => {
                    if (data.role.includes(item.owner)) {
                        const menuItem = items.find((menu) => menu.label === item.belongs);

                        if (!menuItem) {
                            items.push({
                                key: `首页/${item.belongs}`,
                                icon: getIconByBelongs(item.belongs),
                                label: item.belongs,
                                children: [],
                            });
                        }

                        const subItems = items.find((menu) => menu.label === item.belongs).children;
                        subItems.push({
                            key: `首页/${item.belongs}/${item.component}`,
                            label: item.component,
                            component: getComponentByLabel(item.component),
                        });
                    }
                });
                setMenuItems(items);
            }
        } catch (error) {
            console.error("Error fetching data:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const getIconByBelongs = (belongs) => {
        switch (belongs) {
            case "功能":
                return <UserOutlined />;
            case "系统":
                return <LaptopOutlined />;
            case "管理":
                return <MonitorOutlined />;
            default:
                return null;
        }
    };

    const getComponentByLabel = async (label) => {
        switch (label) {
            case "审计日志":
                return (<Log/>)
            case "人员管理":
                return (<UserManage/>)
            case "IP查询":
                return (<IPQuery/>)
            case "文件切分":
                return (<FileMergeCut/>);
            default:
                const component = await fetchComponent(label);
                return component ? () => component : null;
        }
    };

    const fetchComponent = async (label) => {
        const userFuncResponse = await axios.get("func/component/");
        const { userFuncState, userFuncData } = userFuncResponse.data;
        if (userFuncState) {
            const { component } = await userFuncData.json();
            const Component = eval(component); // 使用eval()动态执行组件代码
            return Component;
        } else {
            return null;
        }
    };


    const handleClick = ({ key }) => {
        const [prefix, belongs, suffix] = key.split("/");
        let updatedBreadcrumbItem = "";

        if (prefix === "首页") {
            updatedBreadcrumbItem = `${prefix}/${belongs}/${suffix}`;
        }
        
        setSelectedMenuKey(key);
        console.log(selectedMenuKey)
        setBreadcrumbItem(updatedBreadcrumbItem);

        const menuItem = menuItems.find((item) => item.key === `${prefix}/${belongs}`);
        if (menuItem && menuItem.label) {
            getComponentByLabel(suffix).then((comp) => setShowComponent(comp));
        }

    };


    const showModal = () => {
        setIsModalVisible(true);
    };

    const handleOk = () => {
        form.submit();
    };

    const handleCancel = () => {
        setIsModalVisible(false);
    };


    const handleLogout = async () => {
        try {
            // 清除本地存储
            sessionStorage.removeItem('repliteweb')
            // 导航到登录页面
            navigate('/login')
          } catch (error) {
            console.error(error)
          }
    }
    const Theme = {
        bodyBg: '#f5f',
        footerBg: '#f5f5f5',
        footerPadding: '24px 50px',
        headerBg: '#00000',
        headerColor: 'rgba(0, 0, 0, 0.88)',
        headerHeight: 64,
        headerPadding: '0 50px',
        lightSiderBg: '#ffffff',
        lightTriggerBg: '#ffffff',
        lightTriggerColor: 'rgba(0, 0, 0, 0.88)',
        siderBg: '#00000',
        triggerBg: '#002140',
        triggerColor: '#fff',
        triggerHeight: 48,
        zeroTriggerHeight: 40,
        zeroTriggerWidth: 40
    }

    return (
        <Layout>
            <Header style={{ display: "flex", alignItems: "center",justifyContent: "space-between", 
            background : Theme.headerBg

            }}>
                <div className="demo-logo"></div>
                <Button
                    icon={<LogoutOutlined />}
                    style={{ marginLeft: "auto" }}
                    onClick={handleLogout}
                > 
                退出登录
                </Button>   
            </Header>
            <Layout>
                {!isLoading && menuItems.length > 0 && (
                    <Sider width={200} style={{ background: colorBgContainer }}>
                        <Menu
                            mode="inline"
                            selectedKeys={[selectedMenuKey]}
                            defaultOpenKeys={[]}
                            style={{ height: "100%", borderRight: 0 }}
                            onClick={handleClick}
                        >
                            {menuItems.map((item) => (
                                <Menu.SubMenu key={item.key} icon={item.icon} title={item.label} items={item.children} defaultOpen={false} >
                                    {item.children.map((subItem) => (
                                        <Menu.Item key={subItem.key}>{subItem.label}</Menu.Item>
                                    ))}
                                </Menu.SubMenu>
                            ))}
                            {menuItems.some((item) => item.label === '功能') && (
                                <div key="addFunction" style={{ marginTop: '10px', display: 'flex', justifyContent: 'center' }}>
                                    <Button icon={<PlusOutlined />} onClick={showModal}>
                                        添加
                                    </Button>
                                </div>
                            )}
                        </Menu>
                        <AddFunctionModal
                            isModalVisible={isModalVisible}
                            handleOk={handleOk}
                            handleCancel={handleCancel}
                        />
                    </Sider>
                )}
                <Layout style={{ padding: "0 24px 24px" }}>
                    <Breadcrumb style={{ margin: "16px 0" }}>
                        {breadcrumbItem ? (
                            <Breadcrumb.Item>{breadcrumbItem}</Breadcrumb.Item>
                        ) : (
                            <Breadcrumb.Item>
                                首页
                            </Breadcrumb.Item>
                        )}
                    </Breadcrumb>
                    <div style={{ minHeight: "calc(100vh - 200px)" }}>
                        <Content style={{ padding: 24, margin: 0, minHeight: 280, background: colorBgContainer }}>
                            {selectedMenuKey && !Showcomponent && (
                                <div>
                                您选择的是:
                                    {selectedMenuKey}
                                </div>
                                
                            )}
                            {selectedMenuKey === "首页/管理/功能管理" && (
                                <Search />
                            )}
                            {!selectedMenuKey &&   
                            <div style={{ display: "flex", flexDirection: "column", justifyContent: "center", 
                            alignItems: "center", height: "100%" }}>
                            <h3 style={{ fontSize: "4rem", fontWeight: "bold", margin: 0 ,background: "#00000" }}>欢迎使用</h3>
                            <p style={{ fontSize: "1rem", margin: "1rem 0" }}>这里是 RepliteWeb Util 工具 首页</p>
                            </div>}
                            {Showcomponent}
                        </Content>
                    </div>
                </Layout>
            </Layout>
            <Footer style={{ textAlign: 'center' }}>RepliteWeb Util ©2023 Created by Upsec </Footer>
        </Layout>
    );
};

export default Main;
