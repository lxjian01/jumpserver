package server

import (
	"jumpserver/config"
	"jumpserver/server/httpd"
	"jumpserver/server/sshd"
)

func StartServer(appConf *config.AppConfig){
	httpd.StartHttpdServer(&appConf.Httpd)
	sshd.StartSshdServer(&appConf.Sshd)
}
