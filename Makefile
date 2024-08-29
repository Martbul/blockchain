# # Build the blockchain binary
# buildBlockchain:
# 	go build -o bin/blockchain main.go

# # Build the blockchain server binary
# buildServer:
# 	go build -o bin/blockchain_server ./blockchain_server/main.go

# # Run the blockchain binary
# runBlockchain: buildBlockchain
# 	./bin/blockchain

# # Run the blockchain server binary
# runServer: buildServer
# 	./bin/blockchain_server

# # Run tests
# test:
# 	go test -v ./...


buildBlockchainServer:
	go build -o bin/blockchain_server ./blockchain_server

buildWalletServer:
	go build -o bin/wallet_server ./wallet_server

runBCServer: buildBlockchainServer
	./bin/blockchain_server

runWServer: buildWalletServer
	./bin/wallet_server

test:
	go test -v ./...
