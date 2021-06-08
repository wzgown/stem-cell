export GOPROXY=https://goproxy.io,direct 
export GOPRIVATE=amtlink.cn
export GOSUMDB=off
export CGO_ENABLED=0 

all:
	go mod tidy
	env CGO_ENABLED=0 go build --tags netgo --tags jsoniter
	go test --tags netgo --tags jsoniter
