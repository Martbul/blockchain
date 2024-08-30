package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

func IsFoundHost(host string, port uint16) bool { //checking if a connection can be made to the given `host:port`
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, 1 * time.Second)
	if err != nil{
		fmt.Printf("%s %v\n",target,err)
		return false  	
	}
	return true
}


// 192.168.0.10:5000
// 192.168.0.11:5000
// 192.168.0.12:5000
// 192.168.0.10:5001
// 192.168.0.10:5002
// 192.168.0.10:5003


var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

//! 100% works
func FindNeighbors(myHost string, myPort uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string{
	address := fmt.Sprintf("%s:%d", myHost,myPort) //creating an address for yourself
	// fmt.Println(address) //127.0.0.1:5000

	m := PATTERN.FindStringSubmatch(myHost)
	
	if m == nil {
		return nil
	}

	prefixHost := m[1]  //m[1] is 192.168.0
		// fmt.Println(prefixHost) //[127.0.0.1 127.0.0. 0. 1]

	lastIp,_ := strconv.Atoi(m[len(m)-1]) // this is after the 0 -> 10 / 11 / 12 ... .striconv.Atoi make the strin '10' to int 10
	neighbors := make([]string,0) //mking an empty string slice
	
	for port := startPort; port <= endPort; port++ {
		for ip := startIp; ip<= endIp; ip++ {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp + int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port){ // if guessTarget is not your address then add the address to the neighbours slice
			
				neighbors = append(neighbors, guessTarget) //! should return [127.0.0.1:5001 127.0.0.1:5002] but return [127.0.0.2:5000 127.0.0.3:5000 127.0.0.4:5000 127.0.0.1:5001 127.0.0.2:5001 127.0.0.3:5001 127.0.0.4:5001 127.0.0.1:5002 127.0.0.2:5002 127.0.0.3:5002 127.0.0.4:5002]
			}
		}
	}
	return neighbors
}


func GetHost() string {
	hostname, err := os.Hostname() //getting your own hostname
	if err != nil{
		return "127.0.0.1"
	}
	// fmt.Println("hostname",hostname) //DESKTOP-FEFP1UD
	address, err := net.LookupHost(hostname)
	if err != nil{
		return "127.0.0.1"
	}
	fmt.Println("address",address[len(address)-1]) //[fe80::770a:9f6e:af83:45e%Ethernet 192.168.0.106]
	return address[len(address)-1]
	// return address[0] //! POSSIBLE FUCK UP
}