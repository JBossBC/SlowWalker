import React,{useEffect, useState}from "react"
import { Table,Menu,message,Button, Modal,Popconfirm} from 'antd'
import { TeamOutlined,WarningOutlined } from "@ant-design/icons";
import axios from "../utils/axios";
import moment from "moment";
const UserManage = ()=>{
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
            return (<div><Button username={record.username} index={index} onClick={(target)=>{updateUserInfoClick(target)}}>修改</Button>
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
    function updateUserInfoClick(target){
        setUpdateInfo(users[target.target.index]);
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
    const dataSource = [];
    return(<div style={{width:"100%",height:"100%",display:"flex",justifyContent:"flex-start",alignContent:'center',flexDirection:"row"}}>
        <div style={{width:"20%",height:"100%"}}>
            <Menu items={department} onSelect={currentDepartmentHandle}>        
            </Menu>
        </div>
        <div style={{width:"80%",height:"100%"}}>
            <Table dataSource={users} columns={columns}></Table>
        </div>
        {updateModal&&<Modal>
            
            </Modal>}
    </div>)
}


export default UserManage;