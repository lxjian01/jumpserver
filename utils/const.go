package utils

const (
	ProfieEnv    = "env"
	LocalHostEnv = "fisher.advertise.localhost.ip"
	ConfigRealIP = "config.real.ip"
	DebugModEnv  = "fisher.debug"

	EnvDev  = "dev"
	EnvTest = "test"
	EnvPrd  = "prd"

	LogTypeLogback = "log"
	LogTypeJson    = "json"
	LogTypeLine    = "line"
	LogTypeNginx   = "ngx"

	ReportLogTypeLogback = "log4j"
	ReportLogTypeNginx   = "nginx"

	CollectTypeWatch = "watch"
	CollectTypeConf  = "conf"
)

var Version string
