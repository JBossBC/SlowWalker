import React,{useState, useContext, useEffect} from 'react';
import { Empty, Space, Table, Tag,Button } from 'antd';
import { Input , List } from 'antd';
import axios from 'axios';
import {Backend} from "../App";
import { Select } from 'antd';
const { Option } = Select;
/** 	PRINT LogLevel = "print"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
	INFO  LogLevel = "info"
	PANIC LogLevel = "panic"
    */
const  levelForColor = {"print":"success","info":"blue","error":"error","warn":"warning","panic":"red"}

const columns = [
    {
      title: '日志等级',
      key: 'level',
      dataIndex: 'level',
      sorter: true,   //加一个sorter用于实现排序
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
  },
  {
    title: '操作',
    key: 'action',
    render: (_, record) => (
      <Space size="middle">
        {/* <a>Invite {record.name}</a> */}
        {/* <a>Delete</a> */}
        <Button>删除</Button>
      </Space>
    ),
  },
];
// const data = [
//   {
//     key: '1',
//     name: 'John Brown',
//     age: 32,
//     address: 'New York No. 1 Lake Park',
//     tags: ['nice', 'developer'],
//   },
//   {
//     key: '2',
//     name: 'Jim Green',
//     age: 42,
//     address: 'London No. 1 Lake Park',
//     tags: ['loser'],
//   },
//   {
//     key: '3',
//     name: 'Joe Black',
//     age: 32,
//     address: 'Sydney No. 1 Lake Park',
//     tags: ['cool', 'teacher'],
//   },
// ];


const getRandomuserParams = (params) => ({  
  pageNumber: params.pagination && params.pagination.pageSize,
  page: params.pagination && params.pagination.page,
  ...params,
});

const Log= ()=>{
  const backendURL = useContext(Backend);
  const [data, setData] = useState([]); //存储后端返回的数据
  const [hasData,setHasData]= useState(false);  //判断是否从后端取到了日志数据,有数据为true，没有数据为false
  const [loading, setLoading] = useState(false);//显示正在加载中（转圈圈），当要加载其他页时，会有转圈圈效果，加载完成救变为
  const [tableParams, setTableParams] = useState({ //核心部分，设置当前页和每页最大条目数
    pagination: {
      page: 1,
      pageSize: 3,
      total:10,
    },
  });
  const [currentPage, setCurrentPage] = useState(1);  //页码控制
  const [pageSize, setPageSize] = useState(3);
  //setLoading(true); //表示数据加载中
  console.log(tableParams)
  const fetchLog = async() => {
      try {
          const response = await axios.get("http://localhost:8080/log/query",{
            params: {
              page : currentPage,
              pageNumber : getRandomuserParams(tableParams).pageNumber
            }, 
            headers:{
                'Authorization': 'Bearer '+ sessionStorage.getItem("repliteweb")   //第三个参数是请求头配置，
              }
          } );
          console.log(response.data)
          console.log(response.data.message) //后端把总的页数存在data.message里面
          const result = response.data;
          if (response.status!='200' || data.state!=true) {
            setData(result.data);
            setLoading(false);
            setTableParams({
              ...tableParams,
              pagination: {
                ...tableParams.pagination, 
                total: response.data.message,  
              },
            });
            setHasData(true)
          }
          console.log(result.data.message)
      }catch (error) {
          setLoading(false);
          console.error(error);
      }
    }

    const handleTableChange = (pagination, sorter) => {
      setTableParams({
        ...tableParams,
        pagination: {
          ...tableParams.pagination,
          page: pagination.page,
          pageSize: pagination.pageNumber
        },
        ...sorter
      });
      setCurrentPage(pagination.page);
      setHasData(false); // 重置hasData状态
    };

    const handlePageSizeChange = (value) => {
      setPageSize(value);
      setTableParams({
        ...tableParams,
        pagination: {
          ...tableParams.pagination,
          pageSize: value
        }
      });
    };
    useEffect(() => {
      fetchLog();
    },[JSON.stringify(tableParams)]);

    if (tableParams.pagination.pageSize !== (tableParams.pagination && tableParams.pagination.pageSize)) {  //设置空数据
      setData([]);
    };
    
    return (
      <div>
      <Select defaultValue={tableParams.pagination.pageSize} onChange={handlePageSizeChange}>
      <Option value={3}>3 条/页</Option>
      <Option value={5}>5 条/页</Option>
      <Option value={10}>10 条/页</Option>
      </Select>

      <Table
        columns={columns}  //定义表格的列属性
        dataSource={data} //定义表格的数据源
        pagination = {tableParams.pagination} //分页设置
        loading = {loading}  //是否显示加载中
        onChange = {handleTableChange} //当表格属性改变时调用handleTableChange这个函数
        />
      </div>
    );
  
};


export default Log