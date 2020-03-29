package test

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"testing"
)

type Student struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func (s *Student) run() {
	fmt.Printf("【%s】在跑", s.Name)
}

func TestESInsert(t *testing.T) {
	// 1. 初始化连接，得到一个client
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil{
		panic(err)
	}
	//client, err := elastic.NewClient(elastic.SetURL("http://192.168.0.102:9200"))
	//exists, err := client.IndexExists("student").Do(context.Background())
	//fmt.Printf("%#v\n", exists)
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//fmt.Println("connect to es success")
	p1 := Student{Name: "wuyongqiang", Age: 22, Married: false}
	put1, err := client.Index().Index("student").BodyJson(p1).Do(context.Background())
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Printf("Indexed student %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
