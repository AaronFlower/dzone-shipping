PROJ_PATH = $(GOPATH)/src/github.com/aaronflower/dzone-shipping
build:
	protoc -I. --go_out=plugins=grpc:$(PROJ_PATH)/service.consignment proto/consignment/consignment.proto
