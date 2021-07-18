package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"helloworld/internal/conf"
	"helloworld/internal/data/ent"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo,NewUserRepo)

// Data .
type Data struct {
	db *ent.Client // 数据库操作对象
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	mylog := log.NewHelper(logger)

	client, err := ent.Open(
		c.Database.Driver,
		c.Database.Source,
	)
	if err != nil {
		mylog.Errorf("failed opening connection to mysql: %v", err)
		return nil, nil, err
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		mylog.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}
	d := &Data{
		db: client,
	}
	cleanup := func() {
		logger.Log(log.LevelInfo, "closing the data resources")
	}
	return d, cleanup, nil
}
