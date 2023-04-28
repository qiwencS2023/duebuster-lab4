all: proto build

proto:
	(cd storage && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative storage.proto)
	(cd coordinator && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative coordinator.proto)

build:
	go build -o dist/storage ./storage/*.go
	go build -o dist/coordinator ./coordinator/*.go

test:
	go test ./coordinator/
	go test ./storage/

clean:
	rm -rf ./dist