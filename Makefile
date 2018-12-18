.ONESHELL:
build:
	cd ./consignment-cli && make build
	cd ./service.vessel && make build
	cd ./service.consignment && make build

clean:
	cd ./consignment-cli && go clean
	cd ./service.vessel && go clean 
	cd ./service.consignment && go clean
