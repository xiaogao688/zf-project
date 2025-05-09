package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// 请根据你的环境修改这两个地址：
	// - 内部 DNS 名称：dosec-kafka:9092
	// - 外部 IP+端口：192.168.9.11:31092
	brokers := []string{
		"192.168.4.95:30092",
	}

	// Kafka 客户端配置
	config := sarama.NewConfig()
	// 设置 Kafka 版本，确保与你的集群版本兼容
	config.Version = sarama.V2_5_0_0

	// 启用 SASL/PLAIN（无 TLS）
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = "admin"              // 替换成你的用户名
	config.Net.SASL.Password = "36UdcuHtpZSfRj" // 替换成你的密码
	config.Net.SASL.Handshake = true
	config.Net.TLS.Enable = false

	// 连接超时设置（可选）
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	// 创建一个普通 client（只是建立连接、获取元数据）
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatalf("无法创建 Kafka 客户端: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("关闭客户端时出错: %v", err)
		}
	}()

	// 列出所有 topic
	topics, err := client.Topics()
	if err != nil {
		log.Fatalf("获取 topic 列表失败: %v", err)
	}

	fmt.Println("当前集群中的 Topics：")
	for _, t := range topics {
		fmt.Printf("  • %s\n", t)
	}
}
