package routers

import (
	"github.com/gin-gonic/gin"
	"jumpserver/server/httpd/controllers"
)

func WebsocketRoutes(route *gin.Engine) {
	webshell := route.Group("/webshell")
	{
		webshell.GET("/docker/:project_code/:module_code/:host/:deploy_job_host_id/:token", controllers.WsConnectDocker)

		webshell.GET("/linux/:project_code/:module_code/:host/:deploy_job_host_id/:token", controllers.WsConnectLinux)
	}
}
