package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{
		"192.168.4.95:30092",
	}
	topic := "dosec_asset_change"

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0

	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = "admin"
	config.Net.SASL.Password = "36UdcuHtpZSfRj"
	config.Net.SASL.Handshake = true
	config.Net.TLS.Enable = false

	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	// 生产者配置（异步）
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("无法创建 Kafka 生产者: %v", err)
	}
	defer producer.Close()

	// 发送一条测试消息到 topic
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("key1"),
		Value: sarama.StringEncoder("这是一个测试消息！"),
	}
	producer.Input() <- msg
	log.Println("发送消息成功: 这是一个测试消息！")

	// 消费者配置
	consumerGroup := "test-group"
	client, err := sarama.NewConsumerGroup(brokers, consumerGroup, config)
	if err != nil {
		log.Fatalf("无法创建 Kafka 消费者组: %v", err)
	}
	defer client.Close()

	// 设置消费处理器
	handler := ConsumerGroupHandler{}

	// 捕获退出信号
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, os.Interrupt)
		<-sigterm
		cancel()
	}()

	go func() {
		for {
			if err := client.Consume(ctx, []string{topic}, handler); err != nil {
				log.Fatalf("消费失败: %v", err)
			}
			// 检查是否取消
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-ctx.Done()
	log.Println("消费结束")
}

// 消费者组处理器
type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("接收到消息: topic=%s partition=%d offset=%d key=%s value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
