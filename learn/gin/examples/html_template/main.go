package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	// *.html 加载顶层以html为后缀的模板
	// **/* 只能加载子目录模板
	// /** 理论上可以加载所有的文件和目录，但编译直接报错
	// 最好的方式为filepath.WalkDir  + LoadHTMLFiles() 的方式加载模板
	router.LoadHTMLGlob("D:\\code\\my_repo\\zf-project\\learn\\gin\\examples\\html_template\\templates\\**\\*")
	router.LoadHTMLFiles("D:\\code\\my_repo\\zf-project\\learn\\gin\\examples\\html_template\\templates\\index.tmpl")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})
	router.Run(":8080")
}
