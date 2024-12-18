package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/json", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// 输出: {"lang":"GO语言","tag":"\u003cbr\u003e"}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/ascii_json", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// 返回 ASCII-only 的 JSON 数据，非 ASCII 字符会被转义为 Unicode 序列。
		// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	r.GET("/indented_json", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// 返回带缩进格式的 JSON 数据，便于调试或阅读。
		// 输出 :
		//{
		//    "lang": "GO语言",
		//    "tag": "\u003cbr\u003e"
		//}
		c.IndentedJSON(http.StatusOK, data)
	})

	// SecureJSON
	// 解决 JSON 数据可能被 JSON Hijacking 攻击的问题。
	// 你也可以使用自己的 SecureJSON 前缀
	// r.SecureJsonPrefix(")]}',\n")

	r.GET("/secure_json", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// 将输出：while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})

	// JSON 使用 unicode 替换特殊 HTML 字符，例如 < 变为 \ u003c。如果要按字面对这些字符进行编码，则可以使用 PureJSON。
	r.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
