package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/martbul/golang-blockchain/utils"
)

const (
	MINING_DIFFICULTY = 3 //!mining dificulty IRL will be larger
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20

	BLOCKCHAIN_PORT_RANGE_START       = 5000
	BLOCKCHAIN_PORT_RANGE_END         = 5003
	NEIGHBOR_IP_RANGE_START           = 0
	NEIGHBOR_IP_RANGE_END             = 1
	BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC = 20
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string // used to identify the minier
	port              uint16
	mux               sync.Mutex
	neighbors         []string
	muxNeighbors      sync.Mutex
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

// ! 100% works
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

// ! 100% works
func (b *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	v := &struct {
		Timestamp    *int64          `json:"timestamp"`
		Nonce        *int            `json:"nonce"`
		PreviousHash *string         `json:"previous_hash"`
		Transactions *[]*Transaction `json:"transactions"`
	}{
		Timestamp:    &b.timestamp,
		Nonce:        &b.nonce,
		PreviousHash: &previousHash,
		Transactions: &b.transactions,
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	ph, _ := hex.DecodeString(*v.PreviousHash)
	copy(b.previousHash[:], ph[:32])
	return nil
}

// ! 100% works
func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

// ! 100% works
func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

// ! 100% works
func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
}

// ! 100% works
func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
}

// ! 100% works
func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	v := &struct {
		Blocks *[]*Block `json:"chain"`
	}{
		Blocks: &bc.chain,
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	return nil
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		sender,
		recipient,
		value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockhain_address %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blockhain_address %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		//! important , is that the filed must start with cappital letters to make it a public property
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

//! 100% works
func (t *Transaction) UnmarshalJSON(data []byte) error {
	v := &struct {
		Sender    *string  `json:"sender_blockchain_address"`
		Recipient *string  `json:"recipient_blockchain_address"`
		Value     *float32 `json:"value"`
	}{
		Sender:    &t.senderBlockchainAddress,
		Recipient: &t.recipientBlockchainAddress,
		Value:     &t.value,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

// ! 100% works
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// ! 100% works
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)         // formating a block to json
	return sha256.Sum256([]byte(m)) //hashing the json
}

// ! 100% works
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	for _, n := range bc.neighbors { // this syncs the transaction pool's between all nodes, when a block is created and transaction pool is emptyed it will empty all transaction pools(I belive)
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %d %s \n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}

	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

// ! 100% works
func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

// ! 100% works
func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	if isTransacted {
		for _, n := range bc.neighbors {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &TransactionRequest{&sender, &recipient, &publicKeyStr, &value, &signatureStr}
			m, _ := json.Marshal(bt)
			buf := bytes.NewBuffer(m)
			endpoint := fmt.Sprintf("http://%s/transactions", n)
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", endpoint, buf)
			resp, _ := client.Do(req)
			log.Printf("%v", resp)
		}
	}

	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER { // in case of transaction from the blockchain to a miner a varification is not needed
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: not enougn balance in a wallet")
		// 	return false
		// }
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range transactions {
		transactions = append(transactions,
			NewTransaction(t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value),
		)
	}

	return transactions
}

// ! 100% works
func (bc *Blockchain) Mining() bool {
	// mining cann be called multyple time as once and you DON'T want that to happen
	bc.mux.Lock()         // when the method is called -> locking it
	defer bc.mux.Unlock() // when the method finishes it will be unlocked

	if len(bc.transactionPool) == 0 {
		return false // when the transaction pool is empty there is no need for mining (IRL mining time is min 10mn so the transaction pool always needs mining)
	}

	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	// log.Println("action=mining, status=success")

	for _, n := range bc.neighbors {
		fmt.Println("NNNNNNNN",n) //172.20.160.1:5001  172.20.160.1:5002 (results from the 2 calls)
		endpoint := fmt.Sprintf("http://%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, nil)
		resp, _ := client.Do(req)
		fmt.Println("WWWWWWWWWWWWW", resp)
		log.Printf("%v", resp)
	}
	return true
}

// ! 100% works
func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second * MINING_TIMER_SEC, bc.StartMining) //executing the StartMining in a loop every 20sec(MINING_TIMER_SEC)
}

