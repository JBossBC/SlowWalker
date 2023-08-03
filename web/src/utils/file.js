import {message} from 'antd';
// 断点续传
async function uploadFile(url, file) {

    const fileSize = file.size;
    const CHUNK_SIZE = 5 * 1024 * 1024; // 每个分片的大小为 5MB
    let startByte = 0;
    let endByte = CHUNK_SIZE - 1;
    let chunkIndex = 0;
    while (startByte < fileSize) {
      const chunk = file.slice(startByte, endByte + 1);
      const formData = new FormData();
      formData.append('chunk', chunk);
      const headers = {
        'Content-Range': `bytes ${startByte}-${endByte}/${fileSize}`,
      };
      const response = await fetch(url, {
        method: 'POST',
        body: formData,
        headers: headers,
      });
  
      if (response.ok) {
        console.log(`Uploaded chunk ${chunkIndex + 1}`);
        if (endByte === fileSize - 1) {
          // 已上传完整文件
          message.success("文件上传完成");
        } else {
          // 继续上传下一片段
          startByte = endByte + 1;
          endByte = Math.min(endByte + CHUNK_SIZE, fileSize - 1);
          chunkIndex++;
        }
      }
    }
  }

  export default uploadFile;