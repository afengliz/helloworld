package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"helloworld/internal/data/ent"
)

type UserRepo interface {
	GetUser(context.Context,string) (*ent.User,error)
}


type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) GetUser(ctx context.Context, name string) (*ent.User,error) {
	return uc.repo.GetUser(ctx, name)
}