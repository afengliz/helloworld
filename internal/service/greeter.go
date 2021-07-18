package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	v1 "helloworld/api/helloworld/v1"
	"helloworld/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer
	gu *biz.GreeterUsecase
	uu *biz.UserUsecase
	log *log.Helper
	requestGroup *singleflight.Group
}

// NewGreeterService new a greeter service.
func NewGreeterService(gu *biz.GreeterUsecase,uu *biz.UserUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{gu: gu,uu:uu, log: log.NewHelper(logger),requestGroup: &singleflight.Group{}}
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("SayHello Received: %v", in.GetName())

	if in.GetName() == "error" {
		err := v1.ErrorUserNotFound("user not found: %s", in.GetName())
		merr := errors.Wrap(err,"哈嘿")
		return nil, merr
	}
	return &v1.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func  (s *GreeterService) GetUserInfo(ctx context.Context,in *v1.GetUserRequest) (*v1.GetUserReply, error) {
	ans,err := s.uu.GetUser(ctx,in.GetName())
	if err != nil{
		return nil,err
	}
	return &v1.GetUserReply{Name: ans.Name,Age: int32(ans.Age)},nil
}
