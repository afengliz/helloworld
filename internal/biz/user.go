package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/singleflight"
	"helloworld/internal/data/ent"
)

type UserRepo interface {
	GetUser(context.Context,string) (*ent.User,error)
}


type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
	sGroup *singleflight.Group
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger),sGroup: &singleflight.Group{}}
}

func (uc *UserUsecase) GetUser(ctx context.Context, name string) (*ent.User,error) {
	return uc.repo.GetUser(ctx, name)
}

func (uc *UserUsecase) GetUserBySingleFlight(ctx context.Context, name string) (*ent.User,error) {
	ans,err,_ := uc.sGroup.Do(name, func() (interface{}, error) {
		return uc.repo.GetUser(ctx, name)
	})
	if err != nil{
		return nil,err
	}
	return ans.(*ent.User),nil
}