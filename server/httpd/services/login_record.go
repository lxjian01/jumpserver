package services

import (
	"jumpserver/db"
	"jumpserver/log"
	"jumpserver/utils"
)

// add login record
func InsertLoginRecord(projectCode string, moduleCode string,host string,deployJobJostId int,userCode string) int64 {
	insertSql := "INSERT INTO login_record(project_code, module_code, deploy_job_host_id,ip,user_code) VALUES (?, ?, ?, ?, ?)"
	log.Info(projectCode, moduleCode,deployJobJostId,host,userCode)
	ret,err := db.DbJump.Exec(insertSql, projectCode, moduleCode,deployJobJostId,host,userCode)
	if err != nil{
		log.Errorf("Add %s login %s error by %v \n",userCode,host,err)
		return 0
	}else{
		log.Infof("Add %s login %s ok \n",userCode,host)
		id,err1 := ret.LastInsertId()
		if err1 != nil{
			return 0
		}
		return id
	}
}

// update login record
func UpdateLoginRecord(id int64) {
	datetime := utils.DateNowFormatYmdhms()
	updateSql := "update login_record set logout_time = ? where id = ?"
	_,err := db.DbJump.Exec(updateSql, datetime,id)
	if err != nil{
		log.Errorf("Update login record id is %d error by %v \n",id,err)
	}else{
		log.Infof("Update login record id is %d ok \n",id)
	}
}

// 添加Docker操作记录
func AddDockerOperRecord(cmd string,userCode string,host string) {
	insertSql := "INSERT INTO docker_record(user_code,ip,cmd) VALUES (?, ?, ?)"
	_,err :=db.DbJump.Exec(insertSql, userCode, host,cmd)
	if err != nil{
		log.Errorf("【docker】：%s oper %s cmd %s record error by %v \n",userCode,host,cmd,err)
	}else{
		log.Infof("【docker】：%s oper %s cmd: %s \n",userCode,host,cmd)
	}
}

// 添加lunux操作记录
func AddLinuxOperRecord(cmd string,userCode string,host string) {
	insertSql := "INSERT INTO linux_record(user_code,ip,cmd) VALUES (?, ?, ?)"
	_,err :=db.DbJump.Exec(insertSql, userCode, host,cmd)
	if err != nil{
		log.Errorf("【linux】：%s oper %s cmd %s record error by %v \n",userCode,host,cmd,err)
	}else{
		log.Infof("【linux】：%s oper %s cmd: %s \n",userCode,host,cmd)
	}
}