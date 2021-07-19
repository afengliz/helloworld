package biz

import (
	"container/heap"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"helloworld/internal/utils"
	"sync"
	"time"
)

type BucketCache interface {
	Add(key string)
	Get(key string) interface{}
	transport.Server
}

func NewBucketCache(userRepo UserRepo, logger log.Logger) BucketCache {
	arr := make([]*bucket,5)
	for i := 0; i < 5; i++ {
		arr[i] = newBucket()
	}
	bChan := make(chan string,10)
	b := &bucketCache{buckets: arr, bChan: bChan,bSize:5,limit:2,dMap: make(map[string]interface{}),userRepo:userRepo,log: log.NewHelper(logger)}
	return b
}

type bucketCache struct {
	// 每一秒的数据统计
	buckets []*bucket
	// 默认统计最近5秒
	bSize int
	// chan 默认大小是10
	bChan chan string
	// 锁
	rwMutex sync.RWMutex
	// data
	dMap map[string]interface{}
	// 
	oldIndex int
	// limit
	limit int
	// userRepo
	userRepo UserRepo
	// log
	log    *log.Helper
}

func (b *bucketCache) Add(key string){
	b.bChan<-key
}
func (b *bucketCache) Get(key string) interface{}{
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()
	if item,ok :=b.dMap[key];ok{
		return item
	}
	return nil
}
func (b *bucketCache) Start(ctx context.Context) error{
	go b.run()
	go b.tickerGenBCache(ctx)
	b.log.Info("启动bucket cache")
	return nil
}
func (b *bucketCache) Stop(context.Context) error{
	close(b.bChan)
	b.log.Info("结束bucket cache")
	return nil
}


func (b *bucketCache) run(){
	for item := range b.bChan {
		b.rwMutex.Lock()
		curIndex := time.Now().Second()% b.bSize
		if b.oldIndex != curIndex {
			b.oldIndex = curIndex
			b.buckets[curIndex].clear()
		}
		b.rwMutex.Unlock()
		b.buckets[curIndex].add(item)
	}
}

// 定时重新生成localcache
func (b *bucketCache) tickerGenBCache(ctx context.Context){
	ticker  := time.NewTicker(30 * time.Second)
	for{
		select {
		case <-ticker.C :
			b.genBCache()
		case <-ctx.Done():
			break
		}
	}
}

func (b *bucketCache) genBCache(){
	//生成cache
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()
	cIndex := time.Now().Second()% b.bSize
	// 不算当前秒，只统计最近四秒的数据
	tMap := make(map[string]int)
	for i := 1; i <= 4; i++ {
		cur := cIndex - i
		if cur < 0 {
			cur = 5 + cur
		}
		cMap := b.buckets[cur]
		for cKey, cVal := range cMap.dataMap {
			if c ,ok :=tMap[cKey];ok{
				tMap[cKey] = c + cVal
			}else{
				tMap[cKey] = cVal
			}
		}
	}
	// 构建大跟堆
	bHeap := utils.InitBigRootHeap()
	heap.Init(bHeap)
	for key, val := range tMap {
		if val < b.limit{
			continue
		}
		heap.Push(bHeap,&utils.HeapNode{Key: key,Count: val})
	}
	keys := make([]string,bHeap.Len())
	for index := 0;  bHeap.Len()>0;index++{
		node := heap.Pop(bHeap).(*utils.HeapNode)
		keys[index] = node.Key
	}
	users,err := b.userRepo.GetUsers(context.Background(),keys)
	if err != nil{
		b.log.Error(err)
	}
	// 清除数据
	for item := range b.dMap {
		delete(b.dMap,item)
	}
	// 插入新数据
	for i := 0; i < len(users); i++ {
		b.dMap[users[i].Name] = users[i]
		b.log.Debug("生成缓存：",users[i].Name)
	}

}


type bucket struct {
	mutex sync.Mutex
	dataMap map[string]int
}

func newBucket()  *bucket {
	return &bucket{dataMap: make(map[string]int)}
}

func (b *bucket) add(key string){
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if item,ok:=b.dataMap[key];ok{
		b.dataMap[key] = item+1
	}else{
		b.dataMap[key] = 1
	}
}
func (b *bucket) clear(){
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for item := range b.dataMap {
		delete(b.dataMap,item)
	}
}