import axios from '../utils/axios';
import React,{useState, useEffect} from 'react';
import { Empty, Space, Table, Tag,Button,message, Form, Col, Input, Row, Select, theme } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import moment from 'moment';
import { useNavigate } from 'react-router';

/** 	PRINT LogLevel = "print"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
	INFO  LogLevel = "info"
	PANIC LogLevel = "panic"
    */
const  levelForColor = {"print":"success","info":"blue","error":"error","warn":"warning","panic":"red"}
const { Option } = Select;
const columns = [
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
        {/* <a>Invite {record.name}</a> */}
        {/* <a>Delete</a> */}
        <Button onClick={() => Log.handleRemove(record)}>删除</Button>
      </Space>
    ),
  },
];

const getNewParams = (params) =>({   //当变量参数更新时修改参数
  pageNumber: params.pagination && params.pagination.pageSize,
  page: params.pagination && params.pagination.current,
  ...params
}); 

const Log= ()=>{
    const [selectedRowKeys, setSelectedRowKeys] = useState([]);//设置批量选择的数组,selectedRowKeys是状态变量本身，setSelectedRowKeys是对状态变量进行设置
    const [form] = Form.useForm();
    const {token} = theme.useToken();
    const [expand, setExpand] = useState(false); 
    const navigate=useNavigate();
    const [hasData,setHasData]= useState(false); 
    const [data, setData] = useState([]); //默认data数据为空数组
    const [loading, setLoading] = useState(false); //调用fetchlog函数加载数据时为true，加载完毕为false
    const [tableParams, setTableParams] = useState({  //设置变量参数
      pagination: {
        current: 1,
        pageNumber: 5,
      }
    })
    const fetchData = async() => {
      console.log("log1")
      setLoading(true);//排除嫌疑
      setHasData(false);//排除嫌疑
      try {
        setLoading(true); //当调用fetchData这个异步函数时，表示要进行数据的重载
        setSelectedRowKeys([]); //默认选择为空
        const response = await axios.get("/log/query", {
          params: getNewParams(tableParams),
        });
        const {state, message: resMessage} = response.data;
        console.log(response.data.state)
        if (!response.data.state) {
          //登录失败
          console.log("log2")
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
        console.log(response.data)
        console.log(response.data.data)
        setData(response.data.data.data);
        console.log(data)

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
        console.log(error); //这里可以加一个提示框显示查询报错
      }
    };

    const handleRemove = async () => {
      try {
        const response = await axios.post('/log/remove', selectedRowKeys); //直接传递数组过去
        if (response.data.success) {
          //fetchData(); // 删除成功后重新加载数据
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
          level: filters.level[0],
          ...sorter,
        });
      } else {
        setTableParams({
          pagination,
          level: 'print',
          ...sorter,
        });
      }
      if (tableParams.pagination.pageNumber !== (tableParams.pagination && tableParams.pagination.pageNumber)) {
        setData([]);
      };
    };

    const onSelectchange = (selectedRowKeys) => {
      setSelectedRowKeys([...selectedRowKeys]);
    };
    const rowSelection = {
      selectedRowKeys,
      onChange:onSelectchange,
    }
    const hasSelected = selectedRowKeys.length > 0;

    return (
      hasData? (
      <div>
        <div
          style={{
            marginBottom: 16,
            display: 'flex',
            justifyContent: 'space-between',  // 添加这一行
          }}
        >
          <Button type="primary" onClick={handleRemove} disabled={selectedRowKeys.length === 0} loading={loading}>
            Remove
          </Button>
          <span>
            {hasSelected ? `Selected ${selectedRowKeys.length} items` : ''}
          </span>
        </div>
        <div style={{ marginBottom: 16 }}>
           <Search />
        </div>
        <Table 
        columns={columns} 
        dataSource={data}
        pagination={tableParams.pagination}
        loading={loading}
        onChange={tableParamsChange}
        rowSelection={rowSelection} // 添加rowSelection属性，启用选择功能
        />
      </div>
    ) :null
  );
};

const Search =(fetchData) => {
  const [form] = Form.useForm(); // 定义 form 变量
  const [expand, setExpand] = useState(false); 
  const {token} = theme.useToken();
  //以下为处理搜索的：
  const formStyle = {   //设置样式
    maxWidth: '600px',
    background: token.colorFillAlter,
    borderRadius: token.borderRadiusLG,
    padding: 10,
  };
  const containerStyle = {
    display: 'flex',
    justifyContent: 'flex-end',
  };
  const getFields = () => {
    return (
      <Row gutter={12} justify="end">
        <Col span={12}>
          <Form.Item
            name={`message`}
            label={`message`}
            rules={[
              {
                required: false,
                message: "Input something!",
              },
            ]}
          >
            <Input placeholder="placeholder" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            name={`ip`}
            label={`ip`}
            rules={[
              {
                required: false,
                message: "Input something!",
              },
            ]}
          >
            <Input placeholder="placeholder" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            name={`operator`}
            label={`operator`}
            rules={[
              {
                required: false,
                message: "Input something!",
              },
            ]}
          >
            <Input placeholder="placeholder" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            name={`level`}
            label={`level`}
            rules={[
              {
                required: false,
                message: "Select something!",
              },
            ]}
            initialValue="1"
          >
            <Select>
              <Option value="1">print</Option>
              <Option value="2">warn</Option>
              <Option value="3">error</Option>
              <Option value="4">info</Option>
              <Option value="5">panic</Option>
            </Select>
          </Form.Item>
        </Col>
      </Row>
    );
  };

  const onFinish = (values) => {
    fetchData(values); // 调用fetchData函数，并传递搜索参数
    console.log('Received values of form: ', values);
  }

  return(
    <div style={containerStyle}>
      <div style={formStyle}>
        <Form form={form} name="advanced_search" onFinish={onFinish}>
          <Row gutter={12}>{getFields()}</Row>
          <div
            style={{
              textAlign: 'right',
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
            <DownOutlined rotate={expand ? 130 : 0} /> Collapse
                </a>
          </Space>
          </div>
        </Form>
      </div>
    </div>
  )
};
export default Log


