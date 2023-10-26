import React from "react"
import { Table,Menu} from 'antd'
const UserManage = ()=>{
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
    const dataSource = [];
    return(<div style={{width:"100%",height:"100%",display:"flex",justifyContent:"flex-start",alignContent:'center',flexDirection:"column"}}>
        <div>
            <Menu>
                
            </Menu>
        </div>
        <div>
   <Table columns={columns}></Table>
        </div>
    </div>)
}


export default UserManage;