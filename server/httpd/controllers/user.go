package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Test(c *gin.Context){
	c.HTML(http.StatusOK, "video.html", gin.H{
		"title": "hello Go",
	})
}
