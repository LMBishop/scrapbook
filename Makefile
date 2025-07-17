BINARY_NAME=scrapbook
SYS_CONF_DIR=/etc/scrapbook
SYS_DATA_DIR=/var/lib/scrapbook
SYSTEMD_DIR=/etc/systemd/system

all: build

build:
	go build -ldflags "-X 'github.com/LMBishop/scrapbook/pkg/constants.SysConfDir=${SYS_CONF_DIR}' -X 'github.com/LMBishop/scrapbook/pkg/constants.SysDataDir=${SYS_DATA_DIR}'" -o ${BINARY_NAME} main.go

.PHONY: runlocal
runlocal:
	PWD=$(shell pwd)
	mkdir -p runlocal
	go build -ldflags "-X 'github.com/LMBishop/scrapbook/pkg/constants.SysConfDir=${PWD}/runlocal/config' -X 'github.com/LMBishop/scrapbook/pkg/constants.SysDataDir=${PWD}/runlocal/data'" -o runlocal/${BINARY_NAME} main.go
	cd runlocal; ./${BINARY_NAME}

install: build
	install -Dm755 ${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

installconf:
	install -Dm755 dist/scrapbook.service ${SYSTEMD_DIR}/
	install -Dm755 dist/config.toml ${SYS_CONF_DIR}/config.toml
	if ! getent passwd scrapbook > /dev/null; then \
		useradd --system --create-home --shell /usr/sbin/nologin --home-dir ${SYS_DATA_DIR} scrapbook ;\
	fi

clean:
	go clean
	rm ${BINARY_NAME}