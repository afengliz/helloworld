package test

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"helloworld/internal/biz"
	"helloworld/internal/conf"
	"helloworld/internal/data"
	"os"
	"sync"
	"testing"
	"time"
)

// singleflight 测试
func TestUserUsecase_GetUserBySingleFlight(t *testing.T) {
	var (
		wg1  sync.WaitGroup
		wg2  sync.WaitGroup
		now1 = time.Now()
		n    = 1000
	)
	logger := log.With(log.NewStdLogger(os.Stdout))
	dataData, _, err := data.NewData(&conf.Data{Database:&conf.Data_Database{
		Driver: "mysql",
		Source: "root:root@tcp(127.0.0.1:3306)/demo",
	} }, logger)
	if err != nil {
		t.Fail()
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(userRepo, logger)
	for i := 0; i < n; i++ {
		wg1.Add(1)
		go func() {
			res, err := userUsecase.GetUserBySingleFlight(context.Background(),"liyanfeng")
			if err != nil || res == nil {
				t.Fail()
			}
			wg1.Done()
		}()
	}

	wg1.Wait()
	now2 := time.Now()
	fmt.Printf("同时发起 %d 次请求，singleflight耗时: %s\n", n, time.Since(now1))
	for i := 0; i < n; i++ {
		wg2.Add(1)
		go func() {
			res, err := userUsecase.GetUser(context.Background(),"liyanfeng")
			if err != nil || res == nil {
				t.Fail()
			}
			wg2.Done()
		}()
	}
	wg2.Wait()
	fmt.Printf("同时发起 %d 次请求，不使用singleflight耗时: %s\n", n, time.Since(now2))
}
