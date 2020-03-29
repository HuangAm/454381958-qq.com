package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logtransfer/conf"
	"logtransfer/es"
	"logtransfer/kafka"
)

func main()  {
	// 0. 加载配置文件
	var cfg conf.LogTransferCfg
	err := ini.MapTo(&cfg, "./conf/cfg.ini")
	if err != nil{
		fmt.Printf("init config, err:%v\n", err)
		return
	}
	fmt.Printf("%#v\n", cfg)
	// 1. 初始化es
	// 1.1 初始化一个ES连接的client
	err = es.Init(cfg.EsCfg.Address, cfg.EsCfg.ChanSize, cfg.EsCfg.Nums)
	if err != nil{
		fmt.Printf("init ES consumer failed, err:%v\n", err)
		return
	}
	fmt.Println("init es success.")
	// 2. 初始化
	// 2.1 连接kafka, 创建分区消费者
	// 2.2 每个分区的消费者分别取出数据 通过SendToES()将数据发往ES
	err = kafka.Init(cfg.KafkaCfg.Address, cfg.KafkaCfg.Topic)
	if err != nil{
		fmt.Printf("init kafka failed, err:%v", err)
		return
	}
	fmt.Println("init kafka success.")
	select {}
}
