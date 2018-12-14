PROJ_PATH = $(GOPATH)/src/github.com/aaronflower/dzone-shipping
build:
	protoc -I. --go_out=plugins=micro:$(PROJ_PATH)/service.consignment proto/consignment/consignment.proto
	GOOS=linux GOARCH=amd64 go build
	docker build --rm -t service.consignment .

run:
	docker run -p 50051:50051  \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns service.consignment
	

