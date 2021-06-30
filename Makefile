NAME:=docker
BINARY=terraform-provider-${NAME}
VERSION=0.0.2
OS_ARCH=linux_amd64

default: build

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_v${VERSION}_darwin_amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_v${VERSION}_linux_amd64
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_v${VERSION}_openbsd_amd64
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_v${VERSION}_windows_amd64
	cd bin && \
	zip -r9 ${BINARY}_${VERSION}_darwin_amd64.zip ${BINARY}_v${VERSION}_darwin_amd64 && \
	zip -r9 ${BINARY}_${VERSION}_linux_amd64.zip ${BINARY}_v${VERSION}_linux_amd64 && \
	zip -r9 ${BINARY}_${VERSION}_openbsd_amd64.zip ${BINARY}_v${VERSION}_openbsd_amd64 && \
	zip -r9 ${BINARY}_${VERSION}_windows_amd64.zip ${BINARY}_v${VERSION}_windows_amd64 && \
	: && \
	sha256sum ${BINARY}_${VERSION}_*.zip > ${BINARY}_${VERSION}_SHA256SUMS && \
	gpg --detach-sign ${BINARY}_${VERSION}_SHA256SUMS

