package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

func main() {
	//compatibleFlag()
	//customizeFlag()
	//NormalizeFlag()
	//NoOptDefValFlag()
	MarkDeprecatedOrMarkHidden()
}

// compatibleFlag 兼容原生的flag
func compatibleFlag() {
	// 兼容flag的用法
	var port *int = pflag.Int("port", 1234, "help message for prot")
	var port2 int
	pflag.IntVar(&port2, "port2", 8080, "help message for port2")
	fmt.Printf("port: %d port2: %d \n", *port, port2) // port: 1234 port2: 8080
	pflag.Parse()

	// 非选项参数数也都被 pflag 识别并记录了下来 例如：.\main.exe --port 2222 a b c
	fmt.Printf("NFlag: %v\n", pflag.NFlag()) // 返回已设置的命令行标志个数
	fmt.Printf("NArg: %v\n", pflag.NArg())   // 返回处理完标志后剩余的参数个数
	fmt.Printf("Args: %v\n", pflag.Args())   // 返回处理完标志后剩余的参数列表
	fmt.Printf("Arg(1): %v\n", pflag.Arg(1)) // 返回处理完标志后剩余的参数列表中第 i 项

}

// 自定义结构
type host struct {
	value string
}

func (h *host) String() string {
	return h.value
}

func (h *host) Set(v string) error {
	h.value = v
	return nil
}

func (h *host) Type() string {
	return "host"
}

// customizeFlag 自定义flag
func customizeFlag() {
	// 自定义 FlagSet 对象
	flagset := pflag.NewFlagSet("test", pflag.ExitOnError)
	var ip = flagset.IntP("ip", "i", 1234, "help message for ip")
	var boolVar bool
	flagset.BoolVarP(&boolVar, "boolVar", "b", true, "help message for boolVar")
	var h host
	flagset.VarP(&h, "host", "H", "help message for host")
	flagset.SortFlags = false // 取消排序
	flagset.Parse(os.Args[1:])
	fmt.Printf("ip: %d\n", *ip)
	fmt.Printf("boolVar: %t\n", boolVar)
	fmt.Printf("host: %+v\n", h)
}

// NormalizeFlag 借助 pflag.NormalizedName 我们能够给标志起一个或多个别名、规范化标志名等
// .\main.exe --my_flag 222 输出 myFlag: 222
func NormalizeFlag() {
	flagset := pflag.NewFlagSet("test", pflag.ExitOnError)

	var ip = flagset.IntP("new-flag-name", "i", 1234, "help message for new-flag-name")
	var myFlag = flagset.IntP("my-flag", "m", 1234, "help message for my-flag")

	flagset.SetNormalizeFunc(normalizeFunc)
	flagset.Parse(os.Args[1:])

	fmt.Printf("ip: %d\n", *ip)
	fmt.Printf("myFlag: %d\n", *myFlag)
}

func normalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	// 别名
	switch name {
	case "old-flag-name":
		name = "new-flag-name"
		break
	}

	// --my-flag == --my_flag == --my.flag
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}

	return pflag.NormalizedName(name)
}

// NoOptDefValFlag 命令行上设置了标志而没有参数选项，则标志将设置为 NoOptDefVal 指定的值
func NoOptDefValFlag() {
	var ip = pflag.IntP("flagname", "f", 1234, "help message")
	pflag.Lookup("flagname").NoOptDefVal = "4321"

	pflag.Parse()

	fmt.Println("ip: ", *ip)
}

// MarkDeprecatedOrMarkHidden 弃用/隐藏标志
// .\main.exe --ip 1 --boolVar=false -H localhost
// 弃用和隐藏的标签都可以使用，但在-h时不展示，使用弃用标签会警告: Flag --ip has been deprecated, deprecated
func MarkDeprecatedOrMarkHidden() {
	flags := pflag.NewFlagSet("test", pflag.ExitOnError)

	var ip = flags.IntP("ip", "i", 1234, "help message for ip")

	var boolVar bool
	flags.BoolVarP(&boolVar, "boolVar", "b", true, "help message for boolVar")

	var h string
	flags.StringVarP(&h, "host", "H", "127.0.0.1", "help message for host")

	// 弃用标志
	flags.MarkDeprecated("ip", "deprecated")                              // 弃用 ip
	flags.MarkShorthandDeprecated("boolVar", "please use --boolVar only") // 弃用boolVar的短标签

	// 隐藏标志
	flags.MarkHidden("host")

	flags.Parse(os.Args[1:])

	fmt.Printf("ip: %d\n", *ip)
	fmt.Printf("boolVar: %t\n", boolVar)
	fmt.Printf("host: %+v\n", h)
}
