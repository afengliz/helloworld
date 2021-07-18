package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"helloworld/internal/biz"
	"helloworld/internal/data/ent"
	"helloworld/internal/data/ent/user"
)
var _ biz.UserRepo = (*userRepo)(nil)
type userRepo struct {
	data *Data
	log  *log.Helper
}

func (u userRepo) GetUser(ctx context.Context, name string) (*ent.User,error) {
	 user,err := u.data.db.User.Query().Where(user.Name(name)).First(ctx)
	 return user,err
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}


