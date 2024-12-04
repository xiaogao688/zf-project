package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
)

// -alsologtostderr：同时将日志打印到文件和标准错误输出。
// -log_backtrace_at value：指定代码运行到指定行时，把该代码的栈信息打印出来。
// -log_dir：指定日志存储的文件夹。
// -logtostderr：日志打印到标准错误输出，而不是文件中。
// -stderrthreshold value：指定大于或者等于该级别的日志才会被输出到标准错误输出，默认为ERROR。
// -v value：指定日志级别。
// -vmodule value：对不同的文件使用不同的日志级别。

// glog的日志输出格式为：<header>] <message>，其中header的格式为：Lmmdd hh:mm:ss.uuuuuu threadid file:line：
//
// Lmmdd：L代表了glog的日志级别：I -> INFO、W -> WARNING、E -> ERROR、F -> FATAL。
// hh:mm:ss.uuuuuu：代表了时间信息，例如10:12:32.995956。
// threadid，是进程PID，即os.Getpid()的返回值。
// file:line：指明了打印日志的位置：文件名和行号。

func main() {
	// 解析命令行标志
	flag.Parse()
	defer glog.Flush()

	// 基本信息日志
	glog.Info("这是一个普通的信息日志。")
	glog.Infof("当前时间是: %s", time.Now().Format(time.DateTime))
	glog.Infoln("使用 Infoln 方法输出的日志。")

	// 警告日志
	glog.Warning("这是一个警告日志。")
	glog.Warningf("警告：%s", "内存使用量过高")
	glog.Warningln("使用 Warningln 方法输出的警告日志。")

	// 错误日志
	glog.Error("这是一个错误日志。")
	glog.Errorf("错误：%s", "无法连接到数据库")
	glog.Errorln("使用 Errorln 方法输出的错误日志。")

	// 致命错误日志（会终止程序）
	// 注意：下面的代码会终止程序，请根据需要取消注释以测试
	/*
		glog.Fatal("这是一个致命错误，程序将终止。")
		glog.Fatalln("使用 Fatalln 方法输出的致命错误日志。")
	*/

	// 使用不同的详细程度（verbosity）
	// 通过设置 -v 标志来控制，例如：-v=2，也可以控制指定文件的日志级别，例如 -vmodule=foo=5 可以设置对foo.go文件使用5级别 语法格式为：-vmodule=file1=2,file2=1,fs*=3
	// V越小，说明日志级别越高
	glog.V(1).Info("这是详细程度为1的信息。")
	glog.V(2).Info("这是详细程度为2的信息。")
	glog.V(3).Info("这是详细程度为3的信息。")

	// 模拟长时间运行的程序，生成更多日志
	for i := 0; i < 5; i++ {
		glog.Infoln("循环日志", "iteration", i)
		time.Sleep(500 * time.Millisecond)
	}

	// traceLocation可以打印出指定位置的栈信息（-log_backtrace_at=filename:line_number）
	// 执行时加入参数 -log_backtrace_at=root.go:60
	glog.Info("This is info message")

	// 使用 Flush 确保所有日志被写入
	glog.Flush()
}
