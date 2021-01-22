SHELL = /bin/bash

build:
	@GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -o modbus-go
	@docker build -t ke-modbus-mapper .
	@rm modbus-go

push: build
	@$(eval REV := $(shell git rev-parse HEAD|cut -c 1-8))
	@docker tag ke-modbus-mapper 10.1.11.196/library/ke-modbus-mapper-common:ke-modbus-mapper-$(REV)
	@docker push  10.1.11.196/library/ke-modbus-mapper-common:ke-modbus-mapper-$(REV)
