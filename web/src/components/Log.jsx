

import React,{useState} from 'react';
import { Empty, Space, Table, Tag,Button } from 'antd';
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
  // {
  //   title: '操作',
  //   key: 'action',
  //   render: (_, record) => (
  //     <Space size="middle">
  //       {/* <a>Invite {record.name}</a> */}
  //       {/* <a>Delete</a> */}
  //       <Button>删除</Button>
  //     </Space>
  //   ),
  // },
];
const data = [
  {
    key: '1',
    name: 'John Brown',
    age: 32,
    address: 'New York No. 1 Lake Park',
    tags: ['nice', 'developer'],
  },
  {
    key: '2',
    name: 'Jim Green',
    age: 42,
    address: 'London No. 1 Lake Park',
    tags: ['loser'],
  },
  {
    key: '3',
    name: 'Joe Black',
    age: 32,
    address: 'Sydney No. 1 Lake Park',
    tags: ['cool', 'teacher'],
  },
];
const Log= ()=>{
    const [hasData,setHasData]= useState(false);
    
    return(hasData?<Table columns={columns} dataSource={data} />:<Empty/>)
}





export default Log