import React,{useEffect, useState}from "react"
import { Table,Menu,message} from 'antd'
import { TeamOutlined } from "@ant-design/icons";
import axios from "../utils/axios";
const UserManage = ()=>{
    const [department,setDepartment] =useState([])
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
        key:"department",
    },{
        tilte:"操作",
        dataIndex:"operation",
        key:"operation",
    }];
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
       axios.get('/department/querys',{params:{"resource":"人员管理"}}).then((response)=>{
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
    function currentDepartmentHandle(curElement){
           axios.get("/user/filter")
    }
     function init(){
        getAllDepartments();
    }
    init();
    const dataSource = [];
    return(<div style={{width:"100%",height:"100%",display:"flex",justifyContent:"flex-start",alignContent:'center',flexDirection:"row"}}>
        <div style={{width:"20%",height:"100%"}}>
            <Menu items={department} onSelect={currentDepartmentHandle}>        
            </Menu>
        </div>
        <div style={{width:"50%",height:"100%"}}>
            <Table columns={columns}></Table>
        </div>
    </div>)
}


export default UserManage;