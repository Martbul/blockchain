// package main

// import (
// 	"flag"
// 	"log"
// )

// func init() {
// 	log.SetPrefix("Wallet Server: ")
// }

// func main() {
// 	port := flag.Uint("port", 8080, "TCP Port Number for Wallet Server")
// 	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain Gateway")
// 	flag.Parse()

// 	app:= NewWalletServer(uint16(*port), *gateway)
// 	app.Run()
// }

package main

import (
	"flag"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	port := flag.Uint("port", 8080, "TCP Port Number for Wallet Server")
	gatewayPort := flag.Uint("gatewayPort", 5000, "TCP Port Number for Blockchain Gateway")
	flag.Parse()

	gateway := fmt.Sprintf("http://127.0.0.1:%d", *gatewayPort)

	app := NewWalletServer(uint16(*port), gateway)
	app.Run()
}
