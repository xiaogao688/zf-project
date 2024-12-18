package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

func main() {
	GroupRouter()
}

func httpRequest(c *gin.Context) {
	c.JSON(200, gin.H{"month": c.Request.Method, "url": c.Request.URL})
}

func ParameterParse() {

}

func GroupRouter() {
	var eg errgroup.Group

	// 一进程多端口
	insecureServer := &http.Server{
		Addr: ":8080",
		Handler: func() http.Handler {
			r := gin.Default()
			groupV1 := r.Group("/v1")
			gv1A := groupV1.Group("a")
			gv1A.GET("b", httpRequest)
			gv1A.HEAD("b", httpRequest)
			return r
		}(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	secureServer := &http.Server{
		Addr: ":8443",
		Handler: func() http.Handler {
			r := gin.Default()
			groupV1 := r.Group("/v1")
			gv1A := groupV1.Group("a")
			gv1A.GET("b", httpRequest)
			gv1A.HEAD("b", httpRequest)
			return r
		}(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 一进程多服务
	eg.Go(func() error {
		err := insecureServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		return err
	})

	eg.Go(func() error {
		err := secureServer.ListenAndServeTLS("D:\\code\\my_repo\\zf-project\\learn\\gin\\ca\\server.pem", "D:\\code\\my_repo\\zf-project\\learn\\gin\\ca\\server.key")
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
