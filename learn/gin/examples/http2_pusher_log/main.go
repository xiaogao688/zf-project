package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置日志文件路径
	logFilePath := "D:\\code\\my_repo\\zf-project\\learn\\gin\\examples\\http2_pusher_log\\test.log"

	// 创建 Gin 实例
	router := gin.Default()

	// 设置模板文件
	router.LoadHTMLGlob("D:\\code\\my_repo\\zf-project\\learn\\gin\\examples\\http2_pusher_log\\templates\\*")

	// 首页，提供实时日志展示的 HTML 页面
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	// SSE 路由，用于推送日志更新
	router.GET("/logs", func(c *gin.Context) {
		c.Stream(func(w io.Writer) bool {
			file, err := os.Open(logFilePath)
			if err != nil {
				log.Printf("Failed to open log file: %v\n", err)
				return false
			}
			defer file.Close()

			// 创建一个文件扫描器
			scanner := bufio.NewScanner(file)

			// 跳过已经读取的旧日志
			for scanner.Scan() {
				line := scanner.Text()
				// 发送 SSE 数据到客户端
				fmt.Println(line)
				c.SSEvent("message", line)
				time.Sleep(1 * time.Second)
			}

			// 持续监听日志文件的新内容
			for {
				if scanner.Scan() {
					line := scanner.Text()
					// 发送 SSE 数据到客户端
					fmt.Println(line)
					c.SSEvent("message", line)
				} else {
					// 等待新内容写入
					fmt.Println("sleep")
					time.Sleep(1 * time.Second)
				}
			}
		})
	})

	// 启动 HTTP/2 服务
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Server started on https://localhost:8080")
	log.Fatal(server.ListenAndServeTLS("D:\\code\\my_repo\\zf-project\\learn\\gin\\ca-crt\\server.crt", "D:\\code\\my_repo\\zf-project\\learn\\gin\\ca-crt\\server.key")) // 启用 HTTPS
}
