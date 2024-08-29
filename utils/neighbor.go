package utils

import (
	"fmt"
	"net"
	"time"
)

func IsFoundHost(host string, port uint16) bool { //checking if a connection has been made with tcp
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target,1*time.Second)
	if err != nil{
		fmt.Printf("%s %v\n",target,err)
		return false  	
	}
	return true
}