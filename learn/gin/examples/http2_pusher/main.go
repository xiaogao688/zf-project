package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 路由
	router.GET("/", func(c *gin.Context) {
		// 检查是否支持 HTTP/2 推送
		if pusher := c.Writer.Pusher(); pusher != nil {
			// 推送静态资源
			err := pusher.Push("/static/style.css", &http.PushOptions{
				Method: "GET",
				Header: http.Header{
					"Accept-Encoding": c.Request.Header["Accept-Encoding"],
				},
			})
			if err != nil {
				log.Printf("Failed to push resource: %v", err)
			}
		}

		// 渲染主页面
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "HTTP/2 Server Push Example",
		})
	})

	// 提供静态资源
	router.Static("/static", "D:\\code\\my_repo\\zf-project\\learn\\gin\\examples\\http2_pusher\\static")

	// 加载模板
	router.LoadHTMLGlob("D:\\code\\my_repo\\zf-project\\learn\\gin\\examples\\http2_pusher\\templates\\*")

	// 启动服务
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Starting server on https://localhost:8080")
	log.Fatal(server.ListenAndServeTLS("D:\\code\\my_repo\\zf-project\\learn\\gin\\ca-crt\\server.crt", "D:\\code\\my_repo\\zf-project\\learn\\gin\\ca-crt\\server.key")) // 启用 HTTPS
}
