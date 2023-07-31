import axios from 'axios';
import React,{useState, useEffect,useContext} from 'react';
import { Empty, Space, Table, Tag,Button } from 'antd';
import moment from 'moment';
import {Backend} from "../App"
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


const getNewParams = (params) =>({   //当变量参数更新时修改参数
  pageNumber: params.pagination && params.pagination.pageSize,
  page: params.pagination && params.pagination.current,
  ...params
}); 

const Log= ()=>{
    const [hasData,setHasData]= useState(false); 
    const backendURL = useContext(Backend);
    const [data, setData] = useState([]); //默认data数据为空数组
    const [loading, setLoading] = useState(false); //调用fetchlog函数加载数据时为true，加载完毕为false
    const [tableParams, setTableParams] = useState({  //设置变量参数
      pagination: {
        current: 1,
        pageNumber: 5,
      }
    })
    
    const fetchData = async() => {
      try {
        setLoading(true); //当调用fetchData这个异步函数时，表示要进行数据的重载
        const response = await axios.get(backendURL+ "/log/query", {
          params: getNewParams(tableParams),
          headers:{
            'Authorization':'Bearer '+sessionStorage.getItem("repliteweb")
          }
        });
        console.log(tableParams)
        console.log(response)
        console.log(response.data)
        console.log()
        if (response.status == '200') {
          setData(response.data.data);
          console.log("there")
          console.log(data)
          setLoading(false);
          setHasData(true);
          setTableParams({
            ...tableParams,
            pagination: {
              ...tableParams.pagination,
              total: response.data.total,
            },
          });
        }
      }catch (error){
        setLoading(false);
        console.log(error); //这里可以加一个提示框显示查询报错
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

    return(hasData?
    <Table 
      columns={columns} 
      dataSource={data}
      pagination={tableParams.pagination}
      loading={loading}
      onChange={tableParamsChange}
    />:<Empty/>);
};
export default Log