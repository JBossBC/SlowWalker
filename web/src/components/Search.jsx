import axios from '../utils/axios';
import React, { useState } from 'react';
import { Input, Button, List, message, Select } from 'antd';

const { Option } = Select;

const Search = () => {
  const [searchText1, setSearchText1] = useState('');
  const [searchText2, setSearchText2] = useState('');
  const [searchResults, setSearchResults] = useState([]);

  const handleSearch = async () => {
    try {
      const response = await axios.get("/search/function", {
        params: {
          labels: searchText1.trim(),
          descriptions: searchText2.trim()
        },
      });
      const { state, message: resMessage } = response.data;
      console.log(response.data.data)
      if (!state) {
        let Message = resMessage;
        if (Message === undefined || Message === "") {
          Message = "系统错误";
        }
        message.error(Message);
        return;
      }

      setSearchResults(response.data.data);
    } catch (error) {
      console.log(error);
      message.error("查询错误");
    }
  };


  const inputsStyle = {
    marginBottom: '12px', // 下边距
    width: '80%', // 宽度
    height: '40px', // 高度
    //padding: '8px', // 内边距
    fontSize: '12px', // 字体大小
    fontWeight: 'bold', // 字体粗细
    color: '#333', // 字体颜色
    backgroundColor: '#fff', // 背景色
    border: '1px solid #ccc', // 边框
    borderRadius: '4px', // 边框圆角
    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.2)', // 阴影
    outline: 'none', // 失去焦点时去掉默认的外部边框
    cursor: 'auto', // 鼠标指针样式
    transition: 'all .3s ease-in-out', // 过渡效果
  };

  return (
    <div>
      <Input
        value={searchText1}
        onChange={e => setSearchText1(e.target.value)}
        placeholder="请输入标签"
        style={inputsStyle}
      />
      <Input
        value={searchText2}
        onChange={e => setSearchText2(e.target.value)}
        placeholder="请输入关键词"
        style={inputsStyle}
      />
      <div style={{ display: 'flex', justifyContent: 'flex-start' }}>
      <Button
        type="primary"
        onClick={handleSearch}
        style={{ marginBottom: 12}}

        ghost={true}
        loading={false}
      >
        搜索
      </Button>
    </div>

      <List
        dataSource={searchResults}
        renderItem={item => (
          <List.Item>
            <List.Item.Meta
              title={item.title}
              description={item.description}
              attribute={item.attribute}
            />
            <div>{item.label}</div>
          </List.Item>
        )}
      />
    </div>
  );
};

export default Search;