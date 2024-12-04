package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

// cobra.Command 是一个结构体，代表一个命令
var rootCmd = &cobra.Command{
	Use:   "hugo",                 // 命令的名称。
	Short: "Hugo 是一个非常快速的静态站点生成器", //  代表当前命令的简短描述。
	Long: `由 spf13 和 Go 中的朋友用爱构建的快速灵活的静态站点生成器。 
完整的文档可在  https://gohugo.io 上找到`, // 表示当前命令的完整描述。
	TraverseChildren: true, // 开启此参数，Cobra 将在执行目标命令之前解析每个命令的本地标志
	Args: func(cmd *cobra.Command, args []string) error { // 对命令参数进行合法性验证
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		if len(args) > 4 {
			return errors.New("the number of args cannot exceed 4")
		}
		if args[0] != "a" {
			return errors.New("first argument must be 'a'")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) { // 属性是一个函数，当执行命令时会调用此函数。
		fmt.Println("run hugo...", GlobalLogo, LocalLogo)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) { // 在 PreRun 函数执行之前执行，对此命令的子命令同样生效。
		fmt.Println("hugo PersistentPreRun")
	},
	PreRun: func(cmd *cobra.Command, args []string) { // 在 Run 函数执行之前执行
		fmt.Println("hugo PreRun")
	},
	PostRun: func(cmd *cobra.Command, args []string) { // 在 Run 函数执行之后执行。
		fmt.Println("hugo PostRun")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) { // 在 PostRun 函数执行之后执行，对此命令的子命令同样生效。
		fmt.Println("hugo PersistentPostRun")
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { // PersistentPreRun 报错时执行
		fmt.Println("hugo PersistentPreRunE")
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error { // PostRun 报错时执行
		fmt.Println("hugo PreRunE")
		return errors.New("PreRunE err")
	},
}

/*
Cobra 内置了以下验证函数：
	NoArgs：如果存在任何命令参数，该命令将报错。
	ArbitraryArgs：该命令将接受任意参数。
	OnlyValidArgs：如果有任何命令参数不在 Command 的 ValidArgs 字段中，该命令将报错。
	MinimumNArgs(int)：如果没有至少 N 个命令参数，该命令将报错。
	MaximumNArgs(int)：如果有超过 N 个命令参数，该命令将报错。
	ExactArgs(int)：如果命令参数个数不为 N，该命令将报错。
	ExactValidArgs(int)：如果命令参数个数不为 N，或者有任何命令参数不在 Command 的 ValidArgs 字段中，该命令将报错。
	RangeArgs(min, max)：如果命令参数的数量不在预期的最小数量 min 和最大数量 max 之间，该命令将报错。
*/

var GlobalLogo bool // 持久标志
var LocalLogo bool
var author string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&GlobalLogo, "globalLogo", "g", false, "全局标识") // 该标志将可用于它所分配的命令以及该命令下的所有子命令。
	rootCmd.Flags().BoolVarP(&LocalLogo, "localLogo", "l", false, "局部标识")             // 标志也可以是本地的，这意味着它只适用于该指定命令。
	err := rootCmd.MarkFlagRequired("localLogo")                                      // 默认情况下，标志是可选的。我们可以将其标记为必选.这个只能配置本地
	if err != nil {
		log.Panic(err)
	}
}

// Execute 是命令的执行入口
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 初始化viper
func init() {
	rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}

//// viper的文件初始化
//var cfgFile string
//
//func init() {
//	cobra.OnInitialize(initConfig) // cobra.OnInitialize() 来初始化配置文件
//	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file")
//}
//
//func initConfig() {
//	if cfgFile != "" {
//		viper.SetConfigFile(cfgFile)
//	} else {
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".cobra")
//	}
//
//	if err := viper.ReadInConfig(); err != nil {
//		fmt.Println("Can't read config:", err)
//		os.Exit(1)
//	}
//}
