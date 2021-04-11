package routers

import (
	"github.com/gin-gonic/gin"
	"jumpserver/server/httpd/controllers"
)

func UserRoutes(route *gin.Engine) {
	user := route.Group("/user")
	{
		user.GET("/test1", controllers.Test)
		//user.POST("/test", controllers.Test)
	}
}