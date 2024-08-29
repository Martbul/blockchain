package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/martbul/golang-blockchain/block"
	"github.com/martbul/golang-blockchain/wallet"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain{
	bc, ok := cache["blockchain"] //chacking if the blockchain exist in the cashe
	if !ok { // in the beggining it doesnt exist in the cashe
		minersWallet := wallet.NewWallet() //registering the miner's address
		bc = block.NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc //adding the block, that was just created, in the cashe
		log.Printf("private_key %v", minersWallet.PrivateKeyStr())
		log.Printf("public_key %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address %v", minersWallet.BlockchainAddress())
	}
	return bc
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	 switch req.Method {
	 case http.MethodGet:
			w.Header().Add("Content-Type","application/json")
			bc := bcs.GetBlockchain()
			m, _ := bc.MarshalJSON()
			io.WriteString(w, string(m[:]) )
		default:
			log.Printf("ERROR: Invalid http method")
	 }
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}
