/**************************************
 * @Author: mazhuang
 * @Date: 2021-09-23 14:19:42
 * @LastEditTime: 2021-10-12 11:11:06
 * @Description:
 **************************************/

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seek4self/logger"
)

func main() {
	log, err := logger.NewDefault()
	if err != nil {
		panic(err)
	}

	r := gin.New()

	r.Use(gin.LoggerWithConfig(log.GinLogConfig()))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
