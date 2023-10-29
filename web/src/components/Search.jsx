import axios from '../utils/axios';
import React, { useState } from 'react';
import { Input, Button, List, message, Select } from 'antd';

const { Option } = Select;

const Search = () => {
  const [searchDesc, setSearchDesc] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [selectedTags, setSelectedTags] = useState([]);

  const handleSearch = async () => {
    try {
      const response = await axios.get("/search/function", {
        params: {
          labels: selectedTags,
          descriptions: searchDesc
        },
      });
      const { state, message: resMessage } = response.data;
      if (!state) {
        let Message = resMessage;
        if (Message === undefined || Message === "") {
          Message = "系统错误";
        }
        message.error(Message);
        return;
      }

      setSearchResults(response.data.data.hits);
    } catch (error) {
      console.log(error);
      message.error("查询错误");
    }
  };

  return (
    <div>
      <Select
        mode="multiple"
        style={{ width: '100%', marginBottom: 12 }}
        placeholder="请选择类型"
        value={selectedTags}
        onChange={(value) => setSelectedTags(value)}
      >
        <Option value="golang_scripts">golang_scripts</Option>
        <Option value="python_scripts">python_scripts</Option>
        <Option value="fileSystem">fileSystem</Option>
      </Select>

      <Input
        value={searchDesc}
        onChange={e => setSearchDesc(e.target.value)}
        placeholder="请输入关键词"
        style={{ marginBottom: 12 }}
      />

      <Button type="primary" onClick={handleSearch} style={{ marginBottom: 12 }}>
        搜索
      </Button>

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