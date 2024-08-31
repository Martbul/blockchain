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
	./bin/wallet_server -port=$(PORT) -gatewayPort=$(GATEWAY_PORT)

runCMD: buildCMD
	./bin/CMD

test:
	go test -v ./...
