package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/singleflight"
	"helloworld/internal/data/ent"
	"helloworld/internal/utils"
	"strings"
)

type UserRepo interface {
	GetUser(context.Context, string) (*ent.User, error)
	GetUsers(context.Context, []string) ([]*ent.User, error)
}

type UserUsecase struct {
	repo   UserRepo
	log    *log.Helper
	sGroup *singleflight.Group
	lru    utils.LRU
	bc     BucketCache
}

func NewUserUsecase(repo UserRepo, logger log.Logger,bc BucketCache) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger), sGroup: &singleflight.Group{},lru: utils.NewLRU(3),bc:bc}
}

func (uc *UserUsecase) GetUser(ctx context.Context, name string) (*ent.User, error) {
	name = strings.TrimSpace(name)

	// 从lru缓存获取数据
	cUser := uc.lru.Get(name)
	if cUser != nil{
		return cUser.(*ent.User),nil
	}
	ans,err := uc.repo.GetUser(ctx, name)
	if err != nil{
		return nil,err
	}
	if ans != nil{
		// 设置lru缓存
		uc.lru.Put(name,ans)
	}
	return ans,nil
}

// single flight 写法
func (uc *UserUsecase) GetUserBySingleFlight(ctx context.Context, name string) (*ent.User, error) {
	ans, err, _ := uc.sGroup.Do(name, func() (interface{}, error) {
		return uc.repo.GetUser(ctx, name)
	})
	if err != nil {
		return nil, err
	}
	return ans.(*ent.User), nil
}

// 动态生成热点数据 缓存到本地 的做法
func (uc *UserUsecase) GetUserByBucketCache(ctx context.Context, name string) (*ent.User, error) {
	name = strings.TrimSpace(name)

	cUser := uc.bc.Get(name)
	if cUser != nil{
		uc.bc.Add(name)
		return cUser.(*ent.User),nil
	}
	ans,err := uc.repo.GetUser(ctx, name)
	if err != nil{
		return nil,err
	}
	if ans != nil{
		uc.bc.Add(name)
	}
	return ans,nil
}