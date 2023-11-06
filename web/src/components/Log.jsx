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
        {/* <a>Invite {record.name}</a> */}
        {/* <a>Delete</a> */}
        <Button onClick={() => Log.handleRemove(record)}>删除</Button>
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



//定义一个函数getNewParams，函数接收一个参数params，并且返回一个新对象
//这个新对象包含params对象中的所有属性，并且额外添加了pageNumber和page属性
//具体规则为：如果params对象中存在pagination属性，并且pagination属性中有current和pageSize属性
//则新对象的pageNumber属性的值将设置为params.pagination.pageSize的值
//则新对象的page属性的值将设置为params.pagination.current的值
//并且通过展开语法...params，新对象会将params对象中所有属性添加到新的对象中
const getNewParams = (params) =>({   
  pageNumber: params.pagination && params.pagination.pageSize,
  page: params.pagination && params.pagination.current,
  //...params,
  filters: params.filters
}); 

// const getNewParams2 = (params) => ({
//   pageNumber: params.pagination && params.pagination.pageSize,
//   page: params.pagination && params.pagination.current,
//   filters: Object.entries(params.filters)
//     .map(([key, value]) => `${key}=${encodeURIComponent(value)}`)
//     .join(','),
// });

//Log函数组件
const Log= ()=>{
    //使用useState钩子创建一个状态变量selectedRowKeys，用于存储用户选择的行的键值
    const [selectedRowKeys, setSelectedRowKeys] = useState(() => []);
    //使用useForm钩子创建一个表单实例from，用于处理表单数据
    const [form] = Form.useForm(); 
    //使用theme从上下文中获取token全局身份验证令牌
    const {token} = theme.useToken();
    //使用useState钩子创建一个状态变量expend，用于控制是否展开显示更多内容
    const [expand, setExpand] = useState(false); 
    //使用useNavigate钩子创建一个导航函数navigate，用于在组件内进行页面导航
    const navigate=useNavigate();
    //使用useState钩子创建一个状态变量hasData，用于标识数据是否可供显示
    const [hasData,setHasData]= useState(false); 
    //使用useState钩子创建一个状态变量data，用于存储后端传输过来的日志数据
    const [data, setData] = useState([]); 
    //使用useState钩子创建一个状态变量loading，用于标识数据是否正在加载中
    const [loading, setLoading] = useState(false); 
    //使用useState钩子创建一个状态变量tableParams，初始值为一个对象，其中包含一个pagination子对象
    const [tableParams, setTableParams] = useState({ 
      pagination: {
        current: 1, //当前页码
        pageNumber: 5,  //每页显示条目数量
      },
      filters: {
      } // 添加一个空的过滤器
    })

    //fetchData异步函数处理日志获取与显示
    const fetchData = async(values) => { 
      console.log("execute fetchData function")
      //设置当前为数据加载中
      setLoading(true);
      //设置当前没有日志数据
      setHasData(false);
      try {
        //setLoading(true);    2
        //setSelectedRowKeys([]);   3

        //axios向后端发送get请求到/log/query接口，并且传递tableParams作为参数

        //const encodedParams = qs.stringify(tableParams, { indices: false });
        console.log(getNewParams({...tableParams, filters: values}))
        console.log(values)
        
        const response = await axios.get("/log/query", {
          params: getNewParams({...tableParams, filters: values}), // 将合并后的 filters 赋值给 tableParams
          //params: tableParams, filters: values
        });
        //const response = await axios.get(`/log/query?${encodedParams}`);

        //从后端返回的相应数据中，结构出state和message字段，并将message赋值给resMessage变量
        const {state, message: resMessage} = response.data;
        // console.log(response.data.state)   4

        if (!response.data.state) {
          //登录失败
          console.log("ERROR!!!")
          //创建一个名为Message的变量，并将resMessage赋值给它
          let Message = resMessage;
          //如果Message是未定义或者为空字符串，将Message赋值为"系统错误"
          if (Message == undefined || message == "") {
              Message = "系统错误";
          }
          //显示一个错误提示框，并且在用户关闭提示框之后执行下面代码块
          //如果响应状态码为304，即未修改，使用navigate函数进行页面导航到'/'路径。
          message.error(Message).then(()=>{
              if (response.status==304){
                  navigate('/');
              }
          });
          return
        }
        //console.log(response.data)
        //console.log(response.data.data)
        //将响应数据中的data字段的data属性赋值给data状态变量
        setData(response.data.data.data);
        //console.log(data)
        //将loading状态变量设置为false，表示数据加载完成
        setLoading(false);
        //将hasData状态变量设置为true，表示已有数据
        setHasData(true);
        //？？？？？？
        //使用展开语法更新tableParams状态变量，将其中的pagination属性更新为一个新对象
        //并且设置total属性为后端响应数据中的data字段的total属性
        setTableParams({
          ...tableParams,
          pagination: {
            ...tableParams.pagination,
            total: response.data.data.total,
          },
        });
      }catch (error){  //这个可以进行修改
        setLoading(false);
        console.log(error); //这里可以加一个提示框显示查询报错
      }
    };

    //handleRemove异步函数处理日志删除
    const handleRemove = async () => {
      try {
        //向后端的/log/remove接口发送post请求，携带参数为要删除的日志id
        const response = await axios.post('/log/remove', selectedRowKeys); 
        if (response.data.success) {
          fetchData(); // 删除成功后重新加载数据
          message.success('删除成功');
        } else {
          message.error(response.data.message || '删除失败');
        }
      } catch (error) {
        message.error('删除失败');
      }
    };

    //useEffect钩子函数
    //在组件加载完成之后调用fetchData日志获取显示函数
    //并且只有当tableParams.pagination对象发生便哈才执行
    useEffect(() => {
      fetchData();
    },[JSON.stringify(tableParams.pagination)]);


    //tableParamsChange函数用于处理表格参数的变化
    //该函数接收三个参数，pagination, filters, sorter
    //三个参数用于更新表格的分页、筛选、排序参数
    const tableParamsChange = (pagination, filters, sorter) => {
      //？？？？？
      //如果filters.level存在且长度>0，表明有筛选条件
      //此时设置filters.level[0]为level参数的值
      //如果不存在筛选条件，则将level参数的值设置为print
      if (filters.level && filters.level.length > 0) {
        setTableParams({
          pagination,
          filters: { level: filters.level[0] }, // 设置过滤器
          ...sorter,
        });
      } else {
        setTableParams({
          pagination,
          filters: {}, // 如果没有选择过滤器，则清空过滤器
          ...sorter,
        });
      }
      //？？？？？？
      //如果分页的参数发生了变化，就清空当前表格的数据
      if (tableParams.pagination.pageNumber !== (tableParams.pagination && tableParams.pagination.pageNumber)) {
        setData([]);
      };
    };

    //onSelectchange函数用于更新选中的表格的行
    //传入选中的行
    const onSelectChange = (selectedRowKeys) => {
      setSelectedRowKeys(selectedRowKeys);
      console.log(selectedRowKeys);
    };
    
    // rowSelection对象，包含两个属性：
    // selectedRowKeys：指定当前选中的行
    // onChange：指定onSelectChange函数作为选中行发生变化时执行的回调函数。
    const rowSelection = {
      selectedRowKeys,
      onChange: onSelectChange,
    };

    //hasSelected对象，如果选中行>0，赋值为true，表示有选中的行
    const hasSelected = selectedRowKeys.length > 0;

    //formStyle对象，用于设置表单样式
    const formStyle = {  
      //最大宽度
      maxWidth: 'none', 
      //背景颜色
      background: token.colorFillAlter,
      //圆角半径
      borderRadius: token.borderRadiusLG,
      //内边距
      padding: 18,
    };
    // const formStyle = {
    //   position: 'fixed',
    //   top: 20,
    //   right: 20,
    //   maxWidth: 'none',
    //   background: token.colorFillAlter,
    //   borderRadius: token.borderRadiusLG,
    //   padding: 12,
    // };
    
    //getFields函数，用于生成表单的字段内容
    const getFields = () => {
      //在函数内部创建一个空数组children，用于存储每个字段的内容
      const children = []
      //使用children.push方法将各个表单字段添加到children数组中。
      //每个表单字段使用<Form.Item>组件来包裹，并设置相应属性
      children.push(
        <Col span={8} >
          <Form.Item
            //name指定字段名称
            name = {`message`}
            //label指定字段标签
            label = {`message`}
            //rules指定验证规则
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

    //?????
    // 函数onFinish用于在搜索框中的表单提交的时候被调用，并将搜索参数传给fetchData函数
    const onFinish = (values) => {
      console.log('Success:', values);
      setTableParams({
        ...tableParams,
        filters: {...tableParams.filters, ...values}, // 合并 filters 和表单数据
      });
      fetchData({...tableParams.filters, ...values}); // 将合并后的 filters 和表单数据传递给 fetchData 函数
    };

    return (
      //如果hasData为true，将会执行后面的代码块
      hasData? (
      <div>

        <div
          style={{
            marginBottom: 16,
          }}
        >
                {/* // type="primary"：设置按钮的类型为主要按钮
                // onClick={handleRemove}：指定点击事件处理函数为handleRemove
                // disabled={!hasSelected}：根据hasSelected的值来决定是否禁用按钮。当没有选择行时，按钮将被禁用
                // loading={loading}：根据loading的值来决定按钮是否处于加载状态 */}
          <Button type="primary" onClick={handleRemove} disabled={!hasSelected} loading={loading}>
            Remove
          </Button>

          <span
            style={{
              //设置左边距为8像素
              marginLeft: 8,
            }}
          >
            {/*如果有选中的行，则显示选中的行数*/ }
            {hasSelected ? `Selected ${selectedRowKeys.length} items` : ''}
          </span>
        </div>

        <Table 
        //指定表格的列配置
        columns={columns} 
        //指定表格的数据源
        dataSource={data}
        //指定分页配置
        pagination={tableParams.pagination}
        //指定显示加载状态
        loading={loading}
        //指定表格发生变化时的处理函数
        onChange={tableParamsChange}
        //指定启用选择行功能
        rowSelection={rowSelection} 
        />

        {/*
        创建一个表单组件Form，并将表单实例form与Form组件进行绑定。
        name：设置表单的名称为“advanced_search”
        style：设置表单的样式
        onFinish：指定表单提交时候的回调函数
        */ }
        <Form form={form} name="advanced_search" style={formStyle} onFinish={onFinish}>

          {/*
          使用Row组件来包裹表单字段。gutter属性设置了字段之间的间距为24位
          通过调用getFileds()函数，在Row组件中渲染表单字段
          */}
          <Row gutter={24}>{getFields()}</Row>

          <div
            style={{
              textAlign: 'right',//设置右对齐
            }}
          >

          {/*
          使用Space组件，用于在内部的元素之间添加间距
          size属性设置了间距大小为“small”
          */}
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

