package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"jumpserver/config"
	"jumpserver/db"
	"jumpserver/log"
	"jumpserver/pools"
	"jumpserver/server"
	"os"
	"path/filepath"
)

var (
	daemon bool
	appConf *config.AppConfig
	lockfile = "/var/run/jumpserver.pid"
)

//Execute方法触发init方法
func init() {
	//初始化配置文件转化成对应的结构体
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(versionCmd)
}

//项目启动调用的入口方法
func Execute() {
	//初始化Cobra
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "jumpserver",
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting init system log")
		log.Init(&appConf.Log)
		fmt.Println("Init system log ok")

		log.Info("Starting init pool")
		pools.InitPool(appConf.PoolNum)
		log.Info("Init pool ok")

		log.Info("Starting init mysql")
		db.InitMysqlDb(&appConf.Mysql)
		log.Info("Init mysql ok")

		log.Info("Starting init redis")
		db.InitRedisPool(&appConf.Redis)
		log.Info("Init redis ok")
	},
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer(appConf)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show jumpserver version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Jumpserver version is",appConf.Version)
		os.Exit(0)
	},
}

//通过viper初始化配置文件到结构体
func initConfig() {
	dir,_ := os.Getwd()
	configPath := filepath.Join(dir,"config")
	// 设置读取的文件路径
	viper.AddConfigPath(configPath)
	// 设置读取的文件名
	viper.SetConfigName("jumpserver")
	// 设置文件的类型
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Errorf("Read config error by %v \n",err)
		os.Exit(1)
	}
	appConf = config.InitConfig()
}