# FreeKill - 服务器

FreeKill server意在取代FreeKill客户端一般会自带的服务端功能，提供一个更好的后端实现，
同时也兼容现有的客户端实现。

- [ ] 其他服务应通过消息队列服务实现到网关服务的通信，而非在网关定义RPC服务

## 开坑理由

FreeKill项目中已经自带了一个服务器，使用和客户端一样的Qt架构实现，提供了用户登录、
房间创建、数据存储等等服务端该有的功能。然而其本身只适用于单机游戏或者规模很小的私服，
理由有以下几点：

- 单进程架构，无法横向扩容
- 缺少网关层，没有做流量检测与频率限制之类的
- 随着运行时长增加，内存碎片化严重，容易OOM
- 出于单机游戏的考量，数据库选了sqlite，且没有引入缓存
- 以及许多其他问题……

为了优化上述缺点，故决定采用更适合后端的语言重新设计并开发一个后端。
由于游戏逻辑完全Lua化，因此将服务端从头来过的可行性也并非不存在。

此外为了节约精力，用go-zero库直接做微服务发现和数据库CRUD等。。

## 服务划分

考虑到C++版服务端实现的功能，可能可以这样划分：

- 网关服务：用于维持TCP长连接，将其他服务（主要是游戏逻辑服务）传来的消息合成json发送给客户端；将客户端发来的json数据解析并发送给其他服务。还需要对客户端的流量进行控制和过滤等。
- 用户登录服务：维持着已登录用户列表
- 游戏逻辑服务：维持着Lua VM
- 管理员服务：维持着类似目前shell的功能

## 编译运行

依赖：

安装流程（Linux下）：

Docker应该是个好办法，但目前还是先写个todo罢，下面介绍在本机编译运行的流程。

### 0. 安装好Go和依赖 这些依赖也包括项目之外的服务

```sh
# (均为Debian，记得在普通用户下)
$ sudo apt install golang
# 设置go代理 否则难以下载依赖
$ go env -w GOPROXY=https://goproxy.cn,direct
# 设置PATH变量 最好写进shell配置文件（~/.bashrc之类的）
$ export PATH=$PATH:$(go env GOPATH)/bin

# goctl是编译所需的依赖，用来生成代码
$ go install github.com/zeromicro/go-zero/tools/goctl@latest
# 注：goctl应该安装了protoc之类的protobuf生成器，此命令检查是否漏了
$ goctl env check --install --verbose --force

$ git clone https://github.com/Qsgs-Fans/freekill-server.git
$ cd freekill-server
$ go mod download
```

### 1. 启动除了freekill-server本身之外的其余服务

```sh
# 窗口1之内启动etcd(默认监听2379端口)
$ etcd
# 窗口2之内启动redis
$ redis-server
```

### 2. 启动freekill-server

下面介绍的都是直接启动的方式，你可以通过其他方法简化启动流程。

```sh
# 窗口1
$ cd service/user
$ go run user.go

# 窗口2
$ cd service/router
$ go run router.go
```

可以用FreeKill客户端测试。毕竟目标就是让正常客户端能正常连接。

## 部署到服务器
