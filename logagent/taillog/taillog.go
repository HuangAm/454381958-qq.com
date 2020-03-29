// 专门从日志文件收集日志的模块
package taillog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"logagent/kafka"
)

// 一个日志收集的任务
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
	// 为了能实现退出 t.run()
	ctx context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:  path,
		topic: topic,
		ctx: ctx,
		cancelFunc:cancel,
	}
	tailObj.init()
	return
}

// 根据路径去打开对应操作
func (t *TailTask) init() {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tial file failed, err:", err)
	}
	go t.run() // 直接去采集日志发送到kafka
}

func (t *TailTask) ReadChan() <-chan *tail.Line {
	return t.instance.Lines
}

func (t *TailTask) run() {
	for {
		select {
		case <- t.ctx.Done():
			return
		case line := <-t.instance.Lines:
			//kafka.SendToKafka(t.topic, line.Text) // 同步
			kafka.SendToChan(t.topic, line.Text)
		}
	}
}
