// logagent 入口程序
package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"logagent/utils"
	"sync"
	"time"
)

var cfg = new(conf.AppConf)

func main() {
	// 0. 加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}
	// 1. 初始化 kafka 连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka failed, err:%v\n", err)
		return
	}
	fmt.Println("init kafka success.")
	// 2. 初始化 etcd
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("Init taillog failed, err:%v\n", err)
		return
	}
	// 为了实现每个 logagent 都拉取自己独有的配置，所以要以自己的IP地址作为区分
	ipStr, err := utils.GetOutboundIP()
	if err != nil{
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ipStr)
	fmt.Println("etcdConfKey",etcdConfKey)
	// 2.1 从etcd中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf("etcd.GetConf err:%v", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)
	for index, value := range logEntryConf {
		fmt.Printf("index:%v value:%v\n", index, value)
	}
	// 2.2 派一个哨兵去监视日志收集项的变化（有变化及时通知我的logagent实现热加载配置）

	// 3. 收集日志发往 Kafka
	taillog.Init(logEntryConf)
	newConfChan := taillog.NewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, newConfChan)
	wg.Wait()
	// 3.1 循环每一个收集项，创建 TailObj
}
