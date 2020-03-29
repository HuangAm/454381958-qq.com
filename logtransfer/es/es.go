package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
	"time"
)

type LogData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

var (
	client *elastic.Client
	ch     chan *LogData
)

// 初始化es准备接受kafka发来的数据
func Init(address string, chanSize int, Nums int) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		return
	}
	fmt.Println("connect to es success.")
	ch = make(chan *LogData, chanSize)
	for i:=0; i< Nums; i++{
		go SendToES()
	}
	return
}

func SendToESChan(msg *LogData) {
	ch <- msg
}

func SendToES() {
	for {
		select {
		case msg := <-ch:
			put1, err := client.Index().Index(msg.Topic).BodyJson(msg).Do(context.Background())
			if err != nil {
				fmt.Println("write to es error:", err)
				continue
			}
			fmt.Printf("Indexed student %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}
}
