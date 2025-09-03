# 🐳 RustDesk API Server - Docker 部署指南

## 📋 项目简介

这是一个基于 [lejianwen/rustdesk-api](https://github.com/lejianwen/rustdesk-api) 的定制版本，主要特点：

- ✅ **增强的中文本地化支持** - 完善的中文界面和提示信息
- ✅ **优化的用户管理** - 支持账户有效期和设备数量限制
- ✅ **改进的系统配置** - 更直观的管理界面
- ✅ **Docker 容器化部署** - 简单快速部署

## 🚀 快速开始

### 方式1：使用 Docker Hub 镜像（推荐）

```bash
# 拉取最新镜像
docker pull kayung1012/kustdesk-api:latest

# 或者拉取特定版本
docker pull kayung1012/kustdesk-api:v1.0.0
```

### 方式2：本地构建镜像

```bash
# 克隆项目
git clone https://github.com/Kayungko/kustdesk-server.git
cd kustdesk-server/Kustdesk-api

# 构建镜像
docker build -f Dockerfile.simple -t kayung1012/kustdesk-api:latest .
```

## 📁 项目结构

```
kustdesk-server/
├── Kustdesk-api/          # 后端 API 服务
│   ├── Dockerfile.simple  # 简化版 Dockerfile
│   ├── docker-compose-dev.yaml  # 开发环境配置
│   └── ...
└── Kustdesk-api-web/      # 前端管理界面
    ├── src/views/user/     # 用户管理页面
    ├── src/views/system/   # 系统配置页面
    └── src/utils/i18n/     # 国际化文件
```

## ⚠️ 重要配置说明

### RustDesk 服务端关键配置

- **`MUST_LOGIN`**: 
  - `N` (默认): 无需登录即可连接，适合内网环境
  - `Y`: 必须登录才能连接，适合公网环境，更安全

- **`ENCRYPTED_ONLY`**: 
  - `1` (默认): 仅允许加密连接，更安全
  - `0`: 允许非加密连接，兼容性更好但安全性较低

- **`RELAY`**: 中继服务器地址，用于NAT穿透

## 🐳 Docker Compose 部署

### 方案1：仅API服务部署（推荐新手）

如果你已经有RustDesk服务端，只需要API服务：

```yaml
version: '3.8'

services:
  rustdesk-api:
    image: kayung1012/kustdesk-api:latest
    container_name: rustdesk-api
    environment:
      - TZ=Asia/Shanghai
      - RUSTDESK_API_LANG=zh-CN
      - RUSTDESK_API_RUSTDESK_ID_SERVER=你的ID服务器IP:21116
      - RUSTDESK_API_RUSTDESK_RELAY_SERVER=你的中继服务器IP:21117
      - RUSTDESK_API_RUSTDESK_API_SERVER=http://你的API服务器IP:21114
      - RUSTDESK_API_RUSTDESK_KEY=your_rustdesk_key_here
      - RUSTDESK_API_APP_WEB_CLIENT=1
      - RUSTDESK_API_APP_REGISTER=false
      - RUSTDESK_API_APP_SHOW_SWAGGER=0
      - RUSTDESK_API_APP_MAX_CONCURRENT_DEVICES=3
    ports:
      - "21114:21114"
    volumes:
      - ./data/rustdesk/api:/app/data
      - ./conf:/app/conf
      - ./logs:/app/runtime
    restart: unless-stopped
    networks:
      - rustdesk-net

networks:
  rustdesk-net:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### 方案2：完整RustDesk服务端 + API一体化部署

如果你需要完整的RustDesk服务，包含ID服务器、中继服务器和API：

```yaml
version: '3.8'

networks:
  rustdesk-net:
    external: false

services:
  # RustDesk 服务端 (ID + Relay)
  rustdesk:
    image: lejianwen/rustdesk-server-s6:latest
    container_name: rustdesk-server
    ports:
      - "21115:21115"  # hbbs
      - "21116:21116"  # hbbs
      - "21116:21116/udp"  # hbbs
      - "21117:21117"  # hbbr
      - "21118:21118"  # hbbr
      - "21119:21119"  # hbbr
    environment:
      - RELAY=你的服务器IP:21117
      - ENCRYPTED_ONLY=1
      - MUST_LOGIN=Y  # 默认为 N，设置为 Y 则必须登录才能连接
      - TZ=Asia/Shanghai
      - RUSTDESK_API_RUSTDESK_ID_SERVER=你的服务器IP:21116
      - RUSTDESK_API_RUSTDESK_RELAY_SERVER=你的服务器IP:21117
      - RUSTDESK_API_RUSTDESK_API_SERVER=http://你的服务器IP:21114
      - RUSTDESK_API_KEY_FILE=/data/id_ed25519.pub
      - RUSTDESK_API_JWT_KEY=your_jwt_secret_here
    volumes:
      - ./data/rustdesk/server:/data
      - ./data/rustdesk/api:/app/data
    networks:
      - rustdesk-net
    restart: unless-stopped

  # RustDesk API 服务（使用你的镜像）
  rustdesk-api:
    image: kayung1012/kustdesk-api:latest
    container_name: rustdesk-api
    ports:
      - "21114:21114"
    environment:
      - TZ=Asia/Shanghai
      - RUSTDESK_API_LANG=zh-CN
      - RUSTDESK_API_RUSTDESK_ID_SERVER=你的服务器IP:21116
      - RUSTDESK_API_RUSTDESK_RELAY_SERVER=你的服务器IP:21117
      - RUSTDESK_API_RUSTDESK_API_SERVER=http://你的服务器IP:21114
      - RUSTDESK_API_RUSTDESK_KEY=your_rustdesk_key_here
      - RUSTDESK_API_APP_WEB_CLIENT=1
      - RUSTDESK_API_APP_REGISTER=false
      - RUSTDESK_API_APP_SHOW_SWAGGER=0
      - RUSTDESK_API_APP_MAX_CONCURRENT_DEVICES=3
    volumes:
      - ./data/rustdesk/api:/app/data
      - ./conf:/app/conf
      - ./logs:/app/runtime
    networks:
      - rustdesk-net
    restart: unless-stopped
    depends_on:
      - rustdesk
```

## 🚀 启动服务

```bash
# 启动服务
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f rustdesk-api

# 停止服务
docker-compose down
```

## ⚙️ 环境变量配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `RUSTDESK_API_LANG` | 界面语言 | `zh-CN` | `zh-CN`, `en` |
| `RUSTDESK_API_RUSTDESK_ID_SERVER` | ID服务器地址 | - | `192.168.1.66:21116` |
| `RUSTDESK_API_RUSTDESK_RELAY_SERVER` | 中继服务器地址 | - | `192.168.1.66:21117` |
| `RUSTDESK_API_APP_MAX_CONCURRENT_DEVICES` | 最大并发设备数 | `3` | `5` |
| `RUSTDESK_API_APP_REGISTER` | 是否允许用户注册 | `false` | `true` |
| `MUST_LOGIN` | 是否必须登录才能连接 | `N` | `Y` (必须登录), `N` (无需登录) |
| `ENCRYPTED_ONLY` | 是否仅允许加密连接 | `1` | `1` (仅加密), `0` (允许非加密) |

## 🌐 访问地址

- **管理后台**: http://你的服务器IP:21114/_admin/
- **API 文档**: http://你的服务器IP:21114/swagger/index.html
- **健康检查**: http://你的服务器IP:21114/health

## 🔧 自定义构建

### 使用 Dockerfile.simple

```bash
# 构建镜像
docker build -f Dockerfile.simple -t my-rustdesk-api:latest .

# 运行容器
docker run -d \
  --name rustdesk-api \
  -p 21114:21114 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/conf:/app/conf \
  my-rustdesk-api:latest
```

### 使用 Dockerfile.dev

```bash
# 构建开发版本（包含前端源码）
docker build -f Dockerfile.dev -t my-rustdesk-api:dev .

# 运行开发版本
docker run -d \
  --name rustdesk-api-dev \
  -p 21114:21114 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/conf:/app/conf \
  my-rustdesk-api:dev
```

## 📊 数据持久化

```bash
# 创建数据目录
mkdir -p data/rustdesk/api
mkdir -p conf
mkdir -p logs

# 权限设置
chmod 755 data/rustdesk/api
chmod 755 conf
chmod 755 logs
```

## 🔍 故障排除

### 1. 容器启动失败

```bash
# 查看容器日志
docker logs rustdesk-api

# 检查端口占用
netstat -tlnp | grep 21114

# 检查配置文件
docker exec -it rustdesk-api cat /app/conf/config.yaml
```

### 2. 前端显示异常

```bash
# 重新构建前端
cd Kustdesk-api-web
npm run build

# 更新后端资源
cd ../Kustdesk-api
rm -rf resources/admin
mkdir -p resources/admin
cp -r ../Kustdesk-api-web/dist/* resources/admin/

# 重新构建镜像
docker build -f Dockerfile.simple -t kayung1012/kustdesk-api:latest .
```

### 3. 数据库连接问题

```bash
# 检查数据库文件
ls -la data/rustdesk/api/

# 重置数据库（谨慎操作）
rm -f data/rustdesk/api/*.db
```

## 📝 更新日志

### v1.0.0 (2024-09-03)
- ✨ 完善中文本地化支持
- 🐛 修复设备管理JavaScript错误
- 🔧 优化Docker构建流程
- 📱 改进用户管理界面
- ⚙️ 增强系统配置功能

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

## 🔗 相关链接

- [原版项目](https://github.com/lejianwen/rustdesk-api)
- [Docker Hub 镜像](https://hub.docker.com/r/kayung1012/kustdesk-api)
- [问题反馈](https://github.com/Kayungko/kustdesk-server/issues)

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- GitHub Issues: [kustdesk-server](https://github.com/Kayungko/kustdesk-server/issues)
- Docker Hub: [kayung1012](https://hub.docker.com/u/kayung1012)

---

**⭐ 如果这个项目对你有帮助，请给它一个星标！**
