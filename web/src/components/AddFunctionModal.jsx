import React, { useState } from 'react';
import { Modal, Form, Input, Select, Button, Upload, message } from 'antd';
import { PlusOutlined, DeleteOutlined, InboxOutlined, ArrowLeftOutlined } from '@ant-design/icons';
import {Option} from "antd/es/mentions";
//import { FormattedMessage, useIntl } from "react-intl";

const { Dragger } = Upload;

const AddFunctionModal = ({isModalVisible,handleOk,handleCancel}) => {
    const [params, setParams] = useState([]);
    const [showNextModal, setShowNextModal] = useState(false);
    const [description, setDescription] = useState('');
    const [author, setAuthor] = useState('');
    const [explanation, setExplanation] = useState('');

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

    const handlePush = () => {
        //handleOk();调用父组件传入的handleOk函数or本组件写
        setShowNextModal(false);
    };

    const handleValidate = () => {
        form.validateFields((errors, values) => {
            if (errors) {
                message.warning('请填写必填项！');
            } else {
                setShowNextModal(true);
            }
        });
    };

    const [form] = Form.useForm();

    return (
        <React.Fragment>
            <Modal title= "添加自定义功能" visible={isModalVisible && !showNextModal} onOk={() => setShowNextModal(true)} onCancel={handleCancel}>
                <Form form={form}>
                    <Form.Item
                        label="功能名称"
                        name="functionName"
                        rules={[{ required: true, message: '请输入功能名称!' }]}
                    >
                        <Input.TextArea rows={1} placeholder="请填写功能名称" />
                    </Form.Item>

                    <Form.Item
                        label="代码语言"
                        name="language"
                        rules={[{ required: true, message: '请输入代码语言!' }]}
                    >
                        <Select placeholder="请选择代码语言">
                            <Option value="Python">Python</Option>
                            <Option value="Golang">Golang</Option>
                            <Option value="Java">Java</Option>
                            <Option value="JavaScript">JavaScript</Option>
                            <Option value="Rust">Rust</Option>
                            <Option value="C">C</Option>
                            <Option value="C++">C++</Option>
                            <Option value="PHP">PHP</Option>
                        </Select>
                    </Form.Item>

                    <Form.Item
                        label="上传功能文件"
                        name="content"
                        rules={[{ required: true, message: '请上传功能文件!' }]}
                    >
                        <Dragger name='file' multiple={false} beforeUpload={() => false}>
                            <p className="ant-upload-drag-icon">
                                <InboxOutlined />
                            </p>
                            <p className="ant-upload-text">点击或拖拽文件到这个区域以上传</p>
                        </Dragger>
                    </Form.Item>

                    {params.map((param, index) => (
                        <div key={`param-${index}`} style={{ display: 'flex', marginBottom: 16 }}>
                            <Form.Item style={{ marginRight: 8 }} label={`名称${index + 1}`} labelCol={{ flex: "0 0 auto" }} wrapperCol={{ flex: "1" }}>
                                <Input
                                    value={param.name}
                                    onChange={(e) => handleParamChange(index, 'name', e.target.value)}
                                />
                            </Form.Item>
                            <Form.Item style={{ marginRight: 8, width: "50%" }} label={`类型${index + 1}`}>
                                <Select
                                    value={param.type}
                                    onChange={(value) => handleParamChange(index, 'type', value)}
                                    style={{ width: "100%" }}
                                >
                                    <Option value="string">字符串</Option>
                                    <Option value="number">数字</Option>
                                    <Option value="boolean">布尔值</Option>
                                    <Option value="file">文件</Option>
                                </Select>
                            </Form.Item>
                            <Button icon={<DeleteOutlined />} onClick={() => handleDeleteParam(index)} />
                        </div>
                    ))}

                    <Button type="dashed" style={{ marginTop: 16 }} onClick={handleAddParam}>
                        <PlusOutlined /> 添加参数
                    </Button>
                </Form>
            </Modal>

            <Modal  title="填写信息" visible={showNextModal} onOK={handlePush} onCancel={() => setShowNextModal(false)}>
                <Button type="primary" icon={<ArrowLeftOutlined />} onClick={() => setShowNextModal(false)}>返回</Button>
                <Form.Item
                    label="作者"
                    name="author"
                    rules={[{ required: true, message: '请输入作者!' }]}
                >
                    <Input placeholder="请填写作者" onChange={e => setAuthor(e.target.value)} />
                </Form.Item>
                <Form.Item
                    label="功能描述"
                    name="description"
                    rules={[{ required: true, message: '请输入功能描述!' }]}
                >
                    <Input.TextArea rows={3} placeholder="请填写功能描述" onChange={e => setDescription(e.target.value)} />
                </Form.Item>
                <Form.Item
                    label="使用说明"
                    name="explanation"
                    rules={[{ required: true, message: '请输入使用说明!' }]}
                >
                    <Input.TextArea rows={4} placeholder="请填写使用说明" onChange={e => setExplanation(e.target.value)} />
                </Form.Item>
            </Modal>
        </React.Fragment>
    );
};

export default AddFunctionModal;
