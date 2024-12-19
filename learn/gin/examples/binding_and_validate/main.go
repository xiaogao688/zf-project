package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"net/http"
)

var route *gin.Engine

func main() {

	route = gin.Default()
	bindURI()
	validate()
	route.Run(":8088")
}

// 绑定uri

type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func bindURI() {
	// 绑定uri
	route.GET("/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	})
}

// 自定义验证函数
func validateCoolName(fl validator.FieldLevel) bool {
	return fl.Field().String() == "coolname"
}

type User struct {
	Name string `json:"name" binding:"required,coolname"`
}

func validate() {
	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("coolname", validateCoolName)
	}

	route.POST("/validate", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Validation success!"})
	})

}

type Upload struct {
	Title string                `form:"title" binding:"required"`
	File  *multipart.FileHeader `form:"file" binding:"required"`
}

func UploadHandler(c *gin.Context) {
	var upload Upload
	if err := c.ShouldBind(&upload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := upload.File
	c.SaveUploadedFile(file, "./uploads/"+file.Filename)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully!"})
}
