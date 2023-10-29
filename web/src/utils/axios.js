
// import { message } from 'antd';
import axios from 'axios';
const backendURL = "http://localhost:8080";
const instance =axios.create({
    baseURL: backendURL,
    timeout: 5000,
    headers:{
      'Authorization': 'Bearer ' + sessionStorage.getItem('repliteweb'),
    }
  });


//   instance.interceptors.request.use(function(request){
//     let token=sessionStorage.getItem("repliteweb")!=undefined?sessionStorage.getItem("repliteweb"):"";
//     if (token!=""){
//         request.headers.Authorization ='Bearer '+token;
//     }
//     return request;
// },(error)=>Promise.reject(error));

// instance.interceptors.response.use((response)=>{

//     if (response.state ==304){
//         sessionStorage.removeItem("repliteweb");
//         // navigate("/login");
//       }
//     //    message.error("系统出错啦.....");
//     return response;
// },(error)=>{message.error("系统出错啦");return Promise.reject(error)});

 export default instance ;