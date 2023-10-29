import axios from '../utils/axios';
import React,{useState, useEffect} from 'react';
import { Empty, Space, Table, Tag,Button,message, Form, Col, Input, Row, Select, theme } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import moment from 'moment';
import { useNavigate } from 'react-router';
import qs from 'qs';

/** 	PRINT LogLevel = "print"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
	INFO  LogLevel = "info"
	PANIC LogLevel = "panic"
    */
const  levelForColor = {"print":"success","info":"blue","error":"error","warn":"warning","panic":"red"}   //定义了日志等级的颜色
const { Option } = Select;
const columns = [   //定义了日志审计页面中的表格列
  {
    title: '日志等级',
    key: 'level',
    dataIndex: 'level',
    sorter: true,
    filters: [
      {text:'print', value: 'print'},
      {text:'info', value:'info'},
      {text:'error', value:'error'},
      {text:'warn', value:'warn'},
      {text:'paince', value:'panic'}
    ],
    render: (_, { level }) => {
          return (
            <Tag color={levelForColor[level]} key={levelForColor[level]}>
              {level.toUpperCase()}
            </Tag>
          );
    },
  },
  {
    title: '操作IP',
    dataIndex: 'ip',
    key: 'ip',
    render: (text) => <a>{text}</a>,
  },
  {
    title: '操作人',
    dataIndex: 'operator',
    key: 'operator',
  },
  {
    title: '信息',
    dataIndex: 'message',
    key: 'message',
  },
  {
    title: '时间',
    dataIndex: 'date',
    key: 'date',
    render: (text) => moment(text * 1000).format('YYYY-MM-DD HH:mm:ss'),
  },
  {
    title: '操作',
    key: 'action',
    render: (_, record) => (
      <Space size="middle">
        <Button onClick={() => Log.handleRemove(record)}>删除</Button>
      </Space>
    ),
  },
];



const getNewParams = (params) =>({   
  pageNumber: params.pagination && params.pagination.pageSize,
  page: params.pagination && params.pagination.current,
  filters: params.filters
}); 

const Log= ()=>{
    const [selectedRowKeys, setSelectedRowKeys] = useState(() => []);
    const [form] = Form.useForm(); 
    const {token} = theme.useToken();
    const [expand, setExpand] = useState(false); 
    const navigate=useNavigate();
    const [hasData,setHasData]= useState(false); 
    const [data, setData] = useState([]); 
    const [loading, setLoading] = useState(false); 
    const [tableParams, setTableParams] = useState({ 
      pagination: {
        current: 1, 
        pageNumber: 5,  
      },
      filters: {
      } 
    })

    const fetchData = async(values) => { 
      console.log("execute fetchData function")
      setLoading(true);
      setHasData(false);
      try {

        console.log(getNewParams({...tableParams, filters: values}))
        console.log(values)
        
        const response = await axios.get("/log/query", {
          params: getNewParams({...tableParams, filters: values}),

        });

        const {state, message: resMessage} = response.data;

        if (!response.data.state) {
          let Message = resMessage;
          if (Message == undefined || message == "") {
              Message = "系统错误";
          }
          message.error(Message).then(()=>{
              if (response.status==304){
                  navigate('/');
              }
          });
          return
        }
        setData(response.data.data.data);
        setLoading(false);
        setHasData(true);
        setTableParams({
          ...tableParams,
          pagination: {
            ...tableParams.pagination,
            total: response.data.data.total,
          },
        });
      }catch (error){  
        setLoading(false);
        console.log(error); 
      }
    };

    const handleRemove = async () => {
      try {
        const response = await axios.post('/log/remove', selectedRowKeys); 
        if (response.data.success) {
          fetchData(); 
          message.success('删除成功');
        } else {
          message.error(response.data.message || '删除失败');
        }
      } catch (error) {
        message.error('删除失败');
      }
    };


    useEffect(() => {
      fetchData();
    },[JSON.stringify(tableParams.pagination)]);


    const tableParamsChange = (pagination, filters, sorter) => {

      if (filters.level && filters.level.length > 0) {
        setTableParams({
          pagination,
          filters: { level: filters.level[0] }, 
          ...sorter,
        });
      } else {
        setTableParams({
          pagination,
          filters: {}, 
          ...sorter,
        });
      }

      if (tableParams.pagination.pageNumber !== (tableParams.pagination && tableParams.pagination.pageNumber)) {
        setData([]);
      };
    };


    const onSelectChange = (selectedRowKeys) => {
      setSelectedRowKeys(selectedRowKeys);
      console.log(selectedRowKeys);
    };

    const rowSelection = {
      selectedRowKeys,
      onChange: onSelectChange,
    };
    const hasSelected = selectedRowKeys.length > 0;

    const formStyle = {  
      maxWidth: 'none', 
      background: token.colorFillAlter,
      borderRadius: token.borderRadiusLG,
      padding: 18,
    };

    const getFields = () => {
      const children = []
      children.push(
        <Col span={8} >
          <Form.Item
            name = {`message`}
            label = {`message`}
            rules={[
              {
                required:true,
                message:"Input something!",
              },
            ]}
          >
            <Input placeholder="placeholder"/>
          </Form.Item>

          <Form.Item
            name = {`ip`}
            label = {`ip`}
            rules={[
              {
                required:true,
                message:"Input something!",
              },
            ]}
          >
            <Input placeholder="placeholder"/>
          </Form.Item>

          <Form.Item
            name = {`operator`}
            label = {`operator`}
            rules={[
              {
                required:true,
                message:"Input something!",
              },
            ]}
          >
            <Input placeholder="placeholder"/>
          </Form.Item>

          <Form.Item
            name = {`level`}
            label = {`level`}
            rules={[
              {
                required:true,
                message:"Select something!",
              },
            ]}
            initialValue="print"
          > 
            <Select>
              <Option value="print">print</Option>
              <Option value="warn">warn</Option>
              <Option value="error">error</Option>
              <Option value="info">info</Option>
              <Option value="panic">panic</Option>
            </Select>
          </Form.Item>
        </Col> 
      );
      return children
    };


    const onFinish = (values) => {
      console.log('Success:', values);
      setTableParams({
        ...tableParams,
        filters: {...tableParams.filters, ...values}, 
      });
      fetchData({...tableParams.filters, ...values}); 
    };

    return (
      hasData? (
      <div>

        <div
          style={{
            marginBottom: 16,
          }}
        >
          <Button type="primary" onClick={handleRemove} disabled={!hasSelected} loading={loading}>
            Remove
          </Button>

          <span
            style={{
              marginLeft: 8,
            }}
          >
            {hasSelected ? `Selected ${selectedRowKeys.length} items` : ''}
          </span>
        </div>

        <Table 
        columns={columns} 
        dataSource={data}
        pagination={tableParams.pagination}
        loading={loading}
        onChange={tableParamsChange}
        rowSelection={rowSelection} 
        />

        <Form form={form} name="advanced_search" style={formStyle} onFinish={onFinish}>
          <Row gutter={24}>{getFields()}</Row>

          <div
            style={{
              textAlign: 'right',//设置右对齐
            }}
          >
            <Space size="small">
              <Button type="primary" htmlType="submit">
                Search
              </Button>
              <Button
                onClick={() => {
                  form.resetFields();
                }}
              >
                Clear
              </Button>
              <a
                style={{
                  fontSize: 12,
                }}
                onClick={() => {
                  setExpand(!expand);
                }}
              >
              </a>
            </Space>
          </div>
        </Form>
      </div>
    ) : null);
};
export default Log

