#!C:/software/Git/bin/bash
export GOPROXY="https://goproxy.cn,direct"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server_linux_amd64 .
#CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o server_linux_arm64 .

