package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// 全局上下文
var ctx = context.Background()

// 常量定义，包括 Redis 地址、Stream 名称以及消费者组名称
const (
	redisAddr  = "192.168.6.7:30379"
	streamName = "config_stream"

	consumerGroup1 = "consumer_group_1"
	consumerGroup2 = "consumer_group_2"
	consumerGroup3 = "consumer_group_3"
)

func main() {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "*****",
	})

	// 创建三个消费者组。这里使用 XGROUP CREATE MKSTREAM 来创建组，同时自动创建 stream（如果不存在）。
	groups := []string{consumerGroup1, consumerGroup2, consumerGroup3}
	for _, group := range groups {
		err := rdb.XGroupCreateMkStream(ctx, streamName, group, "$").Err()
		if err != nil {
			// 如果消费者组已存在，忽略 BUSYGROUP 错误
			if err.Error() != "BUSYGROUP Consumer Group name already exists" {
				log.Printf("创建消费者组 %s 失败: %v\n", group, err)
			} else {
				log.Printf("消费者组 %s 已存在，跳过创建\n", group)
			}
		} else {
			log.Printf("成功创建消费者组: %s\n", group)
		}
	}

	// 启动三个消费者，每个消费者加入不同的消费者组，并分别设定消费者名称
	go consumer(rdb, consumerGroup1, "consumer1")
	go consumer(rdb, consumerGroup2, "consumer2")
	go consumer(rdb, consumerGroup3, "consumer3")

	// 启动生产者，定时发布消息到 stream
	go producer(rdb)

	// 阻塞主 goroutine
	select {}
}

// producer 定时向 Redis Stream 中发布消息
func producer(rdb *redis.Client) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	msgCount := 0
	for {
		<-ticker.C
		msgCount++
		values := map[string]interface{}{
			"id":      msgCount,
			"content": fmt.Sprintf("配置更新消息 %d", msgCount),
		}
		// 使用 XADD 命令向 stream 添加消息
		msgID, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: values,
			MaxLen: 100,
			Approx: true, // 可选；开启后修剪操作为近似修剪，可提高写入性能
		}).Result()
		if err != nil {
			log.Printf("消息添加失败: %v\n", err)
		} else {
			log.Printf("生产者发布消息：%s, 内容: %v\n", msgID, values)
		}
	}
}

// consumer 负责从其所属消费者组内阻塞读取并处理消息
func consumer(rdb *redis.Client, group, consumerName string) {
	for {
		// 使用 XREADGROUP 从 stream 中读取消息，等待最长 5 秒
		streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Count:    1,               // 每次读取一条消息
			Block:    5 * time.Second, // 阻塞等待 5 秒消息到来
		}).Result()
		if err != nil {
			// 如果超时未获取到消息，继续阻塞等待
			if err == redis.Nil {
				continue
			}
			log.Printf("[%s] 读取消息出错: %v\n", consumerName, err)
			time.Sleep(time.Second)
			continue
		}
		// 遍历读取到的消息
		for _, stream := range streams {
			for _, message := range stream.Messages {
				log.Printf("[%s] 接收到消息: ID=%s, Values=%v\n", consumerName, message.ID, message.Values)
				// 这里可加入消息处理逻辑
				// 处理完成后，调用 XACK 命令确认消息
				if err = rdb.XAck(ctx, streamName, group, message.ID).Err(); err != nil {
					log.Printf("[%s] 确认消息失败: %v\n", consumerName, err)
				} else {
					log.Printf("[%s] 消息 %s 已确认\n", consumerName, message.ID)
				}
			}
		}
	}
}
