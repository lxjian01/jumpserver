package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"jumpserver/config"
	"jumpserver/log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	DbJump *sqlx.DB
)
// 初始化mysql数据库
func InitMysqlDb(c *config.MysqlConfig){
	jumpConstr := getDbStr(c.DbJump)
	log.Info("Jump db connect info",jumpConstr)
	DbJump = sqlx.MustConnect("mysql", jumpConstr)
	// 设置连接可以被使用的最长有效时间
	DbJump.SetConnMaxLifetime(10 * time.Minute)
	// 设置连接池中的保持连接的最大连接数
	DbJump.SetMaxIdleConns(10)
	// 设置打开数据库的最大连接数
	DbJump.SetMaxOpenConns(30)
}

func getDbStr(c interface{}) string{
	var build strings.Builder
	var constr string
	cValue := reflect.ValueOf(c)
	User := cValue.FieldByName("User").String()
	password := cValue.FieldByName("Password").String()
	Host := cValue.FieldByName("Host").String()
	Port := cValue.FieldByName("Port").Int()
	Db := cValue.FieldByName("Db").String()
	Charset := cValue.FieldByName("Charset").String()
	build.WriteString(User)
	build.WriteString(":")
	build.WriteString(password)
	build.WriteString("@tcp(")
	build.WriteString(Host)
	build.WriteString(":")
	build.WriteString(strconv.FormatInt(Port, 10))
	build.WriteString(")/")
	build.WriteString(Db)
	build.WriteString("?charset=")
	build.WriteString(Charset)
	constr = build.String()
	return constr
}

func CloseMysql() {
	if err := DbJump.Close();err != nil{
		log.Error("Close jump db error by",err)
	}
}