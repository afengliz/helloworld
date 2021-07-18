package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "helloworld/api/helloworld/v1"
	"log"
	"time"
)

func main(){
	// kratos客户端连接服务端代码示例
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:9090"), grpc.WithTimeout(time.Second*60))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewGreeterClient(conn)
	res, err := c.GetUserInfo(context.Background(), &v1.GetUserRequest{Name: "liyanfeng"})
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%+v\n", res)
	res, err = c.GetUserInfo(context.Background(), &v1.GetUserRequest{Name: "dingxinxin"})
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%+v\n", res)
}
