.ONESHELL:
build:
	cd ./consignment-cli && make build
	cd ./service.vessel && make build
	cd ./service.consignment && make build
	cd ./service.user && make build
	cd ./user-cli && make build

clean:
	cd ./consignment-cli && go clean
	cd ./service.vessel && go clean 
	cd ./service.consignment && go clean
	cd ./service.user && go clean 
	cd ./user-cli && go clean