// ! 100% works
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

// ! 100% works
func (bc *Blockchain) ValidChain(chain []*Block) bool { // checks if the entire blovkvhain is valid, iterates through each block and enshures that:1 Hash continuity(each block's previousHash must match the hash of the next block in the chain) and 2 Valid Proof(each block must meet the required difficulty level 000)
	preBlock := chain[0] //the first block in the blockchain
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.previousHash != preBlock.Hash() {
			return false
		}

		if !bc.ValidProof(b.Nonce(), b.PreviousHash(), b.Transactions(), MINING_DIFFICULTY) { //! here is returning false
			fmt.Println("FALLLLLSSEEEEEE")
			return false
		}

		preBlock = b
		currentIndex += 1
	}
	return true
}

//! 100% works
func (bc *Blockchain) Chain() []*Block {
	return bc.chain
}

//! 100% works
func (bc *Blockchain) ResolveConflicts() bool {
	var longestChain []*Block = nil
	maxLength := len(bc.chain)

	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/chain", n)
		resp, _ := http.Get(endpoint)
		if resp.StatusCode == 200 {
			var bcResp Blockchain
			decoder := json.NewDecoder(resp.Body)

			_ = decoder.Decode(&bcResp)

			chain := bcResp.Chain()
			for _ ,v := range chain{
				//ha owete sa im werni, problema e s valid proof
									fmt.Println("DDDDDDDDDDDDDDDDDDDDDDD", *v)

			}

			fmt.Println("q1q1q1q1q1q1q1q1q1q1q1q1q1q1q", bc.ValidChain(chain)) //! chain is not valid. Options: not getting the chain right; wrong code for ValidChain; the chain is build wrong

			if len(chain) > maxLength && bc.ValidChain(chain) {

				maxLength = len(chain)
				longestChain = chain
			}
		}
	}

	if longestChain != nil {
		bc.chain = longestChain
		log.Printf("Resolve conflicts replaced")
		return true
	}
	log.Printf("Resolve conflicts NOT replaced")
	return false
}

// ! 100% works
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool { // check if a given nonce results in a valid hash
	zeros := strings.Repeat("0", difficulty) // the block must start with difficulty number of 0's
	guessBlock := Block{0, nonce, previousHash, transactions} //  creating a guessBlock, using the given nonce, previousHash, and transactions
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())

	return guessHashStr[:difficulty] == zeros
}

// ! 100% works
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}

	return nonce
}

func (b *Block) Print() {
	fmt.Printf("timestamp %d\n", b.timestamp)
	fmt.Printf("nonce %d\n", b.nonce)
	fmt.Printf("previousHas %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return true
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: ar.Amount,
	})
}

// ! 100% works
func (bc *Blockchain) Run() { //this will be called from your web API (blockchain_server.go)
	bc.StartSyncNeighbors()
	bc.ResolveConflicts()
}

// ! 100% works
func (bc *Blockchain) SetNeighbors() {
	bc.neighbors = utils.FindNeighbors(
		utils.GetHost(), bc.port,
		NEIGHBOR_IP_RANGE_START, NEIGHBOR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START, BLOCKCHAIN_PORT_RANGE_END)
	log.Printf("%v", bc.neighbors)
}

// !100% works
func (bc *Blockchain) SyncNeighbors() {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	bc.SetNeighbors()
}

// ! 100% works
func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(time.Second*BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC, bc.StartSyncNeighbors)
}

// ! 100% works
func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

// ! 100% works
func (b *Block) Nonce() int {
	return b.nonce
}

// ! 100% works
func (b *Block) Transactions() []*Transaction {
	return b.transactions
}
