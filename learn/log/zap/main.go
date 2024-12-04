package main

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction() // 适合在生产环境。json格式，高性能、结构化日志、适合与日志收集工具集成
	//zap.NewExample() // 适合用在测试代码中。简化配置、适合展示和测试、不推荐用于实际应用
	//zap.NewDevelopment() // 适合开发环境中使用，可读性强、详细的调试信息、适合快速调试和开发过程

	defer logger.Sync() // 将缓存中的日志刷新到磁盘文件中

	// 为了提高性能，Logger没有使用interface和反射，并且Logger只支持结构化的日志，所以在使用Logger时，需要指定具体的类型和key-value格式的日志字段
	url := "http://marmotedu.com"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	// 用logger.Sugar()即可创建SugaredLogger。SugaredLogger的使用比Logger简单，但性能比Logger低 50% 左右，可以用在调用次数不高的函数中
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	// 自定义方法
	// Level：日志级别。
	// Development：设置Logger的模式为development模式。
	// DisableCaller：禁用调用信息. 该字段值为 true 时, 日志中将不再显示该日志所在的函数调用信息。
	// DisableStacktrace：禁用自动堆栈跟踪捕获。
	// Sampling：流控配置, 也叫采样. 单位是每秒钟, 作用是限制日志在每秒钟内的输出数量, 以防止CPU和IO被过度占用。
	// Encoding：指定日志编码器, 目前仅支持两种编码器：console和json，默认为json。
	// EncoderConfig：编码配置。
	// OutputPaths：配置日志标准输出，可以配置多个日志输出路径, 一般情况可以仅配置标准输出或输出到文件, 如有需求的话, 也可以两者同时配置。
	// ErrorOutputPaths：错误输出路径，可以是多个。
	// InitialFields：初始化字段配置, 该配置的字段会以结构化的形式打印在每条日志输出中。

	// 其中EncoderConfig为编码配置：
	// 常用的设置如下：
	//
	// MessageKey：日志中信息的键名，默认为msg。
	// LevelKey：日志中级别的键名，默认为level。
	// EncodeLevel：日志中级别的格式，默认为小写，如debug/info。
	rawJSON := []byte(`{
    "level":"debug",
    "encoding":"json",
    "outputPaths": ["stdout", "test.log"],
    "errorOutputPaths": ["stderr"],
    "initialFields":{"name":"dj"},
    "encoderConfig": {
      "messageKey": "message",
      "levelKey": "level",
      "levelEncoder": "lowercase"
    }
  }`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	// 创建了一个输出到标准输出和文件test.log的Logger
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("server start work successfully!")

	// zap支持多种选项，选项的使用方式如下
	// AddStacktrace(lvl zapcore.LevelEnabler)：用来在指定级别及以上级别输出调用堆栈。
	//zap.WithCaller(enabled bool)：指定是否在日志输出内容中增加文件名和行号。
	//zap.AddCaller()：与zap.WithCaller(true)等价，指定在日志输出内容中增加行号和文件名。
	//zap. AddCallerSkip(skip int)：指定在调用栈中跳过的调用深度，否则通过调用栈获得的行号可能总是日志组件中的行号。
	//zap. IncreaseLevel(lvl zapcore.LevelEnabler)：提高日志级别，如果传入的lvl比当前logger的级别低，则不会改变日志级别。
	//ErrorOutput(w zapcore.WriteSyncer)：指定日志组件中出现异常时的输出位置。
	//Fields(fs ...Field)：添加公共字段。
	//Hooks(hooks ...func(zapcore.Entry) error)：注册钩子函数，用来在日志打印时同时调用hook方法。
	//WrapCore(f func(zapcore.Core) zapcore.Core)：替换Logger的zapcore.Core。 - Development()：将Logger修改为Development模式。
	logger, _ = zap.NewProduction(zap.AddCaller())
	defer logger.Sync()

	logger.Info("hello world")

	// 预设日志字段,在创建Logger时使用Fields(fs ...Field)选项
	logger = zap.NewExample(zap.Fields(
		zap.Int("userID", 10),
		zap.String("requestID", "fbf54504"),
	))

	logger.Debug("This is a debug message")
	logger.Info("This is a info message")
}
