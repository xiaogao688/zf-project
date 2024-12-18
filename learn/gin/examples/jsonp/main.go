package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	// JSONP 路由
	r.GET("/jsonp", func(c *gin.Context) {
		// 获取客户端传递的回调函数名（通常是 callback 参数）
		callback := c.Query("callback")
		if callback == "" {
			// 如果没有传 callback，返回错误信息
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing callback parameter",
			})
			return
		}

		// 返回的数据
		data := gin.H{
			"message": "Hello, JSONP!",
			"status":  "success",
		}

		// 通过 Gin 的 JSONP 方法返回结果
		c.JSONP(http.StatusOK, data)
	})

	// 启动服务器
	r.Run(":8080") // 默认监听 8080 端口
}
