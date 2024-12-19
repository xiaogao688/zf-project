package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()
	router.GET("/someDataFromReader", func(c *gin.Context) {
		response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})

	router.GET("/test", func(c *gin.Context) {
		// 创建一个管道
		pr, pw := io.Pipe()

		// 启动一个协程，用于写入数据
		go func() {
			for i := 1; i <= 10; i++ {
				// 写入数据
				message := "1234567890\n"
				_, err := pw.Write([]byte(message))
				if err != nil {
					fmt.Println("Error writing to pipe:", err)
					return
				}
				// 模拟写入间隔
				time.Sleep(1 * time.Second)
			}
			// 关闭写入器
			pw.Close()
		}()

		// 当字符数少于contentLength时，应该会报错
		c.DataFromReader(http.StatusOK, 55, "Content-Type: text/plain; charset=utf-8", pr, nil)
	})

	router.Run(":8080")
}
