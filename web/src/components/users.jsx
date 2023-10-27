import React,{useEffect, useState}from "react"
import { Table,Menu,message} from 'antd'
import axios from "axios";
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
    function getItem(label, key, icon, children, type) {
        return {
          key,
          icon,
          children,
          label,
          type,
        };
    };

    async function getAllDepartments(){
       const departments=await axios.get('/department/querys').then((response)=>{
            let data = response.data;
            if(data.state!==true){
                let msg =data.message;
                if(msg == undefined || msg == null){
                    msg = "请求失败";
                }
              message.error(msg);
              return
            }
            return data.data
        })
        return departments
    }

    async function init(){
        let departments= await getAllDepartments()
        setDepartment(departments)

    }
    useEffect(()=>{
        init()
    },[])
    const dataSource = [];
    return(<div style={{width:"100%",height:"100%",display:"flex",justifyContent:"flex-start",alignContent:'center',flexDirection:"column"}}>
        <div>
            <Menu >
               
                
            </Menu>
        </div>
        <div>
            <Table columns={columns}></Table>
        </div>
    </div>)
}


export default UserManage;