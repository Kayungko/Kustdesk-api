# KustDesk API

[English Doc](README_EN.md)


## 🚀 主要改进

- ✅ **增强的中文本地化支持** - 完善的中文界面和提示信息
- ✅ **优化的用户管理** - 支持账户有效期和设备数量限制
- ✅ **改进的系统配置** - 更直观的管理界面
- ✅ **Docker 容器化部署** - 简单快速部署
- ✅ **修复了多个Bug** - 提升系统稳定性
本项目使用 Go 实现了 RustDesk 的 API，并包含了 Web Admin 和 Web 客户端。


<div align=center>
<img src="https://img.shields.io/badge/golang-1.22-blue"/>
<img src="https://img.shields.io/badge/gin-v1.9.0-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-v1.25.7-green"/>
<img src="https://img.shields.io/badge/swag-v1.16.3-yellow"/>
<img src="https://goreportcard.com/badge/github.com/lejianwen/rustdesk-api/v2"/>
<img src="https://github.com/lejianwen/rustdesk-api/actions/workflows/build.yml/badge.svg"/>
</div>

## 搭配[lejianwen/rustdesk-server]使用更佳。
> [lejianwen/rustdesk-server]fork自RustDesk Server官方仓库
> 1. 解决了使用API链接超时问题
> 2. 可以强制登录后才能发起链接
> 3. 支持客户端websocket



# 特性

- PC端API
    - 个人版API
    - 登录
    - 地址簿
    - 群组
    - 授权登录
      - 支持`github`, `google` 和 `OIDC` 登录，
      - 支持`web后台`授权登录
      - 支持`LDAP`(AD和OpenLDAP已测试), 如果API Server配置了LDAP
    - i18n
