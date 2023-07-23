import React from 'react';
import { LaptopOutlined, NotificationOutlined, UserOutlined } from '@ant-design/icons';
import { Breadcrumb, Layout, Menu, theme } from 'antd';

const { Header, Content, Sider } = Layout;

const items1 = ['1', '2', '3'].map((key) => ({
    key,
    label: `nav ${key}`,
}));

const items2 = [
    {
        key: 'sub1',
        icon: <UserOutlined />,
        label: '功能',
        children: new Array(4).fill(null).map((_, j) => {
            const subKey = j + 1;
            return {
                key: `sub1/option${subKey}`,
                label: `选项${subKey}`,
            };
        }),
    },
    {
        key: 'sub2',
        icon: <LaptopOutlined />,
        label: '系统',
        children: new Array(4).fill(null).map((_, j) => {
            const subKey = j + 1;
            return {
                key: `sub2/option${subKey}`,
                label: `选项${subKey}`,
            };
        }),
    },
];

const Main = () => {
    const {
        token: { colorBgContainer },
    } = theme.useToken();

    const [selectedMenuKeys, setSelectedMenuKeys] = React.useState([]);
    const [breadcrumbItems, setBreadcrumbItems] = React.useState([]);

    const handleClick = ({ key }) => {
        const [prefix, suffix] = key.split('/');
        let updatedBreadcrumbItems = [];

        if (prefix === 'sub1') {
            updatedBreadcrumbItems = [
                {
                    key: `sub1/${suffix}`,
                    label: `功能/${suffix.replace('option', '选项')}`,
                },
            ];
        } else if (prefix === 'sub2') {
            updatedBreadcrumbItems = [
                {
                    key: `sub2/${suffix}`,
                    label: `系统/${suffix.replace('option', '选项')}`,
                },
            ];
        }

        setSelectedMenuKeys([key]);
        setBreadcrumbItems(updatedBreadcrumbItems);
    };

    return (
        <Layout>
            <Header style={{ display: 'flex', alignItems: 'center' }}>
                <div className="demo-logo"></div>
                <Menu theme="dark" mode="horizontal" defaultSelectedKeys={['2']} items={items1} />
            </Header>
            <Layout>
                <Sider width={200} style={{ background: colorBgContainer }}>
                    <Menu
                        mode="inline"
                        selectedKeys={selectedMenuKeys}
                        defaultOpenKeys={['sub1']}
                        style={{ height: '100%', borderRight: 0 }}
                        onSelect={handleClick}
                        items={items2}
                    />
                </Sider>
                <Layout style={{ padding: '0 24px 24px' }}>
                    <Breadcrumb style={{ margin: '16px 0' }}>
                        {breadcrumbItems.length > 0 ? (
                            breadcrumbItems.map((item) => (
                                <Breadcrumb.Item key={item.key}>{item.label}</Breadcrumb.Item>
                            ))
                        ) : (
                            <Breadcrumb.Item>首页</Breadcrumb.Item>
                        )}
                    </Breadcrumb>
                    <Content style={{ padding: 24, margin: 0, minHeight: 280, background: colorBgContainer }}>
                        {selectedMenuKeys.length > 0 ? (
                            <div>您选择的是：{selectedMenuKeys[0]}</div>
                        ) : (
                            <div>欢迎访问首页</div>
                        )}
                    </Content>
                </Layout>
            </Layout>
        </Layout>
    );
};

export default Main;
