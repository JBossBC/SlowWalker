import React, { useState, useEffect } from 'react';
import { Layout, Menu, Breadcrumb, theme, Form, Button } from 'antd';
import { LaptopOutlined, UserOutlined, PlusOutlined, MonitorOutlined } from '@ant-design/icons';
import axios from "../utils/axios";
import rbac from "../utils/rbac.json";
import IPQuery from './IPQuery';
import FileMergeCut from "./FileMergeCut";
import AddFunctionModal from './AddFunctionModal';
import UserManage from "./users";
import Log from "./Log";
const { Header, Content, Sider } = Layout;

const Main = () => {
    // const intl = useIntl();
    const { token: { colorBgContainer } } = theme.useToken();
    const [selectedMenuKey, setSelectedMenuKey] = useState("");
    const [breadcrumbItem, setBreadcrumbItem] = useState("");
    const [menuItems, setMenuItems] = useState([]);
    const [Showcomponent, setShowComponent] = useState(null);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [form] = Form.useForm();
    const [params, setParams] = useState([]);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            const response = await axios.get("/rule/query");
            const { state, data } = response.data;
            console.log(data.role);
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
                            // component: getComponentByLabel(item.component),
                        });
                    }
                });
                setMenuItems(items);
            }
        } catch (error) {
            console.error("Error fetching data:", error);
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
            updatedBreadcrumbItem = `${{ id: "首页" }}/${{ id: belongs }}/${{ id: suffix }}`;
        }

        setSelectedMenuKey(key);
        setBreadcrumbItem(updatedBreadcrumbItem);

        const menuItem = menuItems.find((item) => item.key === prefix+"/"+belongs);
        if (menuItem && menuItem.component) {
            getComponentByLabel(suffix).then((comp)=>setShowComponent(comp));
        } else {
            getComponentByLabel(suffix).then((comp)=>setShowComponent(comp));
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

    const handleAddParam = () => {
        const newParams = [...params];
        newParams.push({ name: '', type: '' });
        setParams(newParams);
    };

    const handleDeleteParam = (index) => {
        const newParams = [...params];
        newParams.splice(index, 1);
        setParams(newParams);
    };

    const handleParamChange = (index, field, value) => {
        const newParams = [...params];
        newParams[index][field] = value;
        setParams(newParams);
    };

    return (
        <Layout>
            <Header style={{ display: "flex", alignItems: "center" }}>
                <div className="demo-logo"></div>
            </Header>
            <Layout>
                {menuItems.length > 0 && (
                    <Sider width={200} style={{ background: colorBgContainer }}>
                        <Menu
                            mode="inline"
                            selectedKeys={[selectedMenuKey]}
                            defaultOpenKeys={menuItems.map((item) => item.key)}
                            style={{ height: "100%", borderRight: 0 }}
                            onClick={handleClick}
                        >
                            {menuItems.map((item) => (
                                <Menu.SubMenu key={item.key} icon={item.icon} title={item.label} items={item.children}>
                                    {item.children.map((subItem) => (
                                        <Menu.Item key={subItem.key}>{subItem.label}</Menu.Item>
                                    ))}
                                </Menu.SubMenu>
                            ))}
                            <div key="addFunction" style={{ marginTop: "10px", display: "flex", justifyContent: "center" }}>
                                <Button icon={<PlusOutlined />} onClick={showModal}>
                                    添加
                                </Button>
                            </div>
                        </Menu>
                        <AddFunctionModal
                            isModalVisible={isModalVisible}
                            handleOk={handleOk}
                            handleCancel={handleCancel}
                            params={params}
                            setParams={setParams}
                            handleAddParam={handleAddParam}
                            handleDeleteParam={handleDeleteParam}
                            handleParamChange={handleParamChange}
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
                    <Content style={{ padding: 24, margin: 0, minHeight: 280, background: colorBgContainer }}>
                        {selectedMenuKey && !Showcomponent && (
                            <div>
                               您选择的是:
                                {selectedMenuKey}
                            </div>
                        )}
                        {!selectedMenuKey && <div>欢迎访问首页</div>}
                        {Showcomponent}
                    </Content>
                </Layout>
            </Layout>
        </Layout>
    );
};

export default Main;