package youyoulibs

import (
	"context"
	"encoding/json"

	kafka "github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

// 消息的类型
type MsgType int32

const (
	MsgType_Login MsgType = 0
)

type KafkaMsg struct {
	Type MsgType `json:"type,omitempty"`
	Msg  []byte  `json:"msg,omitempty"`
}

type KafkaConfig struct {
	Topic   string   `toml:"topic"`
	Group   string   `toml:"group"`
	Brokers []string `toml:"brokers"`
}

type KafKa struct {
	kafkaPub kafka.SyncProducer
	consumer *cluster.Consumer
	c        *KafkaConfig
}

// 创建kafka实例
// @param c 配置结构
func NewKafka(c *KafkaConfig) *KafKa {
	obj := &KafKa{
		c:        c,
		kafkaPub: newKafkaPub(c),
		consumer: newKafkaSub(c),
	}

	return obj
}

// 创建生产者
// @param c 配置结构
func newKafkaPub(c *KafkaConfig) kafka.SyncProducer {
	kc := kafka.NewConfig()
	kc.Producer.RequiredAcks = kafka.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                  // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := kafka.NewSyncProducer(c.Brokers, kc)
	if err != nil {
		panic(err)
	}
	return pub
}

// 创建消费者
// @param c 配置结构
func newKafkaSub(c *KafkaConfig) *cluster.Consumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(c.Brokers, c.Group, []string{c.Topic}, config)
	if err != nil {
		panic(err)
	}
	return consumer
}

// 生产消息
// @param c
// @param key 消息的key
// @param m_type 消息类型
// @param msg 	消息
func (d *KafKa) Product(c context.Context, key string, m_type MsgType, msg []byte) (err error) {
	pushMsg := &KafkaMsg{
		Type: m_type,
		Msg:  msg,
	}
	b, err := json.Marshal(pushMsg)
	if err != nil {
		return
	}

	m := &kafka.ProducerMessage{
		Key:   kafka.StringEncoder(key),
		Topic: d.c.Topic,
		Value: kafka.ByteEncoder(b),
	}
	if _, _, err = d.kafkaPub.SendMessage(m); err != nil {
		return
	}

	return
}

// 循环消费消息
// @param errf 	收到错误回调函数
// @param notif 收到通知回调函数
// @param msgf  收到消息毁掉函数
func (j *KafKa) Consume(errf func(error), notif func(interface{}), msgf func(*KafkaMsg)) {
	for {
		select {
		case err := <-j.consumer.Errors():
			errf(err)
		case n := <-j.consumer.Notifications():
			notif(n)
		case msg, ok := <-j.consumer.Messages():
			if !ok {
				return
			}

			j.consumer.MarkOffset(msg, "")
			// process push message
			pushMsg := new(KafkaMsg)
			if err := json.Unmarshal(msg.Value, pushMsg); err != nil {
				continue
			}

			msgf(pushMsg)
		}
	}
}

// 关闭消息队列监听
func (j *KafKa) Close() error {
	if j.consumer != nil {
		return j.consumer.Close()
	}

	return nil
}
