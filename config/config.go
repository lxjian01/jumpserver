package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type AppConfig struct {
	Version    string        `yaml:"version"`
	Env        string        `yaml:"env"`
	Httpd      HttpdConfig   `yaml:"httpd"`
	Sshd       SshdConfig    `yaml:"sshd"`
	RecordDir  string        `yaml:"recordDir"`
	LinuxUser  string        `yaml:"linuxUser"`
	PoolNum    int           `yaml:"poolNum"`
	Log        LogConfig     `yaml:"log"`
	Mysql      MysqlConfig   `yaml:"mysql"`
	Redis      RedisConfig   `yaml:"redis"`
}

type HttpdConfig struct {
	Host string
	Port int
}

type SshdConfig struct {
	Host string
	Port int
}

type LogConfig struct {
	Dir       string
	Name      string
	Format    string
	RetainDay int8
	Level     string
}

type MysqlConfig struct {
	DbJump DbConfig
}
type DbConfig struct {
	Host        string
	Port        int
	Db          string
	User        string
	Password    string
	Charset     string
}

type RedisConfig struct {
	Host        string
	Port        int
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

func InitConfig() *AppConfig {
	var appConf *AppConfig
	if err :=viper.Unmarshal(&appConf); err !=nil{
		fmt.Errorf("Unmarshal config error by %v \n",err)
		os.Exit(1)
	}
	return appConf
}