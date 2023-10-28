//TODO ignore this file change for language
import React, { useState } from 'react';
import { Button, Input, Upload, message, Modal, Form, Select, InputNumber } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import { FormattedMessage } from "react-intl";
const FileCut = () => {
    const [fileList, setFileList] = useState([]);
    const [outputFolderPath, setOutputFolderPath] = useState('');
    const [modalVisible, setModalVisible] = useState(false);
    const [form] = Form.useForm();

    const handleFileChange = (info) => {
        setFileList(info.fileList);
    };

    const handleOutputFolderChange = (e) => {
        setOutputFolderPath(e.target.value);
    };

    const showModal = () => {
        setModalVisible(true);
    };

    const handleModalCancel = () => {
        setModalVisible(false);
    };

    const handleSplit = () => {
        form.validateFields().then((values) => {
            console.log(values);
            console.log("文件拆分");
        });
    };

    const handleMerge = () => {
        console.log("合并文件");
    };

    const beforeUpload = (file) => {
        const isExcel =
            file.type === 'application/vnd.ms-excel' ||
            file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet';
        if (!isExcel) {
            message.error('只能上传表格文件');
        }
        return isExcel;
    };

    const onFinish = (values) => {
        console.log(values);
        setModalVisible(false);
    };

    const onFinishFailed = (errorInfo) => {
        console.log(errorInfo);
    };

    return (
        <div>
            <h2><FormattedMessage id="文件拆分与合并"/></h2>
            <div>
                <Upload onChange={handleFileChange} beforeUpload={beforeUpload} fileList={fileList}>
                    <Button>
                        <UploadOutlined /> <FormattedMessage id="选择文件"/>
                    </Button>
                </Upload>
            </div>
            <div>
                <h1><FormattedMessage id="选择文件保存地址"/></h1>
                <Input placeholder="输出文件夹路径" value={outputFolderPath} onChange={handleOutputFolderChange} />
            </div>
            <div>
                <Button onClick={showModal}><FormattedMessage id="文件拆分"/></Button>
            </div>
            <div>
                <Button onClick={handleMerge}><FormattedMessage id="文件合并"/></Button>
            </div>
            <Modal
                visible={modalVisible}
                onCancel={handleModalCancel}
                footer={[
                    <Button key="cancel" onClick={handleModalCancel}>
                        <FormattedMessage id="取消"/>
                    </Button>,
                    <Button key="split" type="primary" onClick={form.submit}>
                        <FormattedMessage id="拆分"/>
                    </Button>,
                ]}
            >
                <Form
                    form={form}
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    labelCol={{ span: 6 }}
                    wrapperCol={{ span: 16 }}
                >
                    <Form.Item
                        name="splitType"
                        label="拆分方式"
                        rules={[{ required: true, message: '请选择拆分方式' }]}
                    >
                        <Select>
                            <Select.Option value="count">按文件个数</Select.Option>
                            <Select.Option value="size">按文件大小</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        noStyle
                        shouldUpdate={(prevValues, currentValues) =>
                            prevValues.splitType !== currentValues.splitType
                        }
                    >
                        {({ getFieldValue }) => {
                            return getFieldValue('splitType') === 'count' ? (
                                <Form.Item
                                    name="fileCount"
                                    label="拆分为几个文件"
                                    rules={[{ required: true, message: '请输入拆分文件个数' }]}
                                >
                                    <InputNumber min={1} />
                                </Form.Item>
                            ) : (
                                <Form.Item
                                    name="fileSize"
                                    label="拆分文件后每个文件的大小"
                                    rules={[{ required: true, message: '请输入拆分文件大小' }]}
                                >
                                    <InputNumber min={1} />
                                </Form.Item>
                            );
                        }}
                    </Form.Item>
                    <Form.Item
                        name="fileType"
                        label="存储文件类型"
                        rules={[{ required: true, message: '请选择存储文件类型' }]}
                    >
                        <Select>
                            <Select.Option value="xlsx">Excel ( .xlsx )</Select.Option>
                            <Select.Option value="csv">CSV ( .csv )</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        name="encoding"
                        label="文件编码格式"
                        rules={[{ required: true, message: '请选择文件编码格式' }]}
                    >
                        <Select>
                            <Select.Option value="utf-8">UTF-8</Select.Option>
                            <Select.Option value="gbk">GBK</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        name="otherParameter"
                        label="其他参数"
                    >
                        <Input.TextArea />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default FileCut;
