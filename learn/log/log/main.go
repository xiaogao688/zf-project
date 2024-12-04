package main

import (
	"log"
	"os"
	"time"
)

func main() {
	// 基本日志输出
	log.Println("这是一个普通的日志消息。")
	log.Printf("当前时间是: %s\n", time.Now().Format(time.DateTime))
	log.Print("使用 Print 方法输出的日志。")

	// 设置日志前缀
	log.SetPrefix("INFO: ")
	log.Println("设置了日志前缀后的日志消息。")

	// 设置日志格式标志
	// 标志选项包括：
	// log.Ldate：日期（例如 2009/01/23）
	// log.Ltime：时间（例如 01:23:23）
	// log.Lmicroseconds：时间，精确到微秒
	// log.Llongfile：完整的文件名和行号
	// log.Lshortfile：文件名和行号
	// log.LUTC：使用 UTC 时间
	// log.LstdFlags：Ldate | Ltime
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("设置了日志格式标志后的日志消息。")

	// 创建一个自定义日志记录器，输出到文件
	file, err := os.OpenFile("D:\\code\\go-project\\learn\\log\\log\\custom.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	customLogger := log.New(file, "CUSTOM: ", log.Ldate|log.Ltime|log.Lshortfile)
	customLogger.Println("这是通过自定义日志记录器输出的日志消息。")

	// 演示 Fatal 函数（会终止程序）
	// 注意：下面的代码会终止程序，请根据需要取消注释以测试
	/*
		log.Fatal("这是一个致命错误，程序将终止。")
		log.Println("这条日志不会被执行。")
	*/

	// 演示 Panic 函数（会引发恐慌）
	// 注意：下面的代码会引发恐慌，请根据需要取消注释以测试
	/*
		log.Panic("这是一个恐慌消息，程序将引发恐慌。")
		log.Println("这条日志不会被执行。")
	*/

	// 使用不同的日志级别（通过创建不同的记录器）
	errorLogger := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger.Println("这是一个错误级别的日志消息。")

	debugLogger := log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger.Println("这是一个调试级别的日志消息。")

	// 示例完成
	log.Println("日志记录示例程序已完成。")
}
