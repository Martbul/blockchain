package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("blockchain: ")
}


func main() {
	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Server")
	flag.Parse()

	app := NewBlockchainServer(uint16(*port)) //app is type *BlockchainServer, so it has the method Run(), that executes this in the blockchain_server.go: func (bcs *BlockchainServer) Run() {
																																																				// 	bcs.GetBlockchain().Run()
																																																				// 	http.HandleFunc("/", bcs.GetChain)
																																																				// 	http.HandleFunc("/transactions", bcs.Transactions)
																																																				// 	http.HandleFunc("/mine", bcs.Mine)
																																																				// 	http.HandleFunc("/mine/start", bcs.StartMine)
																																																				// 	http.HandleFunc("/amount", bcs.Amount)
																																																				// 	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
																																																				// }

	app.Run()
}
