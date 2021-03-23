SHELL = /bin/bash

build:
	@GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -o modbus-go
	@docker build -t gopherlv/ke-modbus-mapper .
	@rm modbus-go

push: build
	@$(eval REV := $(shell git rev-parse HEAD|cut -c 1-8))
	@docker tag gopherlv/ke-modbus-mapper gopherlv/ke-modbus-mapper-shutter-armv7:$(REV)
	@docker push gopherlv/ke-modbus-mapper-shutter-armv7:$(REV)
