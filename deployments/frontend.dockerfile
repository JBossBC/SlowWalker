# 使用 Node.js 镜像作为基础镜像
FROM node:alpine AS build

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 package-lock.json 到容器
COPY web/package*.json ./

# 安装项目依赖
RUN npm install

# 将项目文件复制到容器
COPY ./web/ .

# 构建 React 项目
RUN npm run build

# 使用 Nginx 镜像作为基础镜像
FROM nginx:alpine

# 将 Nginx 配置文件复制到容器
COPY deployments/nginx.conf /etc/nginx/conf.d/default.conf

# 从前一个构建阶段将构建后的 React 项目复制到 Nginx 默认目录下
COPY --from=build /app/build /usr/share/nginx/html

# 暴露 Nginx 默认端口
EXPOSE 80

# 启动 Nginx 服务
CMD ["nginx", "-g", "daemon off;"]
