package main

import (
	"fmt"

	"github.com/martbul/golang-blockchain/utils"
)

func main() { //ex: your address is 127.0.0.1:5000
	// fmt.Println(utils.FindNeighbors("127.0.0.1", 5000, 0, 3,5000,5003))
	// fmt.Println(utils.IsFoundHost("127.0.0.1",5000))
	fmt.Println(utils.GetHost())
}