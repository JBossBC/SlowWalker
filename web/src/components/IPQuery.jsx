import React, { useState } from 'react';
import { Button, Input, Table, Upload } from 'antd';
import { DownloadOutlined, ExportOutlined, InboxOutlined } from '@ant-design/icons';
import { FormattedMessage, useIntl } from "react-intl";
const { Dragger } = Upload;

const IPQuery = () => {
    const [ipAddresses, setIpAddresses] = useState('');
    const [tableData, setTableData] = useState([]);
    const [showExportButton, setShowExportButton] = useState(false);
    const intl = useIntl();

    const handleIpAddressesChange = (e) => {
        setIpAddresses(e.target.value);
    };

    const handleIpQuery = () => {
        const data = [
            { ip: '192.168.0.1', location: '北京' },
            { ip: '192.168.0.2', location: '上海' },
            { ip: '192.168.0.2', location: '天天' },
            { ip: '192.168.0.2', location: '得到' },
            { ip: '192.168.0.2', location: '阿发' },
            { ip: '192.168.0.2', location: '嘎斯' },
            { ip: '192.168.0.2', location: '结果' },
            { ip: '192.168.0.2', location: '阿斯弗' },
            { ip: '192.168.0.2', location: '多个' },
        ];
        setTableData(data);
        setShowExportButton(true);
    };

    const handleFileUpload = (file) => {
        /* 处理上传文件逻辑
         可以在这里读取文件内容，并将内容设置为ipAddresses*/
        /*const reader = new FileReader();
              reader.onload = (event) => {
               setIpAddresses(event.target.result);
              };
              reader.readAsText(file);*/
    };

    const downloadTemplate = () => {
        /* 下载文件范本逻辑
         可以使用window.open或axios库进行文件下载*/
        /*  axios.get('/template-path', { responseType: 'blob' }).then((response) => {
                   const url = window.URL.createObjectURL(new Blob([response.data]));
                   const link = document.createElement('a');
                   link.href = url;
                   link.setAttribute('download', 'template.csv');
                   document.body.appendChild(link);
                   link.click();
                 });*/
    };

    const exportToCSV = () => {
        /* 导出表格为CSV文件逻辑
         可以使用json2csv库将tableData转换为CSV格式，并进行文件下载*/
        /*const csv = convertToCSV(tableData);
              const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
              const url = URL.createObjectURL(blob);
              const link = document.createElement('a');
              link.href = url;
              link.setAttribute('download', 'ip_addresses.csv');
               link.click();*/
    };

    const columns = [
        {
            title: intl.formatMessage({id: "IP地址"}),
            dataIndex: 'ip',
            key: 'ip',
        },
        {
            title: intl.formatMessage({id: "位置"}),
            dataIndex: 'location',
            key: 'location',
        },
    ];

    return (
        <div>
            <div> <FormattedMessage id="IP查询" /></div>
            <div>
                <Dragger
                    showUploadList={false}
                    beforeUpload={(file) => {
                        handleFileUpload(file);
                        return false;
                    }}
                >
                    <p className="ant-upload-drag-icon">
                        <InboxOutlined />
                    </p>
                    <p className="ant-upload-text"><FormattedMessage id="点击或拖拽上传文件" /></p>
                </Dragger>
                <Button type="primary" style={{ marginRight: 10 }} onClick={downloadTemplate}>
                    <FormattedMessage id="下载格式范本" />
                    <DownloadOutlined />
                </Button>
                <Button type="primary" style={{ marginRight: 10 }} onClick={handleIpQuery}>
                    <FormattedMessage id="查询" />
                </Button>
                {showExportButton && ( // 只有有查询结果时才显示导出按钮
                    <Button type="primary" onClick={exportToCSV}>
                        <FormattedMessage id="导出查询结果" />
                        <ExportOutlined />
                    </Button>
                )}
            </div>
            <Input.TextArea
                rows={4}
                cols={50}
                placeholder={intl.formatMessage({id: "请输入IP地址每个IP地址一行最多一次5000条"})}
                value={ipAddresses}
                onChange={handleIpAddressesChange}
            />
            <Table columns={columns} dataSource={tableData} />
        </div>
    );
};

export default IPQuery;