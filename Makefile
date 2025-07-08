BINARY_NAME=scrapbook
SYS_CONF_DIR=/etc/scrapbook
SYS_CONF_DIR=/var/lib/scrapbook

all: build

build:
	go build -ldflags "-X 'github.com/LMBishop/scrapbook/pkg/constants.SysConfDir=${SYS_CONF_DIR}' -X 'github.com/LMBishop/scrapbook/pkg/constants.SysDataDir=${SYS_DATA_DIR}'" -o ${BINARY_NAME} main.go

runlocal:
	PWD=$(shell pwd)
	mkdir -p out
	go build -ldflags "-X 'github.com/LMBishop/scrapbook/pkg/constants.SysConfDir=${PWD}/out/config' -X 'github.com/LMBishop/scrapbook/pkg/constants.SysDataDir=${PWD}/out/data'" -o out/${BINARY_NAME} main.go
	cd out; ./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}