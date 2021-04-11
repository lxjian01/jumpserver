package httpd

import (
	"github.com/gin-gonic/gin"
	"jumpserver/config"
	"jumpserver/log"
	"jumpserver/server/httpd/middlewares"
	"jumpserver/server/httpd/routers"
	"net"
	"strconv"
)

func StartHttpdServer(c *config.HttpdConfig) {
	router := gin.Default()
	// 添加自定义的 logger 间件
	router.Use(middlewares.Logger(), gin.Recovery())
	router.Use(middlewares.Auth(), gin.Recovery())
	// 添加路由
	routers.UserRoutes(router)      //Added all user routers
	routers.WebsocketRoutes(router) //Added all websocket routers
	// 拼接host
	Host := c.Host
	Port := strconv.Itoa(c.Port)
	addr := net.JoinHostPort(Host, Port)
	log.Info("Start HTTP server at", addr)
	err1 := router.Run(addr)
	if err1 != nil{
		log.Error("Start server error by",err1)
	}
	log.Info("Start server ok")
}
