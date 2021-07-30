# Kratos Project Template

## 燕锋的笔记

### 1、Import "google/api/annotations.proto" was not found or had errors.

```
出现问题：Import "google/api/annotations.proto" was not found or had errors.
解决方式：https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-win64.zip 下载解压后，把里面的include/google/的protobuf文件夹拷贝到kratos项目下的third_party/google文件夹下
```

### 2、生成pb文件命令

```
kratos proto client api/helloworld/v1/greeter.proto
kratos proto client api/helloworld/v1/greeter.proto --go-http_opt=omitempty=false
```

### 3、post请求j记得加  body: "*"

```
  rpc GetUserInfo (GetUserRequest) returns (GetUserReply){
    option (google.api.http) = {
        post: "/v1/greeter/getuserinfo",
        body: "*",
    };
  }
```
### 4、goland 运行 main.go
```
需要将Program argument 选项配置为-conf ./configs
```


## Install Kratos

```
go get -u github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API swagger json files by proto file
make swagger
# Generate API validator files by proto file
make validate
# Generate all files
make all
```
## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

