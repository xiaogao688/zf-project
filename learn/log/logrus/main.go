package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

func main() {
	// logrus设置 指定输出为json格式
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// 指定输出为标准输出
	logrus.SetOutput(os.Stdout)
	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)

	// logrus使用
	logrus.Debug("Useful debugging information.")
	logrus.Info("Something noteworthy happened!")
	logrus.Warn("You should probably take a look at this.")
	logrus.Error("Something failed but I'm not quitting.")

	// 如何选择和定义 logrus.Fields 的参数
	// 1、明确日志目的：首先确定日志的用途，是用于调试、监控还是审计。根据目的选择需要记录的字段。
	// 2、保持一致性：在整个项目中保持字段命名和使用的一致性，便于日志的统一处理和分析。
	// 3、避免敏感信息：确保不在日志中记录敏感信息，如密码、信用卡信息等。对于需要记录的敏感数据，考虑进行脱敏处理。
	// 4、合理使用嵌套结构：虽然 logrus.Fields 支持嵌套结构，但应避免过度复杂的嵌套，保持日志的简洁和可读性。
	// 5、性能考虑：尽量避免在日志中记录过大的数据量，以免影响日志记录的性能和存储
	logrus.WithFields(logrus.Fields{
		"request_id":  "123e4567-e89b-12d3-a456-426614174000",
		"method":      "GET",
		"url":         "/api/v1/users",
		"status_code": 200,
		"ip":          "192.168.1.100",
		"user_agent":  "Mozilla/5.0",
	}).Info("处理请求完成")

	logger := logrus.WithFields(logrus.Fields{"request_id": 123456789})
	logger.Info("something happened on that request") // 也会记录request_id
	logger.Warn("something not great happened")

	// logrus具有hook能力，允许我们自定义一些日志处理逻辑
	// logrus在记录Levels()返回的日志级别的消息时会触发HOOK,
	// 按照Fire方法定义的内容修改logrus.Entry.
	//type Hook interface {
	//	Levels() []Level
	//	Fire(*Entry) error
	//}
	logrus.AddHook(&DefaultFieldHook{})
	logrus.Debug("debug")
	logrus.Info("info")
	logrus.Warn("warn")
	logrus.Error("error")
}

type DefaultFieldHook struct {
	Writer    *os.File
	LogLevels []logrus.Level
	mu        sync.Mutex
}

// 可以实现写入kafka或者其他逻辑
func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["myHook"] = " MyHookTest "

	// 格式化日志（使用标准文本格式）
	line, err := entry.String()
	if err != nil {
		return err
	}
	fmt.Println("----------------------------")
	fmt.Println(line)
	fmt.Println("----------------------------")
	return nil
}

// 指定输出level级别
func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
