GOPATH=$(shell go env GOPATH)
GIN_PORT=3000
APP_PORT=$$(( $(GIN_PORT) + 1 ))

run:
	@$(GOPATH)/bin/gin --port=$(GIN_PORT) --appPort=$(APP_PORT) --bin='app-bin' --immediate --buildArgs='-v -x -mod=readonly -buildvcs=false' run main.go  