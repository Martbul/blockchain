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







PORT ?= 5000

buildBlockchainServer:
	go build -o bin/blockchain_server ./blockchain_server

buildWalletServer:
	go build -o bin/wallet_server ./wallet_server


buildCMD:
	go build -o bin/CMD ./cmd

runBCServer: buildBlockchainServer
	./bin/blockchain_server -port=$(PORT)

runWServer: buildWalletServer
	./bin/wallet_server

runCMD: buildCMD
	./bin/CMD

runMultipleBCServers:
	start cmd /c "$(MAKE) runBCServer PORT=5000" && \
	start cmd /c "$(MAKE) runBCServer PORT=5001" && \
	start cmd /c "$(MAKE) runBCServer PORT=5002"

test:
	go test -v ./...
