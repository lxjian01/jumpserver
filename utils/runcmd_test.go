package utils

import (
	"fmt"
	"testing"
)

func TestRunCommand(t *testing.T) {
	out,err := RunCommand("hostname")
	if err != nil{
		fmt.Errorf("Exec cmd error %v \n",err)
	}else{
		fmt.Println("Exec cmd return",out)
	}
}
