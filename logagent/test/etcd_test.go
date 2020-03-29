package test

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"testing"
	"time"
)

func TestPut(t *testing.T){
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil{
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	value := `[{"path":"D:/tmp/nginx.log","topic":"web_log"},{"path":"D:/tmp/redis.log","topic":"redis_log"}]`
	//value := `[{"path":"D:/tmp/nginx.log","topic":"web_log"},{"path":"D:/tmp/mysql.log","topic":"mysql_log"},{"path":"D:/tmp/redis.log","topic":"redis_log"}]`
	_,err = cli.Put(ctx, "/logagent/192.168.0.102/collect_config", value)
	cancel()
	if err != nil{
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
}
