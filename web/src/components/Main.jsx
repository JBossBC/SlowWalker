import React from 'react';
import { LaptopOutlined, UserOutlined } from '@ant-design/icons';
import { Layout, Menu, Breadcrumb, theme } from 'antd';

const { Header, Content, Sider } = Layout;

const items2 = [
    {
        key: 'sub1',
        icon: <UserOutlined />,
        label: '功能',
        children: [],
    },
    {
        key: 'sub2',
        icon: <LaptopOutlined />,
        label: '系统',
        children: [],
    },
];

const Main = () => {
    const {
        token: { colorBgContainer },
    } = theme.useToken();

    const [selectedMenuKey, setSelectedMenuKey] = React.useState('');
    const [breadcrumbItem, setBreadcrumbItem] = React.useState('');

    const handleClick = ({ key }) => {
        const [prefix, suffix] = key.split('/');
        let updatedBreadcrumbItem = '';

        if (prefix === 'sub1') {
            updatedBreadcrumbItem = `功能/${suffix.replace('option', '选项')}`;
        } else if (prefix === 'sub2') {
            updatedBreadcrumbItem = `系统/${suffix.replace('option', '选项')}`;
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
                <Sider width={200} style={{ background: colorBgContainer }}>
                    <Menu
                        mode="inline"
                        selectedKeys={[selectedMenuKey]}
                        defaultOpenKeys={['sub1']}
                        style={{ height: '100%', borderRight: 0 }}
                        onSelect={handleClick}
                        items={items2}
                    />
                </Sider>
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
