import React,{useEffect, useState}from "react"
import { Table,Menu,message,Button, Modal,Popconfirm, Form, Input} from 'antd'
import { TeamOutlined,WarningOutlined } from "@ant-design/icons";
import axios from "../utils/axios";
import moment from "moment";
import DynamicComponent from "./parseReact";


const UserManage = (props)=>{
    const [department,setDepartment] =useState([]);
    const [users,setUsers]=useState([]);
    const [updateModal,setUpdateModal]=useState(false);
    const [updateInfo,setUpdateInfo]=useState(null);
    const [deleteInfo,setDeleteInfo]=useState(null);
    const [deleteConfirm,setDeleteConfirm]=useState(false);
    const [deleteSend,setDeleteSend] = useState(false);
    const columns=[{
        title:"账号",
        dataIndex: 'username',
        key: 'username',
    },{
        title:"真实姓名",
        dataIndex:"realName",
        key:"realName",
    },{
        title:"联系电话",
        dataIndex:"phoneNumber",
        key:"phoneNumber",
    },{
        tilte:"所属部门",
        dataIndex:"department",
        key:"belongs_department",
    },{
        title:"注册时间",
        dataIndex:"createTime",
        key:"createTime",
        render:(value)=>{
            return  moment(value * 1000).format('YYYY-MM-DD HH:mm:ss');  
        },
    },{
        tilte:"操作",
        dataIndex:"operation",
        key:"operation",
        render:function(value,record,index){
            return (<div><Button username={record.username} userIndex={index} onClick={()=>{updateUserInfoClick(index)}}>修改</Button>
            <Popconfirm
        title="删除"
      description={deleteInfo!=undefined&&`是否要删除用户: ${deleteInfo.username}`}
      open={deleteConfirm}
      cancelText="取消"
      okText="确定"
      icon={<WarningOutlined />}
      onConfirm={deleteUser}
      okButtonProps={{
        loading: deleteSend,
      }}
      onCancel={()=>{setDeleteConfirm(false)}}
    >
            <Button  index={index} onClick={()=>{delConfirm(index)}}>删除</Button>
            </Popconfirm></div>)
        },
    }];
    function deleteUser(){

    }
    function delConfirm(index){
        setDeleteInfo(users[index]);
        setDeleteConfirm(true);
    }
    function updateUserInfoClick(targetIndex){
        setUpdateInfo(users[targetIndex]);
        setUpdateModal(true);
    }
    function getItem(key, icon, children, label, type) {
        return {
          key,
          icon,
          children,
          label,
          type,
        };
    };

     function getAllDepartments(){
       axios.get('/department/querys').then((response)=>{
            let data = response.data;
            if(data.state!==true){
                let msg =data.message;
                if(msg == undefined || msg == null){
                    msg = "请求失败";
                }
              message.error(msg);
              return 
            }
            return data.data;
        }).then((departments)=>{
            if(departments === undefined || departments == null){
                return
            }
            let departmentView = [];
            for(let i=0;i<departments.length;i++){
                departmentView.push(getItem(departments[i].name,<TeamOutlined/>,null,departments[i].name,null));
            }
            console.log(departmentView);
            setDepartment(departmentView)
        })
    }
    //TODO
    function currentDepartmentHandle(curElement){
           axios.get("/user/filter",{params:{department:curElement.key,page:1,pageNumber:5}}).then((response)=>{
             let data=response.data;
             if(data.state!=true){
                let msg =data.message;
                if (msg == undefined || msg==""){
                    msg="请求失败";
                }
                message.error(msg);
                return;
             }
             setUsers(data.data);
           })
    }
     function init(){
        getAllDepartments();
    }
    useEffect(()=>{
        init();
    },[])
    return(<div style={{width:"100%",height:"100%",display:"flex",justifyContent:"flex-start",alignContent:'center',flexDirection:"row"}}>
        <div style={{width:"20%",height:"100%"}}>
            <Menu items={department} onSelect={currentDepartmentHandle}>        
            </Menu>
        </div>
        <div style={{width:"80%",height:"100%"}}>
            <Table bordered sticky dataSource={users} columns={columns}></Table>
        </div>
        {/* <DynamicComponent componentCode={componentCode} /> */}
        {updateModal&&<Modal title={<div></div>} keyboard  centered width="700px" okText="确定" cancelText="取消" onCancel={()=>{setUpdateModal(false)}} open={updateModal}>
            <Form>
                <Form.Item label="账号">
                <Input  placeholder={updateInfo?updateInfo.username:""} disabled ></Input>
                </Form.Item>
                <Form.Item label="真实姓名">
                <Input placeholder={updateInfo?updateInfo.realName:""}></Input>
                </Form.Item>
                <Form.Item label="电话号码">
                <Input placeholder={updateInfo?updateInfo.phoneNumber:""}></Input>
                </Form.Item>
                <Form.Item label="所属部门">
                <Input  placeholder={updateInfo?updateInfo.department:""}></Input>
                </Form.Item>
                </Form> </Modal>}
    </div>)
}


export default UserManage;