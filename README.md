# Traefik 演示项目

这是一个用于学习Traefik反向代理框架的完整演示项目，包含详细的中文注释。

## 🎯 项目目标

- 学习Traefik的基本概念和配置
- 了解反向代理的工作原理
- 掌握Docker与Traefik的集成
- 学习路由规则和中间件的使用

## 📁 项目结构

```
traefikDemo/
├── docker-compose.yml      # Docker Compose配置文件
├── traefik.yml            # Traefik配置文件
├── custom-app/            # 自定义应用（Go或Python）
│   ├── Dockerfile         # 应用Dockerfile
│   ├── main.go           # Go应用代码（Go版本）
│   ├── go.mod            # Go模块配置（Go版本）
│   ├── go.sum            # Go依赖校验文件（Go版本）
│   ├── main.py           # Python应用代码（Python版本）
│   └── requirements.txt  # Python依赖文件（Python版本）
├── custom-app-go/        # Go版本备份
├── custom-app-python/    # Python版本备份
├── start.sh              # 启动脚本
├── switch-to-go.sh       # 切换到Go版本
├── switch-to-python.sh   # 切换到Python版本
└── README.md             # 项目说明文档
```

## 🚀 快速开始

### 前置要求

- Docker 和 Docker Compose
- 确保端口80、443、8080未被占用

### 选择应用版本

项目提供了两个版本的自定义应用：

1. **Go版本**（默认）：使用Gin框架，性能更好
2. **Python版本**：使用Flask框架，开发更简单

#### 切换版本

```bash
# 切换到Python版本
./switch-to-python.sh

# 切换到Go版本
./switch-to-go.sh
```

### 启动服务

1. **克隆或下载项目**
   ```bash
   cd traefikDemo
   ```

2. **使用启动脚本（推荐）**
   ```bash
   ./start.sh
   ```

3. **或手动启动所有服务**
   ```bash
   docker-compose up -d
   ```

4. **查看服务状态**
   ```bash
   docker-compose ps
   ```

5. **查看日志**
   ```bash
   # 查看所有服务日志
   docker-compose logs -f
   
   # 查看特定服务日志
   docker-compose logs -f traefik
   docker-compose logs -f custom-app
   ```

## 🌐 访问服务

启动成功后，您可以通过以下地址访问不同的服务：

### 1. Traefik Dashboard
- **地址**: http://traefik.localhost:8080
- **说明**: Traefik的管理界面，可以查看路由、服务、中间件等配置

### 2. 自定义应用
- **地址**: http://app.localhost
- **API端点**:
  - 主页: http://app.localhost/
  - 健康检查: http://app.localhost/health
  - API信息: http://app.localhost/api/info
  - 用户列表: http://app.localhost/api/users
  - 特定用户: http://app.localhost/api/users/1

### 3. Nginx应用
- **地址**: http://app1.localhost
- **说明**: 标准的Nginx欢迎页面

### 4. Apache应用
- **地址**: http://app2.localhost
- **说明**: 标准的Apache欢迎页面

## 🔧 配置说明

### Docker Compose配置

`docker-compose.yml` 文件定义了以下服务：

1. **traefik**: 反向代理服务
   - 端口映射: 80, 443, 8080
   - 挂载Docker socket用于服务发现
   - 启用Dashboard和API

2. **webapp1**: Nginx示例应用
   - 通过 `app1.localhost` 访问

3. **webapp2**: Apache示例应用
   - 通过 `app2.localhost` 访问

4. **custom-app**: 自定义应用
   - 通过 `app.localhost` 访问
   - 包含完整的REST API示例
   - 支持Go和Python两个版本

### Traefik配置

`traefik.yml` 文件包含：

- **入口点配置**: HTTP(80)和HTTPS(443)端口
- **提供者配置**: Docker和文件提供者
- **日志配置**: JSON格式的访问日志
- **健康检查**: 自动健康检查配置

### 应用特性对比

| 特性 | Go版本 | Python版本 |
|------|--------|------------|
| 框架 | Gin | Flask |
| 性能 | 高 | 中等 |
| 内存占用 | 低 | 中等 |
| 开发速度 | 中等 | 快 |
| 部署大小 | 小 | 中等 |
| 学习曲线 | 中等 | 简单 |

## 🛠️ 自定义配置

### 添加新的服务

1. 在 `docker-compose.yml` 中添加新服务
2. 添加Traefik标签配置路由规则
3. 重启服务: `docker-compose up -d`

### 修改路由规则

通过修改Docker标签来调整路由：

```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.myservice.rule=Host(`myservice.localhost`)"
  - "traefik.http.routers.myservice.entrypoints=web"
  - "traefik.http.services.myservice.loadbalancer.server.port=8080"
```

### 添加中间件

```yaml
labels:
  # 添加自定义头部
  - "traefik.http.middlewares.custom-headers.headers.customrequestheaders.X-Custom-Header=Value"
  # 应用中间件到路由
  - "traefik.http.routers.myservice.middlewares=custom-headers"
```

## 📊 监控和调试

### 查看Traefik日志
```bash
docker-compose logs -f traefik
```

### 查看应用日志
```bash
docker-compose logs -f custom-app
```

### 访问Traefik API
```bash
# 获取所有路由
curl http://traefik.localhost:8080/api/http/routers

# 获取所有服务
curl http://traefik.localhost:8080/api/http/services
```

## 🧪 测试API

### 测试用户API
```bash
# 获取用户列表
curl http://app.localhost/api/users

# 获取特定用户
curl http://app.localhost/api/users/1

# 创建新用户
curl -X POST http://app.localhost/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"测试用户","email":"test@example.com"}'
```

### 测试健康检查
```bash
curl http://app.localhost/health
```

## 🛑 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除所有容器和网络
docker-compose down --volumes --remove-orphans
```

## 🔍 故障排除

### 常见问题

1. **端口被占用**
   - 检查端口80、443、8080是否被其他服务占用
   - 修改 `docker-compose.yml` 中的端口映射

2. **无法访问服务**
   - 确保Docker服务正在运行
   - 检查容器状态: `docker-compose ps`
   - 查看容器日志: `docker-compose logs`

3. **域名解析问题**
   - 在 `/etc/hosts` 文件中添加域名映射
   - 或者使用 `localhost` 替代域名

4. **应用构建失败**
   - **Go版本**: 检查Go版本兼容性，确保网络连接正常以下载依赖
   - **Python版本**: 检查Python版本兼容性，确保依赖安装正确

### 添加hosts映射

在 `/etc/hosts` 文件中添加：

```
127.0.0.1 traefik.localhost
127.0.0.1 app.localhost
127.0.0.1 app1.localhost
127.0.0.1 app2.localhost
```

## 📚 学习资源

- [Traefik官方文档](https://doc.traefik.io/traefik/)
- [Docker Compose文档](https://docs.docker.com/compose/)
- [Go语言官方文档](https://golang.org/doc/)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [Python官方文档](https://docs.python.org/)
- [Flask框架文档](https://flask.palletsprojects.com/)

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个演示项目！

## �� 许可证

MIT License 