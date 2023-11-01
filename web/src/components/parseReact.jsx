import React,{useState,useEffect} from 'react';
import ReactDOM from 'react-dom';
import Babel from 'babel-standalone';

function DynamicComponent(props) {
  const [jsxCode, setJsxCode] = useState('');

  useEffect(() => {
    // 假设从后端获取的数据为：
    const dataFromBackend = `import { Button } from 'antd';

      function dynamicCode() {
        return (
          <div>
            <Button>hello</Button>
          </div>
        );
      }`;

    setJsxCode(dataFromBackend);
  }, []);

  useEffect(() => {
    if (jsxCode) {
     // 后端传入的代码
const code = `import antd from "antd"; function dynamicCode(){ return (<div><Button/>hello</div>); }`;

// 使用babel-standalone解析JSX代码
const transformedCode = Babel.transform(code, {
  presets: ['react']
}).code;

// 创建一个动态执行函数
const execute = new Function('React', 'ReactDOM', transformedCode);

// 执行函数，获取React元素
const element = execute(React, ReactDOM);

// 渲染React元素到DOM节点
ReactDOM.render(element, document.getElementById('dynamicContainer'));

    }
  }, [jsxCode]);

  return (
    <div>
      <div id="dynamicContainer"></div>
    </div>
  );
}

export default DynamicComponent;

