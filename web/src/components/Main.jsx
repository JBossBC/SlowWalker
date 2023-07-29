import React, { useEffect, useState } from 'react';
import { LaptopOutlined, UserOutlined } from '@ant-design/icons';
import { Layout, Menu, Breadcrumb, theme } from 'antd';
import axios from "axios";

const { Header, Content, Sider } = Layout;
const localURL = "http://localhost:8080";
//const testURL = "http://112.124.53.234:8080";d
const Main = () => {
    const { token: { colorBgContainer } } = theme.useToken();

    const [selectedMenuKey, setSelectedMenuKey] = useState('');
    const [breadcrumbItem, setBreadcrumbItem] = useState('');
    const [menuItems, setMenuItems] = useState([]);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            const response = await axios.get(localURL + "/rule/query");
            const { state, data } = response.data;
            console.log(data);
            if (state) {
                const items = [];

                Object.keys(data).forEach((key) => {
                    const subItems = Object.keys(data[key]).map((subKey) => ({
                        key: `sub/${key}/${subKey}`,
                        label: subKey,
                    }));

                    items.push({
                        key: `sub/${key}`,
                        icon: key === 'function' ? <UserOutlined /> : <LaptopOutlined />,
                        label: key === 'function' ? '功能' : '系统',
                        children: subItems,
                    });
                });

                setMenuItems(items);
            }
        } catch (error) {
            console.error("Error fetching data:", error);
        }
    };

    const handleClick = ({ key }) => {
        const [prefix, type, suffix] = key.split('/');
        let updatedBreadcrumbItem = '';

        if (prefix === 'sub') {
            updatedBreadcrumbItem = `首页/${type === 'function' ? '功能' : '系统'}/${suffix.replace('option', '选项')}`;
        }

        setSelectedMenuKey(key);
        setBreadcrumbItem(updatedBreadcrumbItem);
    };

    return (
        <Layout>
            <Header style={{ display: 'flex', alignItems: 'center' }}>
                <div className="demo-logo"></div>
            </Header>
            <Layout>
                {menuItems.length > 0 && (
                    <Sider width={200} style={{ background: colorBgContainer }}>
                        <Menu
                            mode="inline"
                            selectedKeys={[selectedMenuKey]}
                            defaultOpenKeys={menuItems.map((item) => item.key)}
                            style={{ height: '100%', borderRight: 0 }}
                            onClick={handleClick}
                        >
                            {menuItems.map((item) => (
                                <Menu.SubMenu key={item.key} icon={item.icon} title={item.label} items={item.children}>
                                    {item.children.map((subItem) => (
                                        <Menu.Item key={subItem.key}>{subItem.label}</Menu.Item>
                                    ))}
                                </Menu.SubMenu>
                            ))}
                        </Menu>
                    </Sider>
                )}
                <Layout style={{ padding: '0 24px 24px' }}>
                    <Breadcrumb style={{ margin: '16px 0' }}>
                        {breadcrumbItem ? (
                            <Breadcrumb.Item>{breadcrumbItem}</Breadcrumb.Item>
                        ) : (
                            <Breadcrumb.Item>首页</Breadcrumb.Item>
                        )}
                    </Breadcrumb>
                    <Content style={{ padding: 24, margin: 0, minHeight: 280, background: colorBgContainer }}>
                        {selectedMenuKey && <div>您选择的是：{selectedMenuKey}</div>}
                        {!selectedMenuKey && <div>欢迎访问首页</div>}
                    </Content>
                </Layout>
            </Layout>
        </Layout>
    );
};

export default Main;
