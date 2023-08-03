
import axios from 'axios';
const backendURL = "http://localhost:8080";
const instance =axios.create({
    baseURL: backendURL,
    timeout: 5000,
  });
  


axios.interceptors.response.use(null,(error)=>{
      if (error.response.state ==304){
        sessionStorage.removeItem("repliteweb");
        // navigate("/login");
      }
       message.error("系统出错啦.....");
       console.log(error);
       return Promise.reject(error);
    })