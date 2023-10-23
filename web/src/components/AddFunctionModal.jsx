import React from 'react';
import { Modal, Form, Input, Select, Button } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import {Option} from "antd/es/mentions";

const AddFunctionModal = ({ isModalVisible, handleOk, handleCancel, params, setParams, handleAddParam, handleDeleteParam, handleParamChange }) => {
    return (
        <Modal title="添加自定义代码" visible={isModalVisible}  onOk={handleOk} onCancel={handleCancel}>
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
                label="代码内容"
                name="content"
                rules={[{ required: true, message: '请输入代码内容!' }]}
            >
                <Input.TextArea rows={4} />
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
        </Modal>
    );
};

export default AddFunctionModal;