- Web Admin
    - 用户管理
    - 设备管理
    - 地址簿管理
    - 标签管理
    - 群组管理
    - Oauth 管理
    - 配置LDAP, 配置文件或者环境变量
    - 登录日志
    - 链接日志
    - 文件传输日志
    - 快速使用web client
    - i18n
    - 通过 web client 分享给游客
    - server控制(一些官方的简单的指令 [WIKI](https://github.com/lejianwen/rustdesk-api/wiki/Rustdesk-Command))
- Web Client
    - 自动获取API server
    - 自动获取ID服务器和KEY
    - 自动获取地址簿
    - 游客通过临时分享链接直接远程到设备
    - v2 Preview
- CLI
    - 重置管理员密码

## 功能

## 🐳 Docker 快速部署

### 使用预构建镜像（推荐）
```bash
# 拉取镜像
docker pull kayung1012/kustdesk-api:latest

# 使用 docker-compose 部署
git clone https://github.com/Kayungko/kustdesk-server.git
cd kustdesk-server/Kustdesk-api
docker-compose up -d
```

### 本地构建
```bash
# 克隆项目
git clone https://github.com/Kayungko/kustdesk-server.git
cd kustdesk-server/Kustdesk-api

# 构建镜像
docker build -f Dockerfile.simple -t kustdesk-api:latest .
```

详细部署说明请参考 [README_DOCKER.md](README_DOCKER.md)


### API 服务 
基本实现了PC端基础的接口。支持Personal版本接口，可以通过配置文件`rustdesk.personal`或环境变量`RUSTDESK_API_RUSTDESK_PERSONAL`来控制是否启用

<table>
    <tr>
      <td width="50%" align="center" colspan="2"><b>登录</b></td>
    </tr>
    <tr>
        <td width="50%" align="center" colspan="2"><img src="docs/pc_login.png"></td>
    </tr>
     <tr>
      <td width="50%" align="center"><b>地址簿</b></td>
      <td width="50%" align="center"><b>群组</b></td>
    </tr>
    <tr>
        <td width="50%" align="center"><img src="docs/pc_ab.png"></td>
        <td width="50%" align="center"><img src="docs/pc_gr.png"></td>
    </tr>
</table>

### Web Admin:

* 使用前后端分离，提供用户友好的管理界面，主要用来管理和展示。前端代码在[rustdesk-api-web](https://github.com/lejianwen/rustdesk-api-web)

* 后台访问地址是`http://<your server>[:port]/_admin/`
* 初次安装管理员为用户名为`admin`，密码将在控制台打印，可以通过[命令行](#CLI)更改密码

  ![img.png](./docs/init_admin_pwd.png)

1. 管理员界面
   ![web_admin](docs/web_admin.png)
2. 普通用户界面
   ![web_user](docs/web_admin_user.png)

3. 每个用户可以多个地址簿，也可以将地址簿共享给其他用户
4. 分组可以自定义，方便管理，暂时支持两种类型: `共享组` 和 `普通组`
5. 可以直接打开webclient，方便使用；也可以分享给游客，游客可以直接通过webclient远程到设备
6. Oauth,支持了`Github`, `Google` 以及 `OIDC`, 需要创建一个`OAuth App`，然后配置到后台
    - 对于`Google` 和 `Github`, `Issuer` 和 `Scopes`不需要填写.
    - 对于`OIDC`, `Issuer`是必须的。`Scopes`是可选的，默认为 `openid,profile,email`. 确保可以获取 `sub`,`email` 和`preferred_username`
    - `github oauth app`在`Settings`->`Developer settings`->`OAuth Apps`->`New OAuth App`
      中创建,地址 [https://github.com/settings/developers](https://github.com/settings/developers)
    - `Authorization callback URL`填写`http://<your server[:port]>/api/oidc/callback`
      ，比如`http://127.0.0.1:21114/api/oidc/callback`
7. 登录日志
8. 链接日志
9. 文件传输日志
10. server控制

  - `简易模式`,已经界面化了一些简单的指令，可以直接在后台执行
    ![rustdesk_command_simple](./docs/rustdesk_command_simple.png)

  - `高级模式`,直接在后台执行指令
      * 可以官方指令
      * 可以添加自定义指令
      * 可以执行自定义指令

 
11. **LDAP 支持**, 当在API Server上设置了LDAP(已测试AD和LDAP),可以通过LDAP中的用户信息进行登录 https://github.com/lejianwen/rustdesk-api/issues/114 ,如果LDAP验证失败，返回本地用户

### Web Client:

1. 如果已经登录了后台，web client将自动直接登录
2. 如果没登录后台，点击右上角登录即可，api server已经自动配置好了
3. 登录后，会自动同步ID服务器和KEY
4. 登录后，会将地址簿自动保存到web client中，方便使用
5. 现已支持`v2 Preview`，访问路径是`/webclient2`
   ![webclientv2](./docs/webclientv2.png)
6. `v2 preview` 部署，参考[WIKI](https://github.com/lejianwen/rustdesk-api/wiki)


### 自动化文档: 使用 Swag 生成 API 文档，方便开发者理解和使用 API。

1. 后台文档 `<youer server[:port]>/admin/swagger/index.html`
2. PC端文档 `<youer server[:port]>/swagger/index.html`
   ![api_swag](docs/api_swag.png)

### CLI

```bash
# 查看帮助
./apimain -h
```

#### 重置管理员密码
```bash
./apimain reset-admin-pwd <pwd>
```

## 安装与运行

### 相关配置

* [配置文件](./conf/config.yaml)
* 参考`conf/config.yaml`配置文件，修改相关配置。
* 如果`gorm.type`是`sqlite`，则不需要配置mysql相关配置。
* 语言如果不设置默认为`zh-CN`

### 环境变量
环境变量和配置文件`conf/config.yaml`中的配置一一对应，变量名前缀是`RUSTDESK_API`
下面表格并未全部列出，可以参考`conf/config.yaml`中的配置。

| 变量名                                                    | 说明                                                                             | 示例                           |
|--------------------------------------------------------|--------------------------------------------------------------------------------|------------------------------|
| TZ                                                     | 时区                                                                             | Asia/Shanghai                |
| RUSTDESK_API_LANG                                      | 语言                                                                             | `en`,`zh-CN`                 |
| RUSTDESK_API_APP_WEB_CLIENT                            | 是否启用web-client; 1:启用,0:不启用; 默认启用                                               | 1                            |
| RUSTDESK_API_APP_REGISTER                              | 是否开启注册; `true`, `false`  默认`false`                                             | `false`                      |
| RUSTDESK_API_APP_SHOW_SWAGGER                          | 是否可见swagger文档;`1`显示，`0`不显示，默认`0`不显示                                            | `1`                          |
| RUSTDESK_API_APP_TOKEN_EXPIRE                          | token有效时长                                                                      | `168h`                       |
| RUSTDESK_API_APP_DISABLE_PWD_LOGIN                     | 是否禁用密码登录;  `true`, `false`  默认`false`                                          | `false`                      |
| RUSTDESK_API_APP_REGISTER_STATUS                       | 注册用户默认状态; 1 启用，2 禁用, 默认 1                                                      | `1`                          |
| RUSTDESK_API_APP_CAPTCHA_THRESHOLD                     | 验证码触发次数; -1 不启用， 0 一直启用， >0 登录错误次数后启用 ;默认 `3`                                  | `3`                          |
| RUSTDESK_API_APP_BAN_THRESHOLD                         | 封禁IP触发次数; 0 不启用, >0 登录错误次数后封禁IP; 默认 `0`                                        | `0`                          |
| -----ADMIN配置-----                                      | ----------                                                                     | ----------                   |
| RUSTDESK_API_ADMIN_TITLE                               | 后台标题                                                                           | `RustDesk Api Admin`         |
| RUSTDESK_API_ADMIN_HELLO                               | 后台欢迎语，可以使用`html`                                                               |                              |
| RUSTDESK_API_ADMIN_HELLO_FILE                          | 后台欢迎语文件，如果内容多，使用文件更方便。<br>会覆盖`RUSTDESK_API_ADMIN_HELLO`                        | `./conf/admin/hello.html`    |
| -----GIN配置-----                                        | ----------                                                                     | ----------                   |
| RUSTDESK_API_GIN_TRUST_PROXY                           | 信任的代理IP列表，以`,`分割，默认信任所有                                                        | 192.168.1.2,192.168.1.3      |
| -----GORM配置-----                                       | ----------                                                                     | ---------------------------  |
| RUSTDESK_API_GORM_TYPE                                 | 数据库类型sqlite或者mysql，默认sqlite                                                    | sqlite                       |
| RUSTDESK_API_GORM_MAX_IDLE_CONNS                       | 数据库最大空闲连接数                                                                     | 10                           |
| RUSTDESK_API_GORM_MAX_OPEN_CONNS                       | 数据库最大打开连接数                                                                     | 100                          |
| RUSTDESK_API_RUSTDESK_PERSONAL                         | 是否启用个人版API， 1:启用,0:不启用； 默认启用                                                   | 1                            |
| -----MYSQL配置-----                                      | ----------                                                                     | ----------                   |
| RUSTDESK_API_MYSQL_USERNAME                            | mysql用户名                                                                       | root                         |
| RUSTDESK_API_MYSQL_PASSWORD                            | mysql密码                                                                        | 111111                       |
| RUSTDESK_API_MYSQL_ADDR                                | mysql地址                                                                        | 192.168.1.66:3306            |
| RUSTDESK_API_MYSQL_DBNAME                              | mysql数据库名                                                                      | rustdesk                     |
| RUSTDESK_API_MYSQL_TLS                             | 是否启用TLS, 可选值: `true`, `false`, `skip-verify`, `custom` | `false`                      |
| -----RUSTDESK配置-----                                   | ----------                                                                     | ----------                   |
| RUSTDESK_API_RUSTDESK_ID_SERVER                        | Rustdesk的id服务器地址                                                               | 192.168.1.66:21116           |
| RUSTDESK_API_RUSTDESK_RELAY_SERVER                     | Rustdesk的relay服务器地址                                                            | 192.168.1.66:21117           |
| RUSTDESK_API_RUSTDESK_API_SERVER                       | Rustdesk的api服务器地址                                                              | http://192.168.1.66:21114    |
| RUSTDESK_API_RUSTDESK_KEY                              | Rustdesk的key                                                                   | 123456789                    |
| RUSTDESK_API_RUSTDESK_KEY_FILE                         | Rustdesk存放key的文件                                                               | `./conf/data/id_ed25519.pub` |
| RUSTDESK_API_RUSTDESK_WEBCLIENT<br/>_MAGIC_QUERYONLINE | Web client v2 中是否启用新的在线状态查询方法; `1`:启用,`0`:不启用,默认不启用                            | `0`                          |
| RUSTDESK_API_RUSTDESK_WS_HOST                          | 自定义Websocket Host                                                              | `wss://192.168.1.123:1234`   |
| ----PROXY配置-----                                       | ----------                                                                     | ----------                   |
| RUSTDESK_API_PROXY_ENABLE                              | 是否启用代理:`false`, `true`                                                         | `false`                      |
| RUSTDESK_API_PROXY_HOST                                | 代理地址                                                                           | `http://127.0.0.1:1080`      |
| ----JWT配置----                                          | --------                                                                       | --------                     |
| RUSTDESK_API_JWT_KEY                                   | 自定义JWT KEY,为空则不启用JWT<br/>如果没使用`lejianwen/rustdesk-server`中的`MUST_LOGIN`，建议设置为空 |                              |
| RUSTDESK_API_JWT_EXPIRE_DURATION                       | JWT有效时间                                                                        | `168h`                       |


### 运行

#### docker运行

1. 直接docker运行,配置可以通过挂载配置文件`/app/conf/config.yaml`来修改,或者通过环境变量覆盖配置文件中的配置

    ```bash
    docker run -d --name rustdesk-api -p 21114:21114 \
    -v /data/rustdesk/api:/app/data \
    -e TZ=Asia/Shanghai \
    -e RUSTDESK_API_LANG=zh-CN \
    -e RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.66:21116 \
    -e RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.66:21117 \
    -e RUSTDESK_API_RUSTDESK_API_SERVER=http://192.168.1.66:21114 \
    -e RUSTDESK_API_RUSTDESK_KEY=<key> \
    lejianwen/rustdesk-api
    ```

2. 使用`docker compose`，参考[WIKI](https://github.com/lejianwen/rustdesk-api/wiki)

#### 下载release直接运行

[下载地址](https://github.com/lejianwen/rustdesk-api/releases)

#### 源码安装

1. 克隆仓库
   ```bash
   git clone https://github.com/lejianwen/rustdesk-api.git
   cd rustdesk-api
   ```

2. 安装依赖

    ```bash
    go mod tidy
    #安装swag，如果不需要生成文档，可以不安装
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

3. 编译后台前端，前端代码在[rustdesk-api-web](https://github.com/lejianwen/rustdesk-api-web)中
   ```bash
   cd resources
   mkdir -p admin
   git clone https://github.com/lejianwen/rustdesk-api-web
   cd rustdesk-api-web
   npm install
   npm run build
   cp -ar dist/* ../admin/
   ```
4. 运行
    ```bash
    #直接运行
    go run cmd/apimain.go
    #或者使用generate_api.go生成api并运行
    go generate generate_api.go
    ```
   > 注意：使用 `go run` 或编译后的二进制时，当前目录下必须存在 `conf` 和 `resources`
   > 目录。如果在其他目录运行，可通过 `-c` 和环境变量
   > `RUSTDESK_API_GIN_RESOURCES_PATH` 指定绝对路径，例如：
   > ```bash
   > RUSTDESK_API_GIN_RESOURCES_PATH=/opt/rustdesk-api/resources ./apimain -c /opt/rustdesk-api/conf/config.yaml
   > ```
5. 编译，如果想自己编译,先cd到项目根目录，然后windows下直接运行`build.bat`,linux下运行`build.sh`,编译后会在`release`
   目录下生成对应的可执行文件。直接运行编译后的可执行文件即可。

6. 打开浏览器访问`http://<your server[:port]>/_admin/`，默认用户名密码为`admin`，请及时更改密码。


#### 使用`lejianwen/server-s6`镜像运行

- 已解决链接超时问题
- 可以强制登录后才能发起链接
- github https://github.com/lejianwen/rustdesk-server

```yaml
 networks:
   rustdesk-net:
     external: false
 services:
   rustdesk:
     ports:
       - 21114:21114
       - 21115:21115
       - 21116:21116
       - 21116:21116/udp
       - 21117:21117
       - 21118:21118
       - 21119:21119
     image: lejianwen/rustdesk-server-s6:latest
     environment:
       - RELAY=<relay_server[:port]>
       - ENCRYPTED_ONLY=1
       - MUST_LOGIN=N
       - TZ=Asia/Shanghai
       - RUSTDESK_API_RUSTDESK_ID_SERVER=<id_server[:21116]>
       - RUSTDESK_API_RUSTDESK_RELAY_SERVER=<relay_server[:21117]>
       - RUSTDESK_API_RUSTDESK_API_SERVER=http://<api_server[:21114]>
       - RUSTDESK_API_KEY_FILE=/data/id_ed25519.pub
       - RUSTDESK_API_JWT_KEY=xxxxxx # jwt key
     volumes:
       - /data/rustdesk/server:/data
       - /data/rustdesk/api:/app/data #将数据库挂载
     networks:
       - rustdesk-net
     restart: unless-stopped
       
```


## 其他

- [WIKI](https://github.com/lejianwen/rustdesk-api/wiki)
- [链接超时问题](https://github.com/lejianwen/rustdesk-api/issues/92)
- [修改客户端ID](https://github.com/abdullah-erturk/RustDesk-ID-Changer)
- [webclient来源](https://hub.docker.com/r/keyurbhole/flutter_web_desk)


## 鸣谢

感谢所有做过贡献的人!

<a href="https://github.com/lejianwen/rustdesk-api/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=lejianwen/rustdesk-api" />
</a>

## 感谢你的支持！如果这个项目对你有帮助，请点个⭐️鼓励一下，谢谢！

[lejianwen/rustdesk-server]: https://github.com/lejianwen/rustdesk-server