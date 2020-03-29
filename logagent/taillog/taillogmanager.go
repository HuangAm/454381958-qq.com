package taillog

import (
	"fmt"
	"logagent/etcd"
	"time"
)

var tskMgr *tailLogManager

type tailLogManager struct {
	logEntry []*etcd.LogEntry
	tskMap map[string]*TailTask // 配置热更新
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry)  {
	tskMgr = &tailLogManager{
		logEntry:logEntryConf,
		tskMap: make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry), // 无缓冲区的通道
	}
	for _, logEntry := range logEntryConf{
		// 初始化的时候起了多少个tailtask都要记录下来，为了后续判断方便
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tskMgr.tskMap[mk] = tailObj
	}
	go tskMgr.run()
}

// 监听自己的通道 newConfChan 有了新的配置过来之后就做对应的处理
// 1. 配置新增
// 2. 配置删除
// 3. 配置变更
func (t *tailLogManager) run()  {
	for{
		select {
		case newConf := <- t.newConfChan:
			for _, conf := range newConf{
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_,ok := t.tskMap[mk]
				if !ok{
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.tskMap[mk] = tailObj
				}
			}
			// 找出原来t.logEntry有，但是newConf中没有的，要删除
			for _, c1 := range t.logEntry{
				var isDelete bool
				for _, c2 := range newConf{
					if c2.Path == c1.Path && c2.Topic == c1.Topic{
						continue
					}

				}
				if isDelete{
					// 把c1对应的tailObj给停掉
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc()
				}
			}
			fmt.Println("新的配置来了！", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

// 暴漏 taskManager 的 newConfChan
func NewConfChan() chan<- []*etcd.LogEntry{
	return tskMgr.newConfChan
}