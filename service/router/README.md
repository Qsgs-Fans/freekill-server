此为网关服务。需要以下服务：

- Etcd：注册rpc服务
- MySql/Mariadb + Redis：存储IP黑名单

首先动手编写proto文件，定义信息格式和Rpc服务。然后生成代码。

下面的流程只是记录一下，不要在文件夹里面运行

```sh
$ mkdir proto
$ vim proto/router.proto
$ goctl rpc protoc proto/router.proto --go_out=./ --go-grpc_out=./ --zrpc_out=.
```

然后根据需要新增文件、修改生成的文件等。
