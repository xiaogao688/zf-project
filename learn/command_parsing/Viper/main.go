package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"strings"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	//DefaultViper()
	//ViperReadConfigFile()
	//ViperWriteConfigFile()
	//ViperWatchConfigFile()
	//ViperIOReaderConfigFile()
	ViperFlagConfigFile()
}

// DefaultViper 使用默认值
// 覆盖设置
// 建立别名
// 提取子树
// 反序列化
// 序列化成字符串
// 启用其他实例
func DefaultViper() {
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("Port", 12345)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})

	viper.Set("Port", true) // 覆盖Port设置

	viper.RegisterAlias("Port", "Verbose") // 注册别名（此处Port和Verbose建立了别名）

	// 提取子树，例如 app.cache1.abc = 1, 可以使用这种方式
	subv := viper.Sub("app.cache1")
	subv.GetString("abc")

	// 可以将配置反序列化到指定结构体中
	//viper.Unmarshal(&c)

	// 序列化成字符串
	//c := viper.AllSettings()

	// 启用其他实例
	//x := viper.New()

	fmt.Println("Taxonomies: ", viper.GetStringMapString("Taxonomies"))
	fmt.Println("Port: ", viper.GetInt64("Port"))
}

// ViperReadConfigFile 读取配置文件
func ViperReadConfigFile() {
	viper.SetConfigFile("D:\\code\\my_repo\\zf-project\\learn\\command_parsing\\Viper\\config.json") // 指定配置文件路径

	// 和上面SetConfigFile同时使用可能有点问题
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项

	viper.AddConfigPath("/etc/appname/")  // 查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")              // 还可以在工作目录中查找配置

	// 查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			fmt.Println("配置文件未找到错误", err)
		} else {
			// 配置文件被找到，但产生了另外的错误
			fmt.Println("配置文件被找到，但产生了另外的错误", err)
		}
	}
	fmt.Println("datastore.metric.host: ", viper.GetString("datastore.metric.host"))
}

// ViperWriteConfigFile 写入配置文件,将当前配置导出到指定位置
/*
WriteConfig - 将当前的viper配置写入预定义的路径并覆盖（如果存在的话）。如果没有预定义的路径，则报错。
SafeWriteConfig - 将当前的viper配置写入预定义的路径。如果没有预定义的路径，则报错。如果存在，将不会覆盖当前的配置文件。
WriteConfigAs - 将当前的viper配置写入给定的文件路径。将覆盖给定的文件(如果它存在的话)。
SafeWriteConfigAs - 将当前的viper配置写入给定的文件路径。不会覆盖给定的文件(如果它存在的话)。
*/
func ViperWriteConfigFile() {
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("Port", 12345)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})

	viper.SetConfigFile("D:\\code\\my_repo\\zf-project\\learn\\command_parsing\\Viper\\config2.json") // 指定配置文件路径
	viper.SetConfigFile("D:\\code\\my_repo\\zf-project\\learn\\command_parsing\\Viper\\config2.yaml") // 指定配置文件路径
	viper.WriteConfig()                                                                               // 将当前配置写入“viper.AddConfigPath()”和“viper.SetConfigName”设置的预定义路径

	//viper.SafeWriteConfig()
	//viper.WriteConfigAs("/path/to/my/.config")
	//viper.SafeWriteConfigAs("/path/to/my/.config") // 因为该配置文件写入过，所以会报错
	//viper.SafeWriteConfigAs("/path/to/my/.other_config")
}

// ViperWatchConfigFile 监控并重新读取配置文件
func ViperWatchConfigFile() {
	viper.SetConfigFile("D:\\code\\my_repo\\zf-project\\learn\\command_parsing\\Viper\\config.json")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name) // 修改的文件名称
		fmt.Println(e.Op.String())                  // 修改类型
	})
	time.Sleep(1 * time.Hour)
}

// ViperIOReaderConfigFile 从io.reader 读取
func ViperIOReaderConfigFile() {
	viper.SetConfigType("yaml") // 或者 viper.SetConfigType("YAML")

	// 任何需要将此配置添加到程序中的方法。
	var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	viper.Get("name") // 这里会得到 "steve"
}

// ViperENVConfigFile 使用环境变量
func ViperENVConfigFile() {
	viper.SetEnvPrefix("t_")         // 查找环境变量的前缀
	viper.BindEnv("key", "env-name") // 第一个参数是键名称，后面是环境变量的名称。如果只有第一个参数会查找: 前缀+ “_” +键名全部大写, 显式提供ENV变量名（第二个参数）时，它不会自动添加前缀,也是全大写。
	viper.AutomaticEnv()             // 可以与SetEnvPrefix结合使用,Viper会在发出viper.Get请求时随时检查环境变量。

	viper.SetEnvKeyReplacer(&strings.Replacer{}) // 允许使用strings.Replacer对象在一定程度上重写 Env 键,例如将环境变量的'-'变成'_'

	viper.AllowEmptyEnv(false) // 可以将空环境变量视为已设置
}

// ViperFlagConfigFile 允许通过flag设置环境变量
func ViperFlagConfigFile() {
	pflag.Int("flagname", 1234, "help message for flagname")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	i := viper.GetInt("flagname") // 从viper而不是从pflag检索值
	fmt.Println(i)
}

// ViperKVConfigFile 通过KV导入，
// 在Viper中启用远程支持，需要在代码中匿名导入viper/remote这个包 import _ "github.com/spf13/viper/remote"

func ViperKVConfigFile() {
	// 或者你可以创建一个新的viper实例
	var runtime_viper = viper.New()

	runtime_viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001", "/config/hugo.yml")
	runtime_viper.SetConfigType("yaml") // 因为在字节流中没有文件扩展名，所以这里需要设置下类型。支持的扩展名有 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// 第一次从远程读取配置
	err := runtime_viper.ReadRemoteConfig()
	if err != nil {
		fmt.Println(err)
	}

	type A struct {
		a string
		b int
	}
	runtime_conf := A{}

	// 反序列化
	runtime_viper.Unmarshal(&runtime_conf)

	// 开启一个单独的goroutine一直监控远端的变更
	go func() {
		for {
			time.Sleep(time.Second * 5) // 每次请求后延迟一下

			// 目前只测试了etcd支持
			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				log.Errorf("unable to read remote config: %v", err)
				continue
			}

			// 将新配置反序列化到我们运行时的配置结构体中。你还可以借助channel实现一个通知系统更改的信号
			runtime_viper.Unmarshal(&runtime_conf)
		}
	}()
}
