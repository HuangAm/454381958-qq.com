// 专门往 kafka 写日志的模块
package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

var (
	client      sarama.SyncProducer // 声明一个全局的连接kafka的生产者client
	logDataChan chan *logData
)

type logData struct {
	topic string
	data  string
}

// init 初始化 client
func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // leader 和 follower 都回复 ACK 后才会回复 ACK
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个 partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在 success channel 返回
	// 连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	logDataChan = make(chan *logData, maxSize)
	go sendToKafka()
	return
}

// 该函数只把日志数据发送到内部的 channel 中
func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

// 真正往 kafka 发送日志
func sendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送到kafka
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid: %v  offset: %v\n", pid, offset)
			fmt.Println("发送成功!")
		default:
			time.Sleep(time.Millisecond * 50)
		}

	}
}
