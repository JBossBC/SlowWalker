import React, {useState} from 'react';
import {Button, Input, Table, Upload} from 'antd';
import {DownloadOutlined, ExportOutlined, InboxOutlined} from '@ant-design/icons';

const {Dragger} = Upload;

const IPQuery = () => {
    const [ipAddresses, setIpAddresses] = useState('');
    const [tableData, setTableData] = useState([]);
    const [showExportButton, setShowExportButton] = useState(false);

    const handleIpAddressesChange = (e) => {
        setIpAddresses(e.target.value);
    };

    const handleIpQuery = () => {
        const data = [{ip: '192.168.0.1', location: '北京'}, {ip: '192.168.0.2', location: '上海'}, {
            ip: '192.168.0.2',
            location: '天天'
        }, {ip: '192.168.0.2', location: '得到'}, {ip: '192.168.0.2', location: '阿发'}, {
            ip: '192.168.0.2',
            location: '嘎斯'
        }, {ip: '192.168.0.2', location: '结果'}, {ip: '192.168.0.2', location: '阿斯弗'}, {
            ip: '192.168.0.2',
            location: '多个'
        },];
        setTableData(data);
        setShowExportButton(true);
    };

    const handleFileUpload = (file) => {

    };

    const downloadTemplate = () => {

    };

    const exportToCSV = () => {

    };

    const columns = [{
        title: "IP地址", dataIndex: 'ip', key: 'ip',
    }, {
        title: "位置", dataIndex: 'location', key: 'location',
    },];

    return (<div>
            <div> IP查询</div>
            <div>
                <Dragger
                    showUploadList={false}
                    beforeUpload={(file) => {
                        handleFileUpload(file);
                        return false;
                    }}
                >
                    <p className="ant-upload-drag-icon">
                        <InboxOutlined/>
                    </p>
                    <p className="ant-upload-text">点击或拖拽上传文件</p>
                </Dragger>
                <Button type="primary" style={{marginRight: 10}} onClick={downloadTemplate}>
                    下载格式范本
                    <DownloadOutlined/>
                </Button>
                <Button type="primary" style={{marginRight: 10}} onClick={handleIpQuery}>
                    查询
                </Button>
                {showExportButton && (<Button type="primary" onClick={exportToCSV}>
                        导出查询结果
                        <ExportOutlined/>
                    </Button>)}
            </div>
            <Input.TextArea
                rows={4}
                cols={50}
                placeholder={"请输入IP地址每个IP地址一行最多一次5000条"}
                value={ipAddresses}
                onChange={handleIpAddressesChange}
            />
            <Table columns={columns} dataSource={tableData}/>
        </div>);
};

export default IPQuery;