all:
	go mod tidy
	env GOPROXY=https://goproxy.io,direct go get
	env GOPROXY=https://goproxy.io,direct CGO_ENABLED=0 go build --tags netgo --tags jsoniter
	env GOPROXY=https://goproxy.io,direct CGO_ENABLED=0 go test --tags netgo --tags jsoniter
