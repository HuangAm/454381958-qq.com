#项目开发
## logagent 模块
	1. 读日志，使用 github.com/hpcloud/tail 包来实现
	2. 向 kafka 中写日志，使用 sarama 包
	
## 启动
1. 启动 zookeeper：bin\windows\zookeeper-server-start.bat config\zookeeper.properties
2. 启动 kafka: bin\windows\kafka-server-start.bat config\server.properties
3. 启动 etcd: etcd.exe
4. 启动 elastic: bin\elasticsearch.bat
5. 启动 kibana : bin\kibana.bat\
6. 启动 logagent
7. 启动 logtransfer
	
#相关知识
## Kafka要点
1. kafka 集群架构
	1. broker(节点)
	2. topic(日志类别)
	3. partition(分区，把同一个topic分成不同的分区，提高负载)
		1. leader: 分区的主节点
		2. follower: 分区的从节点
	4. Consumer Group 消费者组
2. 生产者往 kafka 发送数据的流程
	1. 生产者从 kafka 集群中获取分区 leader 信息
	2. 生产者将信息发送给 leader
	3. leader 将消息写入本地磁盘
	4. follower 从 leader 拉取消息数据
	5. follower 将消息写入本地磁盘后向 leader 发送 ACK
	6. leader 收到所有的 follower 的 ACK 之后向生产者发送 ACK
3. kafka 选择分区的模式（3种）
	1. 指定往哪个分区写
	2. 制定key, kafka 根据 key 做哈希，决定写哪个分区
	3. 轮询方式
4. 生产者往kafka发送数据的模式
	1. 0：把数据发给 leader 就成功
	2. 1：把数据发送给 leader, 等待 leader 回 ACK
	3. all: 把数据发给 leader, follower 从 leader 拉取数据，回复 ACK 给leader,leader 再回复 ACK 给生产者
5. 每个partition都是一个有序且不可变的消息记录集合。当新的数据写入时，就被追加到partition的末尾。在每个partition中，每条消息都会被分配一个顺序的唯一标识，这个标识被称为offset, 即偏移量。
6. 在同一个消费者组中，每个消费者实例可以消费多个分区，但是每个分区最多只能被同一个消费者组中的某一个消费者消费

## etcd
etcd 是使用Go语言开发的一个开源的、高可用的分布式key-value存储系统，可以用于配置共享和服务的注册和发现。

类似的项目有zookeeper和consul。

etcd 具有以下特点：
- 完全复制：集群中的每个节点都可以使用完整的存档
- 高可用性：etcd可用于避免硬件的单点故障或网络问题
- 一致性：每次读取都会返回跨多主机的最新写入
- 简单：包括一个定义良好、面向用户的API(gRPC)
- 安全：实现了带有可选的客户端证书身份验证的自动化TLS
- 快速：每秒10000次写入的基准速度
- 可靠：使用Raft算法实现了强一致、高可用的服务存储目录